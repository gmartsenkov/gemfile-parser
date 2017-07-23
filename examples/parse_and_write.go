package main

import (
	"fmt"
	"os"

	"github.com/gmartsenkov/gemfile-parser"
)

func main() {
	gf, err := os.OpenFile("Gemfile", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer gf.Close()
	parsedGemfile := gemfile.Gemfile{}
	parsedGemfile.Parse(gf)
	parsedGemfile.Write(os.Stdout)
}
