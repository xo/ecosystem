package proto

import (
	"github.com/kenshaw/snaker"
	"github.com/xo/ecosystem/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// goType converts a proto type to a types.Type.
func (c Converter) goType(field *protogen.Field) (typ types.Type, simple bool) {
	simple = true
	typ.IsArray = field.Desc.IsList()
	if field.Desc.IsMap() {
		// Protobuf map types will be json encoded when stored.
		return types.Type{
			Type: "[]byte",
		}, true
	}

	switch ftype := field.Desc.Kind(); ftype {
	case protoreflect.MessageKind:
		switch field.Desc.Message().FullName() {
		case "google.protobuf.Timestamp":
			typ.Type = "time"
		case "google.protobuf.Duration":
			typ.Type = "duration"
		default:
			// TODO: support gunk field tags.
			typ.Type = string(field.Desc.Message().FullName())
			simple = false
		}

	case protoreflect.DoubleKind:
		typ.Type = "float64"

	case protoreflect.FloatKind:
		typ.Type = "float32"

	case protoreflect.Int64Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Sint64Kind:
		typ.Type = "int64"

	case protoreflect.Int32Kind, protoreflect.Uint32Kind,
		protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Sint32Kind:
		typ.Type = "int32"

	case protoreflect.BoolKind:
		typ.Type = "bool"

	case protoreflect.StringKind:
		typ.Type = "string"

	case protoreflect.BytesKind:
		typ.Type = "[]byte"

	case protoreflect.EnumKind:
		pkgName := c.PackageNames[field.Enum.Desc.ParentFile().Path()]
		enumName := snaker.CamelToSnake(string(field.Enum.Desc.Name()))
		typ.Type = "enum"
		typ.EnumName = pkgName + "_" + enumName
		if slices.Contains(c.SkipPrefixes, pkgName) {
			typ.EnumName = enumName
		}

	default:
		panic("unknown field type: " + ftype.String())
	}
	return
}
