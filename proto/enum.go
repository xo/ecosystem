package proto

import (
	"strconv"

	"github.com/kenshaw/snaker"
	"github.com/xo/ecosystem/types"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
)

// ConvertEnum converts the specified protobuf enum to a types enum.
func (c Converter) ConvertEnum(protoEnum *protogen.Enum) (types.Enum, error) {
	pkgName := c.PackageNames[protoEnum.Desc.ParentFile().Path()]
	enumTypeName := snaker.CamelToSnake(string(protoEnum.Desc.Name()))
	e := types.Enum{
		Name: pkgName + "_" + enumTypeName,
	}
	if slices.Contains(c.SkipPrefixes, pkgName) {
		e.Name = enumTypeName
	}
	// Create slice of values.
	e.Values = make([]string, 0, len(protoEnum.Values))
	for _, v := range protoEnum.Values {
		name := string(v.Desc.Name())
		num := int(v.Desc.Number())
		for len(e.Values) <= num {
			// Fill in with ENUM_NAME_NUMBER for blank spots to preserve
			// numbering of enums.
			nextNumber := strconv.Itoa(len(e.Values))
			e.Values = append(e.Values, enumTypeName+"_"+nextNumber)
		}
		e.Values[v.Desc.Number()] = name
	}
	return e, nil
}
