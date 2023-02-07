package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/iancoleman/strcase"
)

func StrToMd5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func StrToUpperSnake(text string) string {
	snake := matchFirstCap.ReplaceAllString(text, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func StrToCamel(text string) string {
	return strcase.ToCamel(text)
}
