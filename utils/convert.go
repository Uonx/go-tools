package utils

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
)

type convert struct {
}

func Convert() convert {
	return convert{}
}

func (c convert) Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (c convert) Json(a any) string {
	str, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(str)
}

func (c convert) String(v string) *string {
	return &v
}

func (c convert) Bool(v string) bool {
	s := strings.ToLower(v)
	switch s {
	case "true", "1":
		return true
	default:
		return false
	}
}

func (c convert) Unicode(v string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(v), `\\u`, `\u`, -1))
	if err != nil {
		v = ""
	}
	v = str
	return v
}

func (c convert) Base64(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func (c convert) Base64ToString(v string) string {
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return ""
	}
	return string(b)
}
