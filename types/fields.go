package types

import (
	"strings"

	"github.com/kenshaw/inflector"
	xo "github.com/xo/xo/types"
)

// convertEnum converts the provided Enum to a xo.Enum and records it in
// sqlTypeConverter's internal map.
func (s sqlTypeConverter) convertEnum(target string, enum Enum) xo.Enum {
	e := xo.Enum{
		Name: s.id(enum.Name),
	}
	e.Values = make([]xo.Field, 0, len(enum.Values))
	for i, v := range enum.Values {
		// Enum const values are 1-indexed.
		constVal := i + 1
		e.Values = append(e.Values, xo.Field{
			Name:       v,
			ConstValue: &constVal,
		})
	}
	s.Enums[e.Name] = &e
	return e
}

// convertIndex converts the provided Index into a xo.Index.
func (s sqlTypeConverter) convertIndex(t Table, i Index) xo.Index {
	fn := indexFuncName(t.Name, i.Fields)
	return xo.Index{
		Name:      s.id(i.Name),
		Fields:    s.filterFields(i.Fields),
		IsUnique:  i.IsUnique,
		IsPrimary: i.IsPrimary,
		Func:      fn,
	}
}

// convertForeignKey converts the provided ForeignKey to a xo.ForeignKey.
func (s sqlTypeConverter) convertForeignKey(t Table, fk ForeignKey) xo.ForeignKey {
	fn := foreignKeyFuncName(t, fk)
	refFn := indexFuncName(fk.RefTable, fk.RefFields)
	return xo.ForeignKey{
		Name:      s.id(fk.Name),
		Fields:    s.filterFields(fk.Fields),
		RefTable:  s.id(fk.RefTable),
		RefFields: s.filterFields(fk.RefFields),
		Func:      fn,
		RefFunc:   refFn,
	}
}

// indexFuncName creates the name of an index function.
func indexFuncName(tableName string, indexFields []Field) string {
	tableName = inflector.Singularize(tableName)
	fieldNames := make([]string, 0, len(indexFields))
	for _, f := range indexFields {
		fieldNames = append(fieldNames, f.Name)
	}
	return tableName + "_by_" + strings.Join(fieldNames, "_")
}

// foreignKeyFuncName creates the name of a foreign key function.
func foreignKeyFuncName(t Table, fk ForeignKey) string {
	// Generate qualified name.
	tableName := inflector.Singularize(fk.RefTable)
	fieldNames := make([]string, 0, len(fk.Fields))
	for _, f := range fk.Fields {
		fieldNames = append(fieldNames, f.Name)
	}
	qualifiedName := tableName + "_by_" + strings.Join(fieldNames, "_")
	// If there are conflicts with existing field names, use qualified name.
	for _, v := range t.Columns {
		if tableName == v.Name {
			return qualifiedName
		}
	}
	// If there are conflicts in RefTable name, make it unique by using field
	// names.
	for _, v := range t.ForeignKeys {
		if fk.Name != v.Name && fk.RefTable == v.RefTable {
			return qualifiedName
		}
	}
	// Otherwise, just use RefTable name.
	return tableName
}
