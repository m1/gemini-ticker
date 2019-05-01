package main

import (
	"github.com/m1/gemini-ticker/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		return
	}
}