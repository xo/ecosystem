package proto

import (
	"fmt"

	pb "github.com/xo/ecosystem/proto/xo"
	"github.com/xo/ecosystem/types"
	"google.golang.org/protobuf/compiler/protogen"
)

// ConvertedField is the converted form of protogen.Field, containing
// information on items to add to integrate the field into the schema.
type ConvertedField struct {
	Field       *types.Field
	Indexes     []types.Index
	ForeignKeys []types.ForeignKey
	ExtraTables []types.Table
}

// ConvertField converts the provided protogen.Field to ConvertedField.
func (c Converter) ConvertField(field *protogen.Field) (*ConvertedField, error) {
	// Check for ignored fields.
	options := fieldOpts(field)
	if options.Ignore {
		return &ConvertedField{}, nil
	}
	msgOpts := messageOpts(field.Message)
	if options.EmbedAsJson || msgOpts.Ignore || msgOpts.EmbedAsJson {
		// Embed fields ignored reftypes and those explicitly marked as JSON.
		col := &types.Field{
			Name: field.Desc.JSONName(),
			Type: types.Type{Type: "[]byte"},
		}
		return &ConvertedField{
			Field: col,
		}, nil
	}
	typ, basicTypes, err := c.goType(field)
	if err != nil {
		return nil, err
	}
	switch {
	case basicTypes:
		return c.convertSimpleFields(field, typ)
	case typ.IsArray:
		return c.convertArrayFields(field, typ)
	default:
		return c.convertReferencedFields(field, typ)
	}
}

// convertSimpleFields converts fields with basic types such as integers.
func (c Converter) convertSimpleFields(field *protogen.Field, typ types.Type) (*ConvertedField, error) {
	parentTblName, err := c.messageTableName(field.Parent, true)
	if err != nil {
		return nil, err
	}
	// Construct the column.
	options := fieldOpts(field)
	col := types.Field{
		Name:    field.Desc.JSONName(),
		Type:    typ,
		Default: options.DefaultValue,
	}
	// Add unique index if marked in options.
	var idx []types.Index
	if options.Index != pb.FieldOverride_NONE {
		idxName := fmt.Sprintf("%s_%s_idx", parentTblName, col.Name)
		idx = append(idx, types.Index{
			Name:     idxName,
			Fields:   []types.Field{col},
			IsUnique: options.Index == pb.FieldOverride_UNIQUE,
		})
	}
	// Add references if marked in options.
	var fk []types.ForeignKey
	if options.Ref != nil {
		ref, err := c.createRef(field, col, options.Ref)
		if err != nil {
			return nil, err
		}
		fk = append(fk, *ref)
	}
	return &ConvertedField{
		Field:       &col,
		Indexes:     idx,
		ForeignKeys: fk,
	}, nil
}

// createRef creates a reference specified by ref from field.
func (c Converter) createRef(field *protogen.Field, col types.Field, ref *pb.Ref) (*types.ForeignKey, error) {
	parentTblName, err := c.messageTableName(field.Parent, true)
	if err != nil {
		return nil, err
	}
	// Search for the referenced type.
	typ, err := c.searchForType(ref.TypeSuffix)
	if err != nil {
		return nil, err
	}
	// Search for the referenced field.
	var referencedFieldProto *protogen.Field
	for _, f := range typ.Fields {
		if string(f.Desc.Name()) == ref.FieldName {
			referencedFieldProto = f
			break
		}
	}
	if referencedFieldProto == nil {
		return nil, fmt.Errorf(
			"field %q in %q references field %q in %q which does not exist",
			field.Desc.Name(), field.Parent.Desc.FullName(),
			ref.FieldName, typ.Desc.FullName(),
		)
	}
	referencedField, err := c.ConvertField(referencedFieldProto)
	if err != nil {
		return nil, err
	}
	// Check that the referenced field is valid as a reference target.
	if referencedField.Field == nil {
		return nil, fmt.Errorf(
			"field %q in %q references field %q in %q that does not have an associated column",
			field.Desc.Name(), field.Parent.Desc.FullName(),
			ref.FieldName, typ.Desc.FullName(),
		)
	}
	if len(referencedField.Indexes) == 0 || !referencedField.Indexes[0].IsUnique {
		return nil, fmt.Errorf(
			"field %q in %q references field %q in %q that is not unique",
			field.Desc.Name(), field.Parent.Desc.FullName(),
			ref.FieldName, typ.Desc.FullName(),
		)
	}
	fkName := fmt.Sprintf("%s_%s_fk", parentTblName, col.Name)
	referencedTbl, err := c.messageTableName(typ, true)
	if err != nil {
		return nil, err
	}
	return &types.ForeignKey{
		Name:      fkName,
		Fields:    []types.Field{col},
		RefTable:  referencedTbl,
		RefFields: []types.Field{*referencedField.Field},
	}, nil
}

// convertArrayFields converts fields with array types requiring a separate
// table for entries.
func (c Converter) convertArrayFields(field *protogen.Field, typ types.Type) (*ConvertedField, error) {
	parentSingular, err := c.messageTableName(field.Parent, false)
	if err != nil {
		return nil, err
	}
	parentTbl, err := c.messageTableName(field.Parent, true)
	if err != nil {
		return nil, err
	}
	fieldTypeSingular, err := c.messageTableName(field.Message, false)
	if err != nil {
		return nil, err
	}
	fieldTypeTbl, err := c.messageTableName(field.Message, true)
	if err != nil {
		return nil, err
	}
	// Construct name of table from the name of the reference table.
	fieldName := field.Desc.JSONName()
	refTableName := fmt.Sprintf("%s_%s_entries", parentSingular, fieldName)
	// Construct column names.
	leftCol := parentSingular + "_id"
	rightCol := fieldTypeSingular + "_id"
	// Construct the table.
	lookupTable := types.NewRefTable(refTableName, parentTbl, leftCol, fieldTypeTbl, rightCol)
	return &ConvertedField{
		ExtraTables: []types.Table{lookupTable},
	}, nil
}

// convertReferencedFields converts fields with message types referenced,
// storing the ID of the message type in a column allowing a JOIN to retrieve
// the information.
func (c Converter) convertReferencedFields(field *protogen.Field, typ types.Type) (*ConvertedField, error) {
	parentTblName, err := c.messageTableName(field.Parent, true)
	if err != nil {
		return nil, err
	}
	// Construct column.
	options := fieldOpts(field)
	if options.DefaultValue != "" {
		return nil, fmt.Errorf(
			"cannot set default value for type %q in %q",
			typ.Type, field.Desc.FullName(),
		)
	}
	if options.Ref != nil {
		return nil, fmt.Errorf(
			"cannot set references for type %q in %q",
			typ.Type, field.Desc.FullName(),
		)
	}

	// Generate information related to field type.
	fieldTypeTable, err := c.messageTableName(field.Message, true)
	if err != nil {
		return nil, err
	}
	// One-to-one relationship.
	col := &types.Field{
		Name: field.Desc.JSONName(),
		Type: types.Type{Type: "int32"},
	}
	fk := types.ForeignKey{
		Name:     parentTblName + "_" + col.Name + "_fkey",
		Fields:   []types.Field{*col},
		RefTable: fieldTypeTable,
		RefFields: []types.Field{
			{
				Name:       "id",
				Type:       types.Type{Type: "int32"},
				IsPrimary:  true,
				IsSequence: true,
			},
		},
	}
	return &ConvertedField{
		Field:       col,
		ForeignKeys: []types.ForeignKey{fk},
	}, nil
}
