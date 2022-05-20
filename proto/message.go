package proto

import (
	"fmt"
	"strings"

	"github.com/kenshaw/inflector"
	"github.com/kenshaw/snaker"
	pb "github.com/xo/ecosystem/proto/xo"
	"github.com/xo/ecosystem/types"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// ConvertMessage converts the provided protogen message to xo tables.
// The first table always represents the actual message.
func (c Converter) ConvertMessage(pkg *protogen.File, msg *protogen.Message) ([]types.Table, error) {
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
		c.TableName(pkg, name, true),
		tableOpts.Manual,
	)

	var lookupTables []types.Table

	// Add HasMany entries.
	for _, hasMany := range tableOpts.HasMany {
		typ, err := c.searchForType(hasMany.TypeSuffix)
		if err != nil {
			return nil, err
		}
		rightTable, err := c.messageTableName(typ, true)
		if err != nil {
			return nil, err
		}
		rightSingular, err := c.messageTableName(typ, false)
		if err != nil {
			return nil, err
		}
		// Construct name of table from the name of the HasMany.
		leftSingular := c.TableName(pkg, name, false)
		hasManyName := snaker.CamelToSnake(hasMany.Name)
		hasManyTableName := fmt.Sprintf("%s_%s_entries", leftSingular, hasManyName)
		// Construct column names.
		leftCol := leftSingular + "_id"
		rightCol := rightSingular + "_id"
		// Construct the table.
		lookupTable := types.NewRefTable(hasManyTableName, table.Name, leftCol, rightTable, rightCol)
		lookupTables = append(lookupTables, lookupTable)
	}

	// Add message fields.
	for _, f := range msg.Fields {
		converted, err := c.ConvertField(f)
		if err != nil {
			return nil, err
		}
		if converted.Field != nil {
			table.Columns = append(table.Columns, *converted.Field)
		}
		table.Indexes = append(table.Indexes, converted.Indexes...)
		table.ForeignKeys = append(table.ForeignKeys, converted.ForeignKeys...)
		lookupTables = append(lookupTables, converted.ExtraTables...)
	}
	tables := make([]types.Table, 0, len(lookupTables)+1)
	tables = append(tables, table)
	tables = append(tables, lookupTables...)
	return tables, nil
}

// messageTableName is a helper function that returns the table name of the
// provided message.
func (c Converter) messageTableName(msg *protogen.Message, plural bool) (string, error) {
	pkg, err := c.pkgOf(msg.Desc)
	if err != nil {
		return "", err
	}
	return c.TableName(pkg, string(msg.Desc.Name()), plural), nil
}

// TableName returns the table name of the package and name pair.
func (c Converter) TableName(pkg *protogen.File, name string, plural bool) string {
	pkgSingular := inflector.Singularize(string(pkg.GoPackageName))
	pkgTitle := strings.Title(pkgSingular)

	// Prevent pkg_pkg table naming.
	opts := fileOpts(pkg)
	if opts.SkipPrefix || strings.HasPrefix(name, pkgTitle) {
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

// fileOpts returns the file options of the file or an empty FileOverride if
// the message is nil.
func fileOpts(pkg *protogen.File) *pb.FileOverride {
	if pkg == nil {
		return &pb.FileOverride{}
	}
	if proto.HasExtension(pkg.Desc.Options(), pb.E_FileOverrides) {
		return proto.GetExtension(pkg.Desc.Options(), pb.E_FileOverrides).(*pb.FileOverride)
	}
	return &pb.FileOverride{}
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
