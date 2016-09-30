package storage

import (
	"testing"
	"fmt"
)

func Test_ReadConfigContent(t *testing.T)  {
	s := NewPluginStorage("/Users/sjjwind/Downloads/test/zip/raw.zip")
	err := s.Run()
	if err != nil {
		fmt.Println("run error: ", err.Error())
	}
}