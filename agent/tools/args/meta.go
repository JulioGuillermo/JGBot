package args

import (
	"fmt"
	"strings"
)

type TypeMetaData struct {
	Name        string
	Public      bool
	JsonName    string
	Description string

	Type    string
	SubType *TypeMetaData
	KeyType string

	AdminOnly bool
	Fields    []*TypeMetaData
}

func (m *TypeMetaData) IsArray() bool {
	return m.Type == "Array"
}

func (m *TypeMetaData) IsMap() bool {
	return m.Type == "Map"
}

func (m *TypeMetaData) IsStruct() bool {
	return m.Type == "Struct"
}

func (m *TypeMetaData) IsPrimitive() bool {
	return !(m.IsArray() || m.IsMap() || m.IsStruct())
}

func (m *TypeMetaData) FieldName() string {
	if m.JsonName != "" {
		return m.JsonName
	}
	return m.Name
}

func (m *TypeMetaData) String(isAdmin bool) string {
	return m.getString("", isAdmin)
}

func (m *TypeMetaData) getFieldNameString() string {
	fn := m.FieldName()
	if fn == "" {
		return ""
	}
	return fmt.Sprintf("%s: ", fn)
}
func (m *TypeMetaData) getFieldDescriptionString() string {
	if m.Description == "" {
		return ""
	}
	return fmt.Sprintf(" /* %s */", m.Description)
}

func (m *TypeMetaData) getString(prefix string, isAdmin bool) string {
	if m.AdminOnly && !isAdmin {
		return ""
	}

	if m.IsArray() {
		return m.arrayString(prefix, isAdmin)
	}

	if m.IsMap() {
		return m.mapString(prefix, isAdmin)
	}

	if m.IsStruct() {
		return m.structString(prefix, isAdmin)
	}

	return m.primitiveString(prefix)
}

func (m *TypeMetaData) primitiveString(prefix string) string {
	fieldName := m.getFieldNameString()
	description := m.getFieldDescriptionString()
	return fmt.Sprintf("%s- %s%s%s", prefix, fieldName, m.Type, description)
}

func (m *TypeMetaData) arrayString(prefix string, isAdmin bool) string {
	fieldName := m.getFieldNameString()
	description := m.getFieldDescriptionString()
	subtypeString := m.SubType.getString(prefix+"  ", isAdmin)
	return fmt.Sprintf("%s- %sArray [%s\n%s\n]", prefix, fieldName, description, subtypeString)
}

func (m *TypeMetaData) mapString(prefix string, isAdmin bool) string {
	fieldName := m.getFieldNameString()
	description := m.getFieldDescriptionString()
	subtypeString := m.SubType.getString(prefix+"  ", isAdmin)
	return fmt.Sprintf("%s- %sMap [%s ->%s\n%s\n]", prefix, fieldName, m.KeyType, description, subtypeString)
}

func (m *TypeMetaData) getFields(isAdmin bool) []*TypeMetaData {
	var fields []*TypeMetaData
	for _, field := range m.Fields {
		if field.AdminOnly && !isAdmin {
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

func (m *TypeMetaData) structString(prefix string, isAdmin bool) string {
	fieldName := m.getFieldNameString()
	description := m.getFieldDescriptionString()
	fields := m.getFields(isAdmin)

	var sb strings.Builder
	if len(fields) > 0 {
		sb.WriteString("\n")
	} else {
		sb.WriteString(" ")
	}
	for _, field := range fields {
		s := field.getString(prefix+"  ", isAdmin)
		if s == "" {
			continue
		}
		sb.WriteString(s)
		sb.WriteString(",\n")
	}

	return fmt.Sprintf("%s- %s%s {%s%s}", prefix, fieldName, m.Type, description, sb.String())
}
