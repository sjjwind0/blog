package html

import (
	"plugin/build/step"
)

type htmlBuilder struct {
}

func NewHtmlBuilder() *htmlBuilder {
	return new(htmlBuilder)
}

func (h *htmlBuilder) BuildStep() []step.BuildStep {
	return nil
}
