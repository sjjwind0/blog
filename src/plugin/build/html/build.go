package html

import (
	"plugin/build"
)

type HtmlBuilder struct {
}

func (h *HtmlBuilder) BuildStep() []*build.BuildStep {
	// html, no build step
	return nil
}
