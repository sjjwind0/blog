package node

import (
	"framework/base/config"
	"path/filepath"
	"plugin/build/step"
)

type NodeBuilder struct {
	projectPath string
}

func NewNodeBuilder(projectPath string) *NodeBuilder {
	return &NodeBuilder{projectPath: projectPath}
}

/*
1. pwd log current path
2. rm native handler
3. copy plugin handler to src
4. build
5. install
*/
func (g *NodeBuilder) BuildStep() []step.BuildStep {
	toolPath := config.GetDefaultConfigJsonReader().GetString("storage.file.tool")
	handlePath := filepath.Join(toolPath, "build", "node", "handler")
	modulePath := filepath.Join(toolPath, "build", "node", "node_modules")
	return []step.BuildStep{
		step.NewShellCommandStep(nil, g.projectPath, "pwd", nil),
		step.NewShellCommandStep(nil, g.projectPath, "rm", []string{"-r", "handler"}),
		step.NewShellCommandStep(nil, g.projectPath, "cp", []string{
			"-r",
			handlePath,
			"./",
		}),
		step.NewShellCommandStep(nil, g.projectPath, "cp", []string{
			"-r",
			modulePath,
			"./",
		}),
	}
}
