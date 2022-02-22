package util

import (
	"errors"
	"strconv"
)

var (
	ErrDuringConversion = errors.New("error during the conversion to string")
)

func ToString(value interface{}) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil

	case bool:
		return strconv.FormatBool(value), nil

	case float32, float64:
		return strconv.FormatFloat(value.(float64), 'g', -1, 64), nil

	case int, int16, int32, int64:
		return strconv.FormatInt(value.(int64), 10), nil

	case uint, uint16, uint32, uint64:
		return strconv.FormatUint(value.(uint64), 10), nil

	default:
		return "", ErrDuringConversion
	}
}
