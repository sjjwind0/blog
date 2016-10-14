package golang

import (
	"plugin/build"
)

type GolangBuilder struct {
	projectPath string
}

func NewBuilder(projectPath string) *GolangBuilder {
	return &GolangBuilder{projectPath: projectPath}
}

func (g *GolangBuilder) BuildStep() []build.BuildStep {
	return []build.BuildStep{
		build.NewShellCommandStep(nil, g.projectPath, "pwd", nil),
		build.NewShellCommandStep(nil, g.projectPath, "rm", []string{"-r", "src/handler"}),
		build.NewShellCommandStep(nil, g.projectPath, "cp", []string{
			"-r",
			"/home/wind/Project/blog/tools/build/golang/handler",
			"src/",
		}),
		build.NewShellCommandStep([]string{
			"GOPATH=/home/wind/pkg:/home/wind/data/go-plugin-demo",
			"GOROOT=/home/wind/Application/go",
			"PATH=/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/home/wind/Application/go/bin:/home/wind/pkg/bin:/home/wind/Application/node-v4.5.0-linux-x64/bin:/home/wind/.local/bin:/home/wind/bin",
		}, g.projectPath, "./run.sh", nil),
		build.NewShellCommandStep(nil, g.projectPath, "./plugin", nil),
	}
}
