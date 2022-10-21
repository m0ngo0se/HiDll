package preload

import (
	"HiDll/util"
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
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
			CreatePreProject(args[1], args[2])
		} else {
			fmt.Println("input error!")
		}
	}
}

func PreCompleter2(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "noprex", Description: "noprex"},
		{Text: "nostuff", Description: "nostuff"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func PreExecutor2(s string) {

}

func CopyExes(srcpath, dstpath string) {
	files, err := os.ReadDir(srcpath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			newsrcpath := fmt.Sprintf("%s\\%s", srcpath, file.Name())
			CopyExes(newsrcpath, dstpath)
		} else {
			newsrcpath := fmt.Sprintf("%s\\%s", srcpath, file.Name())
			newdstpath := fmt.Sprintf("%s\\%s", dstpath, file.Name())
			util.CopyFile(newsrcpath, newdstpath)
		}
	}
}

func CreatePreProject(name, path string) {
	util.CreateDir(name)
	CopyExes(path, name)
	project := fmt.Sprintf("%s%s>", "[HiDll]>preload>", name)
	p := prompt.New(
		PreExecutor2,
		PreCompleter2,
		prompt.OptionPrefix(project),
	)
	p.Run()
}
