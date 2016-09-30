package build

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Builder interface {
	BuildStep() []*BuildStep
}

type BuilderMgr struct {
	builder Builder
}

func NewBuildMgr(builder Builder) *BuilderMgr {
	return &BuilderMgr{builder: builder}
}

func (b *BuilderMgr) Run() (io.Reader, io.Reader, chan bool) {
	var outputStr string = ""
	var errorStr string = ""
	var outReader io.Reader = strings.NewReader(outputStr)
	var errReader io.Reader = strings.NewReader(errorStr)
	close := make(chan bool)

	go func() {
		buildStepList := b.builder.BuildStep()
		for index, step := range buildStepList {
			fmt.Println("run: ", index)
			step.Run()
			tmpOutStr, err := ioutil.ReadAll(step.Outer)
			if err != nil {
				errorStr += "\nexec error: " + err.Error()
				return
			}
			outputStr += tmpOutStr

			tmpErrStr, err := ioutil.ReadAll(step.Error)
			if err != nil {
				errorStr += "\nexec error: " + err.Error()
				return
			}
			errorStr += tmpErrStr
		}
	}()

	return outReader, errReader, close
}
