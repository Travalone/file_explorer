package utils

import (
	"file_explorer/common/logger"
	"fmt"
	"reflect"
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
		logger.Error("Conv2Str failed, data=%v, type=%s", data, reflect.TypeOf(data))
		return fmt.Sprintf("%v", data)
	}
}

func PtrInt(num int) *int {
	return &num
}
func PtrStr(str string) *string {
	return &str
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
