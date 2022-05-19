package main

import (
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
	PkgName string `json:"pkg_name,omitempty"`
	PkgPath string `json:"pkg_path,omitempty"`
	Type    string `json:"type"`
	Kind    string `json:"kind"`
	Array   bool   `json:"array,omitempty"`
	Map     bool   `json:"map,omitempty"`
}

func (x *xoPlugin) addMessageType(c proto.Converter, msg *protogen.Message, tbl types.Table) {
	tblFields := make(map[string]bool, len(tbl.Columns))
	for _, col := range tbl.Columns {
		tblFields[col.Name] = true
	}

	// Add fields.
	entries := make(map[string]fieldEntry, len(msg.Fields))
	for _, v := range msg.Fields {
		fieldName := v.GoName
		var kind, fieldPkgName, fieldPkgPath, typ string
		switch v.Desc.Kind() {
		case protoreflect.EnumKind:
			kind = "enum"
			fieldPkgName, fieldPkgPath = x.goPkgName(v.Enum.Desc)
			typ = v.Enum.GoIdent.GoName
		case protoreflect.MessageKind:
			kind = "message"
			fieldPkgName, fieldPkgPath = x.goPkgName(v.Message.Desc)
			typ = v.Message.GoIdent.GoName
		default:
			kind = "basic"
			typ = v.Desc.Kind().String()
		}
		entries[fieldName] = fieldEntry{
			PkgName: fieldPkgName,
			PkgPath: fieldPkgPath,
			Type:    typ,
			Kind:    kind,
			Array:   v.Desc.IsList(),
			Map:     v.Desc.IsMap(),
		}
	}

	msgPkgName, msgPkgPath := x.goPkgName(msg.Desc)
	tblName := tbl.Name
	x.protobufNames[tblName] = typeEntry{
		PkgName:  msgPkgName,
		PkgPath:  msgPkgPath,
		TypeName: msg.GoIdent.GoName,
		Fields:   entries,
	}
}

func (x *xoPlugin) addEnumType(enum *protogen.Enum, goEnum types.Enum) {
	pkgName, pkgPath := x.goPkgName(enum.Desc)
	tblName := goEnum.Name
	x.protobufNames[tblName] = typeEntry{
		PkgName:  pkgName,
		PkgPath:  pkgPath,
		TypeName: enum.GoIdent.GoName,
	}
}

func (x *xoPlugin) parentFile(d protoreflect.Descriptor) *protogen.File {
	return x.plugin.FilesByPath[d.ParentFile().Path()]
}

// goPkgName returns the package name and the import path of the package
// containing the object declared by the descriptor.
func (x *xoPlugin) goPkgName(d protoreflect.Descriptor) (string, string) {
	f := x.parentFile(d)
	return string(f.GoPackageName), string(f.GoImportPath)
}
