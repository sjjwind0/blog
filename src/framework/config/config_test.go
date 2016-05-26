package config

import (
	"encoding/json"
	"testing"
)

func Test_ReadConfigContent(t *testing.T) {
	configContent := `
	{
		"A": {
			"B": {
				"C": 1,
				"D": "2",
				"E": [1, 2, 3]
			}
		},
		"F": 2
	}
	`
	config := NewConfigContentManager(configContent)

	if config.ReadConfig("A.B.C").(int64) != 1 {
		t.Error("Error")
	}

	if config.ReadConfig("A.B.D").(string) != "2" {
		t.Error("Error")
	}

	retList := config.ReadConfig("A.B.E").([]interface{})
	comList := []int{1, 2, 3}
	for i := range retList {
		retValue, _ := retList[i].(json.Number).Int64()
		comValue := int64(comList[i])
		if retValue != comValue {
			t.Error("Error")
		}
	}
}
