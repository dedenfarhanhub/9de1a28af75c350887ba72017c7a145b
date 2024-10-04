package helpers

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConvertToString converts various unit types (int, uint, etc.) to string
func ConvertToString(input interface{}) (string, error) {
	switch v := input.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(input))
	}
}
