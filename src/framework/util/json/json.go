package json

import (
	"fmt"
	"strconv"
	"strings"
)

func ToJsonString(data interface{}) string {
	return valueToString(data)
}

func valueToString(v interface{}) string {
	switch v.(type) {
	// number
	case int8:
		return strconv.Itoa(int(v.(int8)))
	case int16:
		return strconv.Itoa(int(v.(int16)))
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case uint8:
		return strconv.Itoa(int(v.(uint8)))
	case uint16:
		return strconv.Itoa(int(v.(uint16)))
	case uint:
		return strconv.Itoa(int(v.(uint)))
	case uint32:
		return strconv.Itoa(int(v.(uint32)))
	case uint64:
		return fmt.Sprintf("%lld", v.(uint64))
	case float32:
		return fmt.Sprintf("%.4f", v.(float32))
	case float64:
		return fmt.Sprintf("%.8f", v.(float64))
	case string:
		return fmt.Sprintf(`"%s"`, v.(string))
	//list
	case []interface{}:
		return interfaceListToString(v.([]interface{}))
	case []int8:
		return int8ListToString(v.([]int8))
	case []int16:
		return int16ListToString(v.([]int16))
	case []int:
		return intListToString(v.([]int))
	case []int32:
		return int32ListToString(v.([]int32))
	case []int64:
		return int64ListToString(v.([]int64))
	case []uint8:
		return uint8ListToString(v.([]uint8))
	case []uint16:
		return uint16ListToString(v.([]uint16))
	case []uint32:
		return uint32ListToString(v.([]uint32))
	case []uint64:
		return uint64ListToString(v.([]uint64))
	case []float32:
		return float32ListToString(v.([]float32))
	case []float64:
		return float64ListToString(v.([]float64))
	case []string:
		return stringListToString(v.([]string))
	// map
	case map[string]int8:
		return int8MapToString(v.(map[string]int8))
	case map[string]int16:
		return int16MapToString(v.(map[string]int16))
	case map[string]int:
		return intMapToString(v.(map[string]int))
	case map[string]int32:
		return int32MapToString(v.(map[string]int32))
	case map[string]int64:
		return int64MapToString(v.(map[string]int64))
	case map[string]uint8:
		return uint8MapToString(v.(map[string]uint8))
	case map[string]uint16:
		return uint16MapToString(v.(map[string]uint16))
	case map[string]uint32:
		return uint32MapToString(v.(map[string]uint32))
	case map[string]uint64:
		return uint64MapToString(v.(map[string]uint64))
	case map[string]float32:
		return float32MapToString(v.(map[string]float32))
	case map[string]float64:
		return float64MapToString(v.(map[string]float64))
	case map[string]string:
		return stringMapToString(v.(map[string]string))
	case map[string]interface{}:
		return interfaceMapToString(v.(map[string]interface{}))
	default:
		return ""
	}
}

func stringListToString(l []string) string {
	ret := "["
	for i, v := range l {
		ret += fmt.Sprintf(`"%s"`, v)
		if i != len(l)-1 {
			ret += ","
		}
	}
	ret += "]"
	return ret
}

func integerListToString(l []string) string {
	ret := "["
	ret += strings.Join(l, ",")
	ret += "]"
	return ret
}

func byteListToString(l []byte) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func int8ListToString(l []int8) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func int16ListToString(l []int16) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func intListToString(l []int) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func int32ListToString(l []int32) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func int64ListToString(l []int64) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func uint8ListToString(l []uint8) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func uint16ListToString(l []uint16) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func uint32ListToString(l []uint32) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func uint64ListToString(l []uint64) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func float32ListToString(l []float32) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func float64ListToString(l []float64) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func interfaceListToString(l []interface{}) string {
	var stringList []string = nil
	for _, v := range l {
		stringList = append(stringList, valueToString(v))
	}
	return integerListToString(stringList)
}

func interfaceMapToString(data map[string]interface{}) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func byteMapToString(data map[string]byte) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func int8MapToString(data map[string]int8) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func int16MapToString(data map[string]int16) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func intMapToString(data map[string]int) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func int32MapToString(data map[string]int32) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func int64MapToString(data map[string]int64) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func uint8MapToString(data map[string]uint8) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func uint16MapToString(data map[string]uint16) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func uint32MapToString(data map[string]uint32) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func uint64MapToString(data map[string]uint64) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func float32MapToString(data map[string]float32) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func float64MapToString(data map[string]float64) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}

func stringMapToString(data map[string]string) string {
	ret := "{"
	index := 0
	for k, v := range data {
		ret += fmt.Sprintf(`"%s": %s`, k, valueToString(v))
		if index != len(data)-1 {
			ret += ","
		}
		index++
	}
	ret += "}"
	return ret
}
