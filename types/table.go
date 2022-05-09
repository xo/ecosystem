package types

import (
	"github.com/kenshaw/inflector"
	xo "github.com/xo/xo/types"
)

// convertTable converts the Table type to a fledged-out xo.Table.
// It automatically generates extra tables for columns with array types if the
// target is not postgres.
func (s sqlTypeConverter) convertTable(table Table) []xo.Table {
	xoTable := xo.Table{
		Name:        s.id(table.Name),
		Columns:     s.filterFields(table.Columns),
		ForeignKeys: make([]xo.ForeignKey, len(table.ForeignKeys)),
		Indexes:     make([]xo.Index, len(table.Indexes)),
	}
	for i, fk := range table.ForeignKeys {
		xoTable.ForeignKeys[i] = s.convertForeignKey(table, fk)
	}
	for i, idx := range table.Indexes {
		xoTable.Indexes[i] = s.convertIndex(table, idx)
		if idx.IsUnique {
			xoTable.Columns = s.filterUnique(xoTable.Columns, idx.Fields)
		}
	}
	tables := []xo.Table{xoTable}
	if s.Target != "postgres" {
		arrayTables := s.genArrayTables(s.id(table.Name), table.Columns)
		tables = append(tables, arrayTables...)
	}
	return tables
}

// genArrayTables creates lookup tables mapping from the parent's ID column to
// the content of the array. This is done for arrays of all simple types.
func (s sqlTypeConverter) genArrayTables(parentName string, f []Field) []xo.Table {
	tbl := make([]xo.Table, 0)
	for _, v := range f {
		if !v.Type.IsArray {
			continue
		}
		// Key referencing original table to join on.
		parentKeyField := Field{
			Name:      parentName + "_id",
			Type:      Type{Type: "int32"},
			IsPrimary: true,
		}
		// Values of the array.
		typ := v.Type
		typ.IsArray = false
		valueField := Field{
			Name: inflector.Singularize(v.Name),
			Type: typ,
		}
		lookupName := s.id(parentName + "_" + v.Name)
		newTbl := Table{
			Name:        lookupName,
			Columns:     []Field{parentKeyField, valueField},
			PrimaryKeys: []Field{parentKeyField},
			ForeignKeys: []ForeignKey{
				{
					Name:     s.id(lookupName + "_parent_fkey"),
					Fields:   []Field{parentKeyField},
					RefTable: parentName,
					RefFields: []Field{
						{
							Name:      "id",
							Type:      Type{Type: "int32"},
							IsPrimary: true,
						},
					},
				},
			},
			Indexes: []Index{
				{
					Name:   s.id(lookupName + "_lookup_key"),
					Fields: []Field{parentKeyField},
				},
			},
		}
		tbl = append(tbl, s.convertTable(newTbl)...)
	}
	return tbl
}

// filterFields converts the provided slice of fields to xo.Fields by mapping
// the types to the concrete SQL types.
// It also removes all fields with array types if present in non-Postgres
// databases.
func (s sqlTypeConverter) filterFields(f []Field) []xo.Field {
	converted := make([]xo.Field, 0, len(f))
	for _, v := range f {
		if v.Type.IsArray && s.Target != "postgres" {
			// Arrays of simple type in non-Postgres database.
			// Should be handled in genArrayTables.
			continue
		}
		field := xo.Field{
			Name:       v.Name,
			Type:       s.sqlType(v.Type),
			Default:    v.Default,
			IsPrimary:  v.IsPrimary,
			IsSequence: v.IsSequence,
		}
		converted = append(converted, field)
	}
	return converted
}

// filterUnique filters fields that have been unique to compatible types for
// sqlserver. It does nothing for other database targets.
func (s sqlTypeConverter) filterUnique(cols []xo.Field, indexes []Field) []xo.Field {
	if s.Target != "sqlserver" {
		return cols
	}
	for i, col := range cols {
		for _, idx := range indexes {
			if idx.Name != col.Name {
				continue
			}
			indexTyp, ok := indexTypes[col.Type.Type]
			if !ok {
				break
			}
			cols[i].Type = indexTyp
			break
		}
	}
	return cols
}
