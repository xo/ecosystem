package types

import (
	"fmt"

	xo "github.com/xo/xo/types"
	"golang.org/x/exp/maps"
)

// sqlTypeConverter is an internal struct that contains data related to current
// conversion.
type sqlTypeConverter struct {
	// Enums is a map of the enums' name to the underlying enum.
	Enums map[string]*xo.Enum
	// Target is the database driver the converter is targeting.
	Target string

	// idUsedCache stores the number of types the identifier has been used.
	idUsedCache map[string]int
	// idMappedCache stores the current mapped identifier.
	idMappedCache map[string]string
}

// sqlType converts the provided Type to a xo.Type. It panics if provided an
// array type on non-Postgres databases.
func (s sqlTypeConverter) sqlType(goType Type) xo.Type {
	if goType.IsArray {
		if s.Target != "postgres" {
			panic("unsupported database type for arrays: " + s.Target)
		}
		// Postgres supports arrays natively so re-run the code and treat as
		// usual singular types except setting the array flag after.
		goType.IsArray = false
		singleTyp := s.sqlType(goType)
		singleTyp.IsArray = true
		return singleTyp
	}
	if s.Target == "sqlserver" && goType.Type == "string" {
		// Check unique.
	}
	t, ok := typeMap[s.Target][goType.Type]
	if ok {
		return t
	}
	if goType.Type != "enum" {
		panic("unexpected go type: " + goType.Type)
	}
	return s.buildEnum(goType)
}

// buildEnum creates an enum type by assigning an appropriate type for the
// database and setting the Enum field. It panics if the Type provided is of
// unknown enum type.
func (s sqlTypeConverter) buildEnum(goType Type) xo.Type {
	enumName := s.id(goType.EnumName)
	var typ xo.Type
	switch s.Target {
	case "postgres", "mysql":
		typ = buildType(enumName)
	case "sqlite3", "sqlserver":
		typ = buildType("int") // enum index
	default:
		panic("unsupported database type for enum creation: " + s.Target)
	}
	var ok bool
	typ.Enum, ok = s.Enums[enumName]
	if !ok {
		panic(fmt.Sprintf(
			"missing enum type for %q\n\tEnums present: %v",
			enumName, maps.Keys(s.Enums),
		))
	}
	return typ
}

// To use instead of true and false to aid with readability.
var (
	single = false
	array  = true
)

var typeMap = map[string]map[string]xo.Type{
	"postgres": {
		"int64":    buildType("bigint"),
		"int32":    buildType("integer"),
		"float64":  buildType("double precision"),
		"float32":  buildType("real"),
		"bool":     buildType("boolean"),
		"string":   buildType("text"),
		"time":     buildType("timestamp with time zone"),
		"duration": buildType("bigint"),
		"[]byte":   buildType("json"),
	},
	"mysql": {
		"int64":    buildType("bigint"),
		"int32":    buildType("int"),
		"float64":  buildType("double"),
		"float32":  buildType("real"),
		"bool":     buildType("boolean"),
		"string":   buildType("text"),
		"time":     buildType("timestamp"),
		"duration": buildType("bigint"),
		"[]byte":   buildType("json"),
	},
	"sqlite3": {
		"int64":    buildType("bigint"),
		"int32":    buildType("integer"),
		"float64":  buildType("double"),
		"float32":  buildType("real"),
		"bool":     buildType("boolean"),
		"string":   buildType("text"),
		"time":     buildType("datetime"),
		"duration": buildType("bigint"),
		"[]byte":   buildType("blob"),
	},
	"sqlserver": {
		"int64":    buildType("bigint"),
		"int32":    buildType("int"),
		"float64":  buildType("decimal"),
		"float32":  buildType("real"),
		"bool":     buildType("tinyint"),
		"string":   buildType("text"),
		"time":     buildType("datetime2"),
		"duration": buildType("bigint"),
		"[]byte":   buildType("binary"),
	},
}

var indexTypes = map[string]xo.Type{
	"text": buildType("varchar", 255),
}

// buildType is a helper to create a type with the provided type string.
func buildType(typ string, opts ...int) xo.Type {
	t := xo.Type{
		Type: typ,
	}
	switch len(opts) {
	case 2:
		t.Scale = opts[1]
		fallthrough
	case 1:
		t.Prec = opts[0]
	}
	return t
}
