package util

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/go-openapi/inflect"
	"github.com/iancoleman/strcase"
)

var rules *inflect.Ruleset

func init() {
	rules = ruleset()
}
func ruleset() *inflect.Ruleset {
	rules := inflect.NewDefaultRuleset()
	//acronyms := make(map[string]struct{})
	//// Add common initialism from golint and more.
	//for _, w := range []string{
	//	"ACL", "API", "ASCII", "AWS", "CPU", "CSS", "DNS", "EOF", "GB", "GUID",
	//	"HCL", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "KB", "LHS", "MAC",
	//	"MB", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "SSO",
	//	"TCP", "TLS", "TTL", "UDP", "UI", "UID", "URI", "URL", "UTF8", "UUID",
	//	"VM", "XML", "XMPP", "XSRF", "XSS",
	//} {
	//	acronyms[w] = struct{}{}
	//	rules.AddAcronym(w)
	//}
	return rules
}

// UpperToLowerCamel 大驼峰转小驼峰
func UpperToLowerCamel(s string) string {
	if len(s) == 0 {
		return s
	}

	str := bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		if r := rune(s[i]); r >= 'A' && r <= 'Z' {
			str.WriteString(strings.ToLower(string(r)))
		} else {
			str.WriteString(s[i:])
			break
		}
	}
	return str.String()
}

// LowerToUpperCamel 小驼峰转大驼峰
func LowerToUpperCamel(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// UnderlineStr2Strikethrough 将下划线字符串转换为中划线
func UnderlineStr2Strikethrough(str string) string {
	return strings.ReplaceAll(str, "_", "-")
}

// ReplaceTime2TimesISOTime 替换时间time.Time -> times.ISOTime
func ReplaceTime2TimesISOTime(str string) string {
	return strings.ReplaceAll(str, "time.Time", "times.ISOTime")
}

func ToLowerSnakeCase(str string) string {
	return strcase.ToSnake(str)
}

func ToUpperCamelCase(str string) string {
	return strcase.ToCamel(str)
}

func ArrayToPlural(str string) (string, bool) {
	slog.Info("ArrayToPlural", "str", str)
	switch {
	case strings.HasPrefix(str, "[]"):
		str = strings.TrimPrefix(str, "[]")
		return str, true
	case strings.HasPrefix(str, "*"):
		str = strings.TrimPrefix(str, "*")
		return str, false
	default:

	}
	return str, false
}

func ToPlural(name string) string {
	return rules.Pluralize(name)
}
