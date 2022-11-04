package preload

import (
	"HiDll/util"
	"fmt"
	"github.com/c-bata/go-prompt"
	peparser "github.com/saferwall/pe"
	"os"
	"strings"
)

type preinfo struct {
	path     string
	importab map[string]int
}

var (
	project    string
	prexes     map[string]preinfo
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
		{Text: "get", Description: "get iat funcs"},
		{Text: "exit", Description: "exit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func getPreFuncs(exename, dllname string) {
	output := fmt.Sprintf("%s.txt", dllname)
	if exeinfo, ok := prexes[exename]; ok {
		if _, ok := exeinfo.importab[dllname]; ok {
			pe, err := peparser.New(exeinfo.path, &peparser.Options{})
			if err != nil {
				pe.Close()
				println("get funcs error")
				os.Exit(0)
			}
			err = pe.Parse()
			if err != nil {
				pe.Close()
				println("get funcs error")
				os.Exit(0)
			}
			for _, value := range pe.IAT {
				dll, funcname, _ := strings.Cut(value.Meaning, "!")
				if dll == dllname {
					functext := fmt.Sprintf("extern \"C\" __declspec(dllexport) void %s()\n{\n\treturn;\n}\n", funcname)
					util.Writedata(output, functext)
				}
			}
		}
	}
}

func preExecutor2(s string) {
	if s == "exit" {
		os.Exit(0)
	} else {
		args := util.ParseCmd(s)
		if args[0] == "get" && len(args) == 3 {
			getPreFuncs(args[1], args[2])
		} else {
			fmt.Println("input error!")
		}
	}
}

func NoPrex(prex string) {
	for name, info := range prexes {
		imptable := info.importab
		for dllname, _ := range imptable {
			if strings.HasPrefix(strings.ToLower(dllname), prex) {
				err := os.Remove(info.path)
				if err != nil {
					println(err)
				}
				delete(prexes, name)
				break
			}
		}
	}
}

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
		if pe.Is32 || !pe.HasIAT || !pe.IsSigned || pe.HasCLR {
			pe.Close()
			os.Remove(filepath)
			continue
		}
		exeinfo := preinfo{
			path:     filepath,
			importab: make(map[string]int),
		}
		for _, value := range pe.IAT {
			dllname, _, _ := strings.Cut(value.Meaning, "!")
			if dllname != "" {
				if _, ok := exeinfo.importab[dllname]; !ok {
					exeinfo.importab[dllname] = 1
				}
			}
		}
		pe.Close()
		if len(exeinfo.importab) == 0 {
			os.Remove(filepath)
		} else {
			prexes[file.Name()] = exeinfo
		}
	}
	for _, prex := range blackprexs {
		NoPrex(prex)
	}
}

func FilterKnown() {
	for name, info := range prexes {
		deldlls := make([]string, 0)
		for dllname, _ := range info.importab {
			path := fmt.Sprintf("C:\\Windows\\System32\\%s", dllname)
			path2 := fmt.Sprintf("C:\\Windows\\SysWOW64\\%s", dllname)
			if util.PathExist(path) || util.PathExist(path2) {
				deldlls = append(deldlls, dllname)
			}
		}
		for _, dllname := range deldlls {
			delete(info.importab, dllname)
		}
		if len(info.importab) == 0 {
			os.Remove(info.path)
			delete(prexes, name)
		}
	}
}

func CheckNum() {
	delexes := make([]string, 0)
	for name, info := range prexes {
		if len(info.importab) >= 2 {
			os.Remove(info.path)
			delexes = append(delexes, name)
		}
	}
	for _, exename := range delexes {
		delete(prexes, exename)
	}
}

func show() {
	if len(prexes) == 0 {
		fmt.Println("all exes has deleted!")
		os.Exit(0)
	}
	fmt.Printf("Exe num:%d\n", len(prexes))
	for _, value := range prexes {
		fmt.Println(value)
	}
}

func init() {
	prexes = make(map[string]preinfo, 0)
}

func createProject(name, path string) {
	project = ".\\" + name
	if util.CreateDir(project) {
		util.CopyExes(path, project)
	}
	initList(project)
	FilterKnown()
	CheckNum()
	show()
	project := fmt.Sprintf("%s%s>", "[HiDll]>preload>", name)
	p := prompt.New(
		preExecutor2,
		preCompleter2,
		prompt.OptionPrefix(project),
	)
	p.Run()
}
