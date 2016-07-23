package json

import (
	"fmt"
	"testing"
)

func Test_IntergerMap(t *testing.T) {
	var data map[string]int = make(map[string]int)
	data["test_1"] = 1
	data["test_2"] = 2
	data["test_100"] = 100
	transferString := ToJsonString(data)
	fmt.Println(transferString)
}

func Test_StringMap(t *testing.T) {
	var data map[string]string = make(map[string]string)
	data["test_1"] = "1"
	data["test_2"] = "2"
	data["test_100"] = "100"
	transferString := ToJsonString(data)
	fmt.Println(transferString)
}

func Test_IntListMap(t *testing.T) {
	var data map[string]interface{} = make(map[string]interface{})
	data["test_1"] = []int{1, 2, 3}
	data["test_2"] = []int{4, 5, 6}
	data["test_100"] = []int{100, 200, 300}
	transferString := ToJsonString(data)
	fmt.Println(transferString)
}

func Test_InterfaceMap(t *testing.T) {
	var data map[string]interface{} = make(map[string]interface{})
	data["test_1"] = []int{1, 2, 3}
	data["test_2"] = []int{"4", "5", "6"}
	data["test_100"] = []int{100, 200, 300}
	data["test_200"] = 5
	data["test_300"] = map[string]interface{}{"1": 2, "3": 4, "4": "6"}
	transferString := ToJsonString(data)
	fmt.Println(transferString)
}
