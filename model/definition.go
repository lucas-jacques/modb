package model

import (
	"reflect"
)

type ModelDefinition[Model any, Schema any, PK comparable] struct {
	Table  string
	Schema Schema
	PK     func(Schema) TypedCol[Model, PK]
}

func Define[M any, PK comparable, S any](def ModelDefinition[M, S, PK]) TypedModelCols[M, PK, S] {
	m := &model[M, PK, S]{
		table:      def.Table,
		primaryKey: def.PK(def.Schema),
		schema:     def.Schema,
	}

	schemaType := reflect.TypeOf(def.Schema)
	schemaValue := reflect.ValueOf(def.Schema)

	if schemaType.Kind() != reflect.Struct {
		panic("Schema must be a struct")
	}

	columns := make([]Column, schemaType.NumField())
	for i := 0; i < schemaType.NumField(); i++ {
		field := schemaType.Field(i)
		// check if the field implement Column
		if !field.Type.Implements(reflect.TypeFor[ModelCol[M]]()) {
			continue
		}
		column := schemaValue.Field(i).Interface().(ModelCol[M])
		column.SetInnerModel(m)
		column.setTable(def.Table)
		columns[i] = column

	}

	m.columns = columns

	return m
}
