package templ

import (
	"strings"
)

const NativeToolDescriptionTemplate = `TOOL NAME: {{Name}}

DESCRIPTION:
{{Description}}

ARGUMENTS (JSON Schema):
{{Args}}

USAGE GUIDELINES:
1. Provide input strictly in JSON format.
2. Ensure types match the schema above.
3. If an error occurs, the tool will return a "TOOL ERROR" prefix.`

func GetNativeToolDescription(name, description, args string) string {
	str := NativeToolDescriptionTemplate
	str = strings.ReplaceAll(str, "{{Name}}", name)
	str = strings.ReplaceAll(str, "{{Description}}", description)
	str = strings.ReplaceAll(str, "{{Args}}", args)
	return str
}
