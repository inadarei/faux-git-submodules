package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"path"
	"errors"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func abort(err string) {
	fmt.Println("ERROR: " + err)
	os.Exit(1)
}

type cfgOpts struct {
	filepath string
}

type gitModule struct {
	path   string
	branch string
	url    string
}

type pathState int

const (
	pathMissing pathState = iota
	pathFolder
	pathNonFolder
)

func main() {

	modules := getModules(cfgOpts{})

	execGitCommands(modules)

	for _, module := range modules {
		addIfMissing(".gitignore", module.path)
	}

}

func displayHelp(missingConfig bool) {

	colorReset := "\033[0m"
    colorRed := "\033[31m"
    //colorGreen := "\033[32m"
    colorYellow := "\033[33m"
    //colorBlue := "\033[34m"
    //colorPurple := "\033[35m"
    //colorCyan := "\033[36m"
    //colorWhite := "\033[37m"

	if (missingConfig) {
		fmt.Println(string(colorRed), "Missing: ", string(colorReset), "fgs.json configuration file!")
		fmt.Println("For an example fgs.json, see: https://github.com/inadarei/faux-git-submodules")
		fmt.Println(" 	")
	}
    
	progName := path.Base(os.Args[0])
	fmt.Println("Command utility to easily check-out faux git submodules.");
	fmt.Println("For more information: ","https://github.com/inadarei/faux-git-submodules/");
	fmt.Println(" 	")
	fmt.Println("	", string(colorYellow), "Usage: ", string(colorReset), progName);
}

/**
* Check if file exists. If yes and a string not already contained, add at
* the end
 */
func addIfMissing(filePath string, needle string) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		abort(fmt.Sprintf("%s", err))
	}
	content := string(b)
	if !strings.Contains(content, needle) {
		content += fmt.Sprintf("%s\n", needle)
		ioutil.WriteFile(filePath, []byte(content), os.ModeAppend)
	}
}

/**
* Properly execute a shell command without swallowing outputs or any of the
* silly things Go does by default
 */
func execute(command string, inDir string) string {
	//cmd = cmd + " 1>&2"

	var output string
	var lastDir string
	var err error

	if inDir != "" {
		lastDir, err = os.Getwd()
		check(err)
		os.Chdir(inDir)
	}

	args := strings.Split(command, " ")
	app := args[0]
	args = args[1:]
	cmd := exec.Command(app, args...)

	stderr, err := cmd.StderrPipe()
	check(err)

	stdout, err2 := cmd.StdoutPipe()
	check(err2)

	err3 := cmd.Start()
	check(err3)

	slurpErr, _ := ioutil.ReadAll(stderr)
	if len(slurpErr) > 0 {
		output += fmt.Sprintf("%s\n", slurpErr)
	}

	slurpOut, _ := ioutil.ReadAll(stdout)
	if len(slurpOut) > 0 {
		output += fmt.Sprintf("%s\n", slurpOut)
	}

	errW := cmd.Wait()

	if inDir != "" {
		os.Chdir(lastDir)
	}

	if errW != nil {
		_ = err // can add to ouput but not useful for git commands
	}

	return output
}

/* debugCurrPath is used for debugging.
*  It should be the folder where command was
* invoked from, even if executable is in another folder
 */
func debugCurrPath() {
	cmd := "echo $PWD"
	out, err := exec.Command("sh", "-c", cmd).Output()
	check(err)
	fmt.Println(string(out))
}

/**
* Parse config file and return configuration
 */
func getModules(opts cfgOpts) []gitModule {
	modules := []gitModule{}

	if opts.filepath == "" {
		opts.filepath = "./fgs.json"
	}

	if _, err := os.Stat(opts.filepath); errors.Is(err, os.ErrNotExist) {
		displayHelp(true)
		os.Exit(1);
	}
	bRepos, err := ioutil.ReadFile(opts.filepath)
	check(err)

	var msgMapTemplate interface{}
	err3 := json.Unmarshal([]byte(bRepos), &msgMapTemplate)
	check(err3)
	repos := msgMapTemplate.(map[string]interface{})

	var sBranch, sURL string

	for path, cfg := range repos {
		url := cfg.(map[string]interface{})["url"]
		if url == nil {
			abort("module configuration with missing URL: " + path)
		} else {
			sURL = url.(string)
		}

		branch := cfg.(map[string]interface{})["branch"]
		if branch == nil {
			sBranch = "main"
		} else {
			sBranch = branch.(string)
		}

		module := gitModule{
			path:   path,
			branch: sBranch,
			url:    sURL}

		modules = append(modules, module)
	}

	return modules
}

/*
* Turn config into appropriate git commands. If a destination already
* exists then we do git pull, if not: git clone
 */
func execGitCommands(modules []gitModule) {

	var cmd string
	for _, module := range modules {
		pathIs := checkPath(module.path)
		if pathIs == pathFolder {
			cmd = fmt.Sprintf("git pull")
			fmt.Println(cmd)
			fmt.Print(execute(cmd, module.path))
		} else if pathIs == pathMissing {
			cmd = fmt.Sprintf("git clone -b %s %s %s", module.branch, module.url, module.path)
			fmt.Println(cmd)
			fmt.Print(execute(cmd, ""))
		} else { // pathIs == pathNonFolder
			abort("Path already exists and is not a folder: " + module.path)
		}
	}

}

/**
* check whether path exists and is a folder.
*
 */
func checkPath(path string) pathState {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return pathMissing
	}

	if err == nil {
		if fi.IsDir() {
			return pathFolder
		}
		return pathNonFolder
	}

	panic(err) // some other error
}
