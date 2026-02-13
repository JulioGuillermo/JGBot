package args

import (
	"reflect"
)

func GetArgsMetaData(args any) *TypeMetaData {
	return inspectType(reflect.TypeOf(args))
}

func GetArgsMetaDataString(args any, isAdmin bool) string {
	meta := GetArgsMetaData(args)
	return meta.String(isAdmin)
}

func inspectType(typ reflect.Type) *TypeMetaData {
	switch typ.Kind() {
	case reflect.Pointer:
		return inspectType(typ.Elem())
	case reflect.Array, reflect.Slice:
		return inspectArray(typ)
	case reflect.Map:
		return inspectMap(typ)
	case reflect.Struct:
		return inspectStruct(typ)
	default:
		return inspectPrimitive(typ)
	}
}

func inspectPrimitive(typ reflect.Type) *TypeMetaData {
	var primitive string
	switch typ.Kind() {
	case reflect.Bool:
		primitive = "Bool"
	case reflect.String:
		primitive = "String"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		primitive = "IntNumber"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		primitive = "UIntNumber"
	case reflect.Float32, reflect.Float64:
		primitive = "FloatNumber"
	case reflect.Complex64, reflect.Complex128:
		primitive = "ComplexNumber"
	default:
		primitive = typ.Name()
	}
	return &TypeMetaData{
		Type: primitive,
	}
}

func inspectArray(typ reflect.Type) *TypeMetaData {
	return &TypeMetaData{
		Type:    "Array",
		SubType: inspectType(typ.Elem()),
	}
}

func inspectMap(typ reflect.Type) *TypeMetaData {
	keyT := inspectPrimitive(typ.Key()).Type
	return &TypeMetaData{
		Type:    "Map",
		KeyType: keyT,
		SubType: inspectType(typ.Elem()),
	}
}

func inspectStruct(typ reflect.Type) *TypeMetaData {
	meta := &TypeMetaData{
		Type:   "Struct",
		Fields: make([]*TypeMetaData, 0, typ.NumField()),
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldMeta := inspectType(field.Type)
		fieldMeta.Name = field.Name
		fieldMeta.Public = field.IsExported()
		fieldMeta.JsonName = field.Tag.Get("json")
		fieldMeta.Description = field.Tag.Get("description")
		fieldMeta.AdminOnly = field.Tag.Get("admin") == "true"

		meta.Fields = append(meta.Fields, fieldMeta)
	}
	return meta
}
