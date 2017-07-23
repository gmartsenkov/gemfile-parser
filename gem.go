package gemfile

import (
	"fmt"
	"io"
)

type Gem struct {
	Name    string
	Git     string
	Version string
	Tag     string
	Groups  string
	Require string
}

func (gem *Gem) write(writer io.Writer) {
	var line string
	line += fmt.Sprintf(`gem "%s"`, gem.Name)
	if gem.Version != "" {
		line += fmt.Sprintf(`, "%s"`, gem.Version)
	}
	if gem.Git != "" {
		line += fmt.Sprintf(`, git: "%s"`, gem.Git)
	}
	if gem.Tag != "" {
		line += fmt.Sprintf(`, tag: "%s"`, gem.Tag)
	}
	if gem.Require != "" {
		line += fmt.Sprintf(`, require: %s`, gem.Require)
	}
	io.WriteString(writer, line+"\n")
}
