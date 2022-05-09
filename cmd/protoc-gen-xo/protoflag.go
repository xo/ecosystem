package main

import (
	"fmt"

	"github.com/xo/ecosystem/proto"
	"github.com/xo/ecosystem/types"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type typeEntry struct {
	PkgName  string                `json:"pkg_name"`
	PkgPath  string                `json:"pkg_path"`
	TypeName string                `json:"type"`
	Fields   map[string]fieldEntry `json:"fields,omitempty"`
}

type fieldEntry struct {
	PkgName   string `json:"pkg_name,omitempty"`
	PkgPath   string `json:"pkg_path,omitempty"`
	Type      string `json:"type"`
	TableType string `json:"table_type,omitempty"`
	LinkTable string `json:"link_table,omitempty"`
	Kind      string `json:"kind"`
	Array     bool   `json:"repeated,omitempty"`
	Map       bool   `json:"map,omitempty"`
	Embedded  bool   `json:"embedded,omitempty"`
}

func (x *xoPlugin) addMessageType(c proto.Converter, msg *protogen.Message, tbl types.Table) {
	filePath := msg.Desc.ParentFile().Path()
	file := x.plugin.FilesByPath[filePath]
	pkgPath := string(file.GoImportPath)
	pkgName := string(file.GoPackageName)
	tblName := tbl.Name

	tblFields := make(map[string]bool, len(tbl.Columns))
	for _, col := range tbl.Columns {
		tblFields[col.Name] = true
	}

	// Add fields.
	entries := make(map[string]fieldEntry, len(msg.Fields))
	for _, v := range msg.Fields {
		fieldName := v.GoName
		var kind, fieldPkgName, fieldPkgPath, typ, tblType, linkTable string
		var embedded bool
		switch v.Desc.Kind() {
		case protoreflect.EnumKind:
			kind = "enum"
			parentPath := v.Enum.Desc.ParentFile().Path()
			parentFile := x.plugin.FilesByPath[parentPath]
			fieldPkgName = string(parentFile.GoPackageName)
			fieldPkgPath = string(parentFile.GoImportPath)
			typ = v.Enum.GoIdent.GoName
		case protoreflect.MessageKind:
			kind = "message"
			parentPath := v.Message.Desc.ParentFile().Path()
			parentFile := x.plugin.FilesByPath[parentPath]
			fieldPkgName = string(parentFile.GoPackageName)
			fieldPkgPath = string(parentFile.GoImportPath)
			typ = v.Message.GoIdent.GoName
			tblType = c.TableName(fieldPkgName, typ, true)
			embedded = tblFields[v.Desc.JSONName()]
			// Entries table name.
			msgName := string(msg.Desc.Name())
			msgPrefix := c.TableName(pkgName, msgName, false)
			jsonName := v.Desc.JSONName()
			linkTable = fmt.Sprintf("%s_%s_entries", msgPrefix, jsonName)
		default:
			kind = "basic"
			typ = v.Desc.Kind().String()
		}
		entries[fieldName] = fieldEntry{
			PkgName:   fieldPkgName,
			PkgPath:   fieldPkgPath,
			Type:      typ,
			TableType: tblType,
			LinkTable: linkTable,
			Embedded:  embedded,
			Kind:      kind,
			Array:     v.Desc.IsList(),
			Map:       v.Desc.IsMap(),
		}
	}

	x.protobufNames[tblName] = typeEntry{
		PkgName:  pkgName,
		PkgPath:  pkgPath,
		TypeName: msg.GoIdent.GoName,
		Fields:   entries,
	}
}

func (x *xoPlugin) addEnumType(enum *protogen.Enum, goEnum types.Enum) {
	filePath := enum.Desc.ParentFile().Path()
	file := x.plugin.FilesByPath[filePath]
	pkgName := string(file.GoPackageName)
	pkgPath := string(file.GoImportPath)
	tblName := goEnum.Name
	x.protobufNames[tblName] = typeEntry{
		PkgName:  pkgName,
		PkgPath:  pkgPath,
		TypeName: enum.GoIdent.GoName,
	}
}
