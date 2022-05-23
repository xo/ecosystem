package proto

import (
	"github.com/kenshaw/snaker"
	"github.com/xo/ecosystem/types"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var simpleType = map[protoreflect.Kind]string{
	protoreflect.DoubleKind:   "float64",
	protoreflect.FloatKind:    "float32",
	protoreflect.Int64Kind:    "int64",
	protoreflect.Uint64Kind:   "int64",
	protoreflect.Fixed64Kind:  "int64",
	protoreflect.Sfixed64Kind: "int64",
	protoreflect.Sint64Kind:   "int64",
	protoreflect.Int32Kind:    "int32",
	protoreflect.Uint32Kind:   "int32",
	protoreflect.Fixed32Kind:  "int32",
	protoreflect.Sfixed32Kind: "int32",
	protoreflect.Sint32Kind:   "int32",
	protoreflect.BoolKind:     "bool",
	protoreflect.StringKind:   "string",
	protoreflect.BytesKind:    "[]byte",
}

var wellKnownType = map[protoreflect.FullName]string{
	"google.protobuf.Timestamp": "time",
	"google.protobuf.Duration":  "duration",
	"google.protobuf.Value":     "[]byte",
}

// goType converts a proto type to a types.Type.
func (c Converter) goType(field *protogen.Field) (typ types.Type, simple bool, err error) {
	simple = true
	typ.IsArray = field.Desc.IsList()
	typ.Nullable = fieldOpts(field).Nullable
	if field.Desc.IsMap() {
		// Protobuf map types will be json encoded when stored.
		typ := types.Type{
			Type: "[]byte",
		}
		return typ, true, nil
	}

	// Handle simple types.
	ftype := field.Desc.Kind()
	if simple, ok := simpleType[ftype]; ok {
		typ.Type = simple
		return typ, true, nil
	}

	switch ftype {
	case protoreflect.MessageKind:
		fullName := field.Desc.Message().FullName()
		if typeName, ok := wellKnownType[fullName]; ok {
			typ.Type = typeName
			return typ, true, nil
		}
		// TODO: support gunk field tags.
		typ.Type = string(fullName)
		return typ, false, nil

	case protoreflect.EnumKind:
		pkg, err := c.pkgOf(field.Enum.Desc)
		if err != nil {
			return typ, false, err
		}
		enumName := snaker.CamelToSnake(string(field.Enum.Desc.Name()))
		typ.Type = "enum"
		typ.EnumName = string(pkg.GoPackageName) + "_" + enumName
		if fileOpts(pkg).SkipPrefix {
			typ.EnumName = enumName
		}
		return typ, true, nil

	default:
		panic("unknown field type: " + ftype.String())
	}
}
