package preload

import (
	"HiDll/util"
	"fmt"
	"github.com/c-bata/go-prompt"
	peparser "github.com/saferwall/pe"
	"os"
	"strings"
)

type prenfo struct {
	name     string
	path     string
	importab map[string]int
}

var (
	project    string
	prexes     []prenfo
	knowndlls  = []string{"kernel32.exe"}
	blackprexs = []string{"api-ms-win", "vcruntime"}
)

func PreCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "create", Description: "create a project to start: create [projectname] [folderpath]"},
		{Text: "exit", Description: "exit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func PreExecutor(s string) {
	if s == "exit" {
		os.Exit(0)
	} else {
		args := util.ParseCmd(s)
		if args[0] == "create" && len(args) == 3 {
			createProject(args[1], args[2])
		} else {
			fmt.Println("input error!")
		}
	}
}

func preCompleter2(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "noprex", Description: "noprex"},
		{Text: "nostuff", Description: "nostuff"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func preExecutor2(s string) {

}

func NoPrex(prex string) {
	for idx, info := range prexes {
		imptable := info.importab
		for name, _ := range imptable {
			if strings.HasPrefix(strings.ToLower(name), prex) {
				delete(imptable, name)
			}
		}

	}
}

//
//func NoStuff() {
//	for idx, info := range prexes {
//
//	}
//}

func initList(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filepath := fmt.Sprintf("%s\\%s", path, file.Name())
		pe, err := peparser.New(filepath, &peparser.Options{})
		if err != nil {
			pe.Close()
			os.Remove(filepath)
			continue
		}
		err = pe.Parse()
		if err != nil {
			pe.Close()
			os.Remove(filepath)
			continue
		}
		if pe.Is32 || !pe.HasIAT {
			pe.Close()
			os.Remove(filepath)
			continue
		}
		exe := prenfo{
			name:     file.Name(),
			path:     filepath,
			importab: make(map[string]int),
		}
		for _, value := range pe.IAT {
			dllname, _, _ := strings.Cut(value.Meaning, "!")
			if _, ok := exe.importab[dllname]; !ok {
				exe.importab[dllname] = 1
			}
		}
		prexes = append(prexes, exe)
		pe.Close()
	}
}

func filterKnown() {

}

func show() {
	if len(prexes) == 0 {
		fmt.Println("all exes has deleted!")
		os.Exit(0)
	}
	for _, value := range prexes {
		fmt.Println(value)
	}
}

func createProject(name, path string) {
	project = ".\\" + name
	if util.CreateDir(project) {
		util.CopyExes(path, project)
	}
	initList(project)
	show()
	project := fmt.Sprintf("%s%s>", "[HiDll]>preload>", name)
	p := prompt.New(
		preExecutor2,
		preCompleter2,
		prompt.OptionPrefix(project),
	)
	p.Run()
}
