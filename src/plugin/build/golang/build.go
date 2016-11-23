package golang

import (
	"fmt"
	"framework/base/config"
	"os"
	"path/filepath"
	"plugin/build/step"
)

type GolangBuilder struct {
	projectPath string
}

func NewGolangBuilder(projectPath string) *GolangBuilder {
	return &GolangBuilder{projectPath: projectPath}
}

/*
1. pwd log current path
2. rm native handler
3. copy plugin handler to src
4. build
5. install
*/
func (g *GolangBuilder) BuildStep() []step.BuildStep {
	toolPath := config.GetDefaultConfigJsonReader().GetString("storage.file.tool")
	toolPath = filepath.Join(toolPath, "build", "golang", "handler")
	var GOPATH string = "/home/wind/pkg" + ":" + g.projectPath
	var GOROOT string = "/home/wind/Application/go"
	var PATH string = os.Getenv("PATH") + ":" + "/home/wind/Application/go/bin"
	var golangEnv = []string{"GOROOT=" + GOROOT, "GOPATH=" + GOPATH, "PATH=" + PATH}
	fmt.Println(golangEnv)
	return []step.BuildStep{
		step.NewShellCommandStep(nil, g.projectPath, "pwd", nil),
		step.NewShellCommandStep(nil, g.projectPath, "rm", []string{"-r", "src/handler"}),
		step.NewShellCommandStep(nil, g.projectPath, "cp", []string{
			"-r",
			toolPath,
			"src/",
		}),
		step.NewShellCommandStep(golangEnv, g.projectPath, "/bin/bash", []string{"build.sh"}),
		step.NewShellCommandStep(golangEnv, g.projectPath, "/bin/bash", []string{"install.sh"}),
	}
}
