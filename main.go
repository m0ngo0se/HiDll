package main

import (
	"HiDll/preload"
	"fmt"
	"github.com/c-bata/go-prompt"
)

func MainCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "preload", Description: "select preload mode"},
		{Text: "postload", Description: "select postload mode"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func MainExecutor(s string) {
	if s == "preload" {
		p := prompt.New(
			preload.PreExecutor,
			preload.PreCompleter,
			prompt.OptionPrefix("[HiDll]>preload>"),
		)
		p.Run()
	} else if s == "postload" {

	} else {
		fmt.Println("No Mode:" + s)
	}
}

func main() {

	p := prompt.New(
		MainExecutor,
		MainCompleter,
		prompt.OptionTitle("Hi dll"),
		prompt.OptionPrefix("[HiDll]>"),
	)
	p.Run()
}
