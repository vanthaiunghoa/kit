package generator

import (
	"os"
)

const ENV_NAME_KIT_PROTO_FILE_INDENT_CHAR = "KIT_PROTO_FILE_INDENT_CHAR"
const DEFAULT_KIT_PROTO_FILE_INDENT_CHAR = "    "

func getProtoFileIndentChar() string {
	if s := os.Getenv(ENV_NAME_KIT_PROTO_FILE_INDENT_CHAR); s != "" {
		return s
	}
	return DEFAULT_KIT_PROTO_FILE_INDENT_CHAR
}
