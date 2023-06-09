package utils

import (
	"file_explorer/common/logger"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func Conv2Str(data interface{}) string {
	if data == nil {
		return ""
	}

	switch data.(type) {
	case int:
		return fmt.Sprintf("%d", data)
	case int64:
		return fmt.Sprintf("%d", data)
	case int32:
		return fmt.Sprintf("%d", data)
	case byte:
		return fmt.Sprintf("%d", data)
	case float32:
		return fmt.Sprintf("%f", data)
	case float64:
		return fmt.Sprintf("%f", data)
	case []string:
		str := ""
		for _, item := range data.([]string) {
			if len(str) == 0 {
				str += item
			} else {
				str += ", " + item
			}
		}
		return str
	default:
		if dataStr, ok := data.(string); ok {
			return dataStr
		}
		logger.Error("Conv2Str failed, data=%v, type=%v", data, reflect.TypeOf(data))
		return fmt.Sprintf("%v", data)
	}
}

func PtrInt(num int) *int {
	return &num
}
func PtrStr(str string) *string {
	return &str
}

func ReplaceWrap(str string) string {
	return strings.ReplaceAll(str, "\n", "\\n")
}

func ReplaceWildcard(str string, rule map[string]string) string {
	pattern := regexp.MustCompile("{\\w+}")
	matchList := pattern.FindAllString(str, -1)
	for _, match := range matchList {
		key := match[1 : len(match)-1]
		if value, ok := rule[key]; ok {
			str = strings.ReplaceAll(str, match, value)
		}
	}
	return str
}

var SizeUnits = []string{"B", "K", "M", "G"}

func ConvSize(sizeB int64) string {
	var size = float64(sizeB)
	count := 0
	for size > 1024 && count < len(SizeUnits) {
		size /= 1024
		count++
	}
	return fmt.Sprintf("%.2f%s", size, SizeUnits[count])
}

func Strings2Interfaces(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, str := range strings {
		interfaces[i] = str
	}
	return interfaces
}
func Interfaces2Strings(interfaces []interface{}) []string {
	strings := make([]string, len(interfaces))
	for i, item := range interfaces {
		strings[i] = item.(string)
	}
	return strings
}

func hasCh(s string) bool {
	for _, ch := range s {
		if unicode.Is(unicode.Han, ch) {
			return true
		}
	}
	return false
}

func ch2Py(s string) string {
	if !hasCh(s) {
		return s
	}
	return strings.Join(pinyin.LazyConvert(s, nil), " ")
}

func CmpText(s1, s2 string) bool {
	a := ch2Py(strings.ToLower(s1))
	b := ch2Py(strings.ToLower(s2))
	return a < b
}
