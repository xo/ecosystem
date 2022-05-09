package proto

import (
	"fmt"
	"strings"

	"github.com/kenshaw/inflector"
	"github.com/kenshaw/snaker"
	pb "github.com/xo/ecosystem/proto/xo"
	"github.com/xo/ecosystem/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// ConvertMessage converts the provided protogen message to xo tables.
// The first table always represents the actual message.
func (c Converter) ConvertMessage(pkgName string, msg *protogen.Message) ([]types.Table, error) {
	// Skip messages with Request or Response as suffix.
	name := string(msg.Desc.Name())
	if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") {
		// TODO: add case to check for ignored message entry types
		return nil, nil
	}

	tableOpts := messageOpts(msg)
	if tableOpts.Ignore {
		return nil, nil
	}

	table := types.NewTable(
		c.TableName(pkgName, name, true),
		tableOpts.Manual,
	)

	var lookupTables []types.Table

	// Add message fields.
	for _, field := range msg.Fields {
		options := fieldOpts(field)
		if options.Ignore {
			continue
		}
		msgOpts := messageOpts(field.Message)
		if options.EmbedAsJson || msgOpts.Ignore || msgOpts.EmbedAsJson {
			// Embed fields ignored reftypes and those explicitly marked as JSON.
			col := types.Field{
				Name: field.Desc.JSONName(),
				Type: types.Type{Type: "[]byte"},
			}
			table.Columns = append(table.Columns, col)
			continue
		}
		typ, resolved := c.goType(field)
		if resolved {
			col := types.Field{
				Name:    field.Desc.JSONName(),
				Type:    typ,
				Default: options.DefaultValue,
			}
			table.Columns = append(table.Columns, col)
			// Add unique index if marked in options.
			idx := options.Index
			if idx != pb.FieldOverride_NONE {
				table.Indexes = append(table.Indexes, types.Index{
					Name:     table.Name + "_" + col.Name + "_idx",
					Fields:   []types.Field{col},
					IsUnique: idx == pb.FieldOverride_UNIQUE,
				})
			}
			continue
		}

		if options.DefaultValue != "" {
			return nil, fmt.Errorf(
				"cannot set default value for type %q in %q",
				typ.Type, field.Desc.FullName(),
			)
		}

		// If resolved is false, the field is not a simple type, and references
		// a different table.
		path := field.Message.Desc.ParentFile().Path()
		refType := string(field.Message.Desc.Name())
		rightTbl := c.TableName(c.PackageNames[path], refType, true)
		rightTblSingular := c.TableName(c.PackageNames[path], refType, false)
		if typ.IsArray {
			// One-to-many relationship.
			leftTblSingular := c.TableName(pkgName, name, false)
			lookupTable := types.NewRefTable(
				leftTblSingular+"_"+field.Desc.JSONName()+"_entries",
				table.Name, leftTblSingular+"_id",
				rightTbl, rightTblSingular+"_id",
			)
			lookupTables = append(lookupTables, lookupTable)
			continue
		}
		// One-to-one relationship.
		field := types.Field{
			Name: field.Desc.JSONName(),
			Type: types.Type{Type: "int32"},
		}
		// if field.Name == rightTblSingular {
		// 	// If the name matches, add the ID suffix to the name.
		// 	field.Name += "_id"
		// }
		table.Columns = append(table.Columns, field)
		table.ForeignKeys = append(table.ForeignKeys, types.ForeignKey{
			Name:     table.Name + "_" + field.Name + "_fkey",
			Fields:   []types.Field{field},
			RefTable: rightTbl,
			RefFields: []types.Field{
				{
					Name:       "id",
					Type:       types.Type{Type: "int32"},
					IsPrimary:  true,
					IsSequence: true,
				},
			},
		})
	}
	tables := make([]types.Table, 0, len(lookupTables)+1)
	tables = append(tables, table)
	tables = append(tables, lookupTables...)
	return tables, nil
}

// TableName returns the table name of the package and name pair.
func (c Converter) TableName(pkg string, name string, plural bool) string {
	pkgSingular := inflector.Singularize(string(pkg))
	pkgTitle := strings.Title(pkgSingular)

	// Prevent pkg_pkg table naming.
	if slices.Contains(c.SkipPrefixes, pkg) || strings.HasPrefix(name, pkgTitle) {
		snake := snaker.CamelToSnake(name)
		if !plural {
			return snake
		}
		return inflector.Pluralize(snake)
	}

	suffix := strings.TrimPrefix(name, pkgTitle)
	if plural {
		suffix = inflector.Pluralize(suffix)
	}
	return pkgSingular + "_" + snaker.CamelToSnake(suffix)
}

// messageOpts returns the message options of the message or an empty
// MessageOverride if the message is nil.
func messageOpts(msg *protogen.Message) *pb.MessageOverride {
	if msg == nil {
		return &pb.MessageOverride{}
	}
	if proto.HasExtension(msg.Desc.Options(), pb.E_MsgOverrides) {
		return proto.GetExtension(msg.Desc.Options(), pb.E_MsgOverrides).(*pb.MessageOverride)
	}
	return &pb.MessageOverride{}
}

// fieldOpts returns the field options of the field or an empty FieldOverride
// if the field is nil.
func fieldOpts(field *protogen.Field) *pb.FieldOverride {
	if field == nil {
		return &pb.FieldOverride{}
	}
	if proto.HasExtension(field.Desc.Options(), pb.E_FieldOverrides) {
		return proto.GetExtension(field.Desc.Options(), pb.E_FieldOverrides).(*pb.FieldOverride)
	}
	return &pb.FieldOverride{}
}
