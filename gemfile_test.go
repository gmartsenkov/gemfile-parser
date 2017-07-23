package gemfile_test

import (
	"bytes"
	"captain/parsers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var validGemfile = `
source "rubygems"
ruby "2.2.3"

gem "roo", "~>5", require: "lib"
gem "boo", git: "github.com", require: true
gem "moo", git: "github.com", tag: "1.0", require: false
group :development, :test do
  gem "soo", git: "github.com", tag: "1.0"
end
group :test do
  gem "faker", "> 5.0"
end

gem "goo", "~>1"
`

var expectedGemfile = `
source "rubygems"
ruby "2.2.3"

gem "roo", "~>5", require: "lib"
gem "boo", git: "github.com", require: true
gem "moo", git: "github.com", tag: "1.0", require: false
gem "goo", "~>1"

group :development, :test do
  gem "soo", git: "github.com", tag: "1.0"
end

group :test do
  gem "faker", "> 5.0"
end
`

var _ = Describe("Gemfile", func() {
	Describe("Parse", func() {
		It("parses the gemfile correctly", func() {
			gemfile := parsers.Gemfile{}
			gemfileBuffer := bytes.NewBufferString(validGemfile)
			gemfile.Parse(gemfileBuffer)
			Expect(gemfile).To(BeEquivalentTo(
				parsers.Gemfile{
					Source: "rubygems",
					Ruby:   "2.2.3",
					Gems: []*parsers.Gem{
						&parsers.Gem{
							Name: "roo", Git: "", Version: "~>5", Tag: "", Groups: "", Require: "\"lib\"",
						},
						&parsers.Gem{
							Name: "boo", Git: "github.com", Version: "", Tag: "", Groups: "", Require: "true",
						},
						&parsers.Gem{
							Name: "moo", Git: "github.com", Version: "", Tag: "1.0", Groups: "", Require: "false",
						},
						&parsers.Gem{
							Name: "soo", Git: "github.com", Version: "", Tag: "1.0", Groups: ":development, :test", Require: "",
						},
						&parsers.Gem{
							Name: "faker", Git: "", Version: "> 5.0", Tag: "", Groups: ":test", Require: "",
						},
						&parsers.Gem{
							Name: "goo", Git: "", Version: "~>1", Tag: "", Groups: "", Require: "",
						},
					},
				},
			))
		})
	})

	Describe("FindGem", func() {
		It("returns the correct gem", func() {
			gemfile := parsers.Gemfile{}
			gemfileBuffer := bytes.NewBufferString(validGemfile)
			gemfile.Parse(gemfileBuffer)
			err, gem := gemfile.FindGem("faker")
			Expect(err).To(BeNil())
			Expect(gem).To(BeEquivalentTo(
				&parsers.Gem{Name: "faker", Git: "", Version: "> 5.0", Tag: "", Groups: ":test", Require: ""},
			))
		})

		It("updates the gem correctly", func() {
			gemfile := parsers.Gemfile{}
			gemfileBuffer := bytes.NewBufferString(validGemfile)
			gemfile.Parse(gemfileBuffer)
			err, gem := gemfile.FindGem("faker")
			Expect(err).To(BeNil())
			Expect(gem).To(BeEquivalentTo(
				&parsers.Gem{Name: "faker", Git: "", Version: "> 5.0", Tag: "", Groups: ":test", Require: ""},
			))
			gem.Version = "6.0"
			gem.Require = "true"
			err, gem = gemfile.FindGem("faker")
			Expect(err).To(BeNil())
			Expect(gem).To(BeEquivalentTo(
				&parsers.Gem{Name: "faker", Git: "", Version: "6.0", Tag: "", Groups: ":test", Require: "true"},
			))
		})
		Context("with unexisting gem", func() {
			It("returns an error", func() {
				gemfile := parsers.Gemfile{}
				gemfileBuffer := bytes.NewBufferString(validGemfile)
				gemfile.Parse(gemfileBuffer)
				err, gem := gemfile.FindGem("unexisting_gem")
				Expect(err).To(MatchError("Gem not found unexisting_gem"))
				Expect(gem).To(BeEquivalentTo(&parsers.Gem{}))
			})
		})
	})

	Describe("Write", func() {
		It("writes correctly", func() {
			gemfile := parsers.Gemfile{}
			gemfileBuffer := bytes.NewBufferString(validGemfile)
			buffer := bytes.NewBuffer([]byte{})
			gemfile.Parse(gemfileBuffer)
			gemfile.Write(buffer)
			Expect(buffer.String()).To(BeEquivalentTo(expectedGemfile))
		})
	})
})
