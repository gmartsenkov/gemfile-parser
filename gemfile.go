package gemfile

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type Gemfile struct {
	Source string
	Ruby   string
	Gems   []*Gem
}

var (
	// Gem attributes
	REGEX_GEM            = regexp.MustCompile(`gem\s*['|"](.*?)['|"]`)
	REGEX_VERSION        = regexp.MustCompile(`,\s*["|'](.*?)['|"]`)
	REGEX_GIT            = regexp.MustCompile(`git:\s*["|'](.*?)['|"]`)
	REGEX_TAG            = regexp.MustCompile(`tag:\s*["|'](.*?)['|"]`)
	REGEX_REQUIRE_STRING = regexp.MustCompile(`require:\s*["|'](.*?)['|"]`)
	REGEX_REQUIRE_BOOL   = regexp.MustCompile(`require:\s*(false|true)`)
	// Gemfile header
	REGEX_SOURCE = regexp.MustCompile(`source\s*['|"](.*?)['|"]`)
	REGEX_RUBY   = regexp.MustCompile(`ruby\s*['|"](.*?)['|"]`)
	// Gem Groups
	REGEX_GROUP     = regexp.MustCompile(`group\s(.*?)\sdo`)
	REGEX_END_GROUP = regexp.MustCompile(`^end$`)
)

func (gf *Gemfile) Parse(file io.Reader) {
	scanner := bufio.NewScanner(file)
	gf.parseGemfile(scanner, "")
}

func (gf *Gemfile) parseGemfile(scanner *bufio.Scanner, groups string) {
	for scanner.Scan() {
		line := scanner.Text()
		gem := Gem{}

		if REGEX_RUBY.MatchString(line) {
			gf.Ruby = REGEX_RUBY.FindStringSubmatch(line)[1]
		}
		if REGEX_SOURCE.MatchString(line) {
			gf.Source = REGEX_SOURCE.FindStringSubmatch(line)[1]
		}
		if REGEX_GEM.MatchString(line) {
			gem.Name = REGEX_GEM.FindStringSubmatch(line)[1]
		}
		if REGEX_VERSION.MatchString(line) {
			gem.Version = REGEX_VERSION.FindStringSubmatch(line)[1]
		}
		if REGEX_GIT.MatchString(line) {
			gem.Git = REGEX_GIT.FindStringSubmatch(line)[1]
		}
		if REGEX_TAG.MatchString(line) {
			gem.Tag = REGEX_TAG.FindStringSubmatch(line)[1]
		}
		if REGEX_GROUP.MatchString(line) {
			gf.parseGemfile(scanner, REGEX_GROUP.FindStringSubmatch(line)[1])
		}
		if REGEX_REQUIRE_STRING.MatchString(line) {
			gem.Require = fmt.Sprintf(`"%s"`, REGEX_REQUIRE_STRING.FindStringSubmatch(line)[1])
		}
		if REGEX_REQUIRE_BOOL.MatchString(line) {
			gem.Require = REGEX_REQUIRE_BOOL.FindStringSubmatch(line)[1]
		}
		if REGEX_END_GROUP.MatchString(line) && groups != "" {
			gf.parseGemfile(scanner, "")
		}
		if gem.Name != "" {
			gem.Groups = groups
			gf.Gems = append(gf.Gems, &gem)
		}
	}
}

func (gf *Gemfile) writeGems(writer io.Writer) {
	for _, gem := range gf.GemsWithoutGroups() {
		gem.write(writer)
	}

	for _, groups := range gf.UniqueGroups() {
		io.WriteString(writer, "\n")
		io.WriteString(writer, fmt.Sprintf(`group %s do`, groups)+"\n")
		for _, gem := range gf.GemsByGroups(groups) {
			io.WriteString(writer, "  ")
			gem.write(writer)
		}
		io.WriteString(writer, "end\n")
	}
}

func (gf *Gemfile) FindGem(name string) (error, *Gem) {
	for _, gem := range gf.Gems {
		if gem.Name == name {
			return nil, gem
		}
	}
	return fmt.Errorf("Gem not found %s", name), &Gem{}
}

func (gf *Gemfile) Write(writer io.Writer) {
	io.WriteString(writer, "\n") // just for tests
	if gf.Source != "" {
		io.WriteString(writer, fmt.Sprintf(`source "%s"`, gf.Source)+"\n")
	}
	if gf.Ruby != "" {
		io.WriteString(writer, fmt.Sprintf(`ruby "%s"`, gf.Ruby)+"\n")
	}
	io.WriteString(writer, "\n")
	gf.writeGems(writer)
}

func (gf *Gemfile) GemsWithoutGroups() []Gem {
	gems := []Gem{}
	for _, gem := range gf.Gems {
		if gem.Groups == "" {
			gems = append(gems, *gem)
		}
	}
	return gems
}

func (gf *Gemfile) UniqueGroups() []string {
	groups := []string{}
	for _, gem := range gf.Gems {
		if gem.Groups != "" && !isInArray(groups, gem.Groups) {
			groups = append(groups, gem.Groups)
		}
	}
	return groups
}

func isInArray(groups []string, group string) bool {
	for _, g := range groups {
		if g == group {
			return true
		}
	}
	return false
}

func (gf *Gemfile) GemsByGroups(groups string) []Gem {
	gems := []Gem{}
	for _, gem := range gf.Gems {
		if gem.Groups == groups {
			gems = append(gems, *gem)
		}
	}
	return gems
}
