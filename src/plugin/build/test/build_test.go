package test

import (
	"fmt"
	"plugin/build"
	"plugin/build/golang"
	"testing"
	"time"
)

func Test_BuildGolang(t *testing.T) {
	testBuildRootPath := "/home/wind/data/go-plugin-demo"
	builder := golang.NewBuilder(testBuildRootPath)
	builderMgr := build.NewBuildMgr(builder)
	complete := false
	builderMgr.Run(func(info string, err string, isComplete bool) {
		if isComplete {
			fmt.Println("complete")
			complete = true
			return
		}
		fmt.Println(info)
		if err != "" {
			fmt.Println(err)
		}
	})
	for {
		if complete {
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}
}
