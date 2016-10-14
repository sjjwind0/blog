package build

import (
	"fmt"
)

type ProgressCallback func(info string, err string, isComplete bool)

type Builder interface {
	BuildStep() []BuildStep
}

type BuilderMgr struct {
	builder Builder
}

func NewBuildMgr(builder Builder) *BuilderMgr {
	return &BuilderMgr{builder: builder}
}

func (b *BuilderMgr) Run(callback ProgressCallback) {
	go func() {
		var outputStr string = ""
		var errorStr string = ""
		buildStepList := b.builder.BuildStep()
		for index, step := range buildStepList {
			description := fmt.Sprintf("%d. %s\n", index+1, step.Description())
			callback(description, "", false)
			outString, errString, err := step.Run()
			if err != nil {
				errString += err.Error()
				return
			}
			outputStr += outString
			errorStr += errString
			callback(outString, errString, false)
		}
		callback(outputStr, errorStr, true)
	}()
}
