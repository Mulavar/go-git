package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var gitUserName = "Mulavar"

func main() {
	//"/Users/lam/go/src/github.com/apache/dubbo-go"
	// TODO args
	err := scan("/Users/lam/IdeaProjects")
	if err != nil {
		fmt.Println("cuowu,", err)
		return
	}
}

func execScript(workDir, gitUserName string) {
	var gitLog []byte
	var err error

	if err := os.Chdir(workDir); err != nil {
		log.Fatal(err)
	}

	commandLine := fmt.Sprintf("git log --author=%s --pretty=tformat: --numstat ", gitUserName)
	cmd := exec.Command("/bin/sh", "-c", commandLine)
	if gitLog, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO filter
	if gitLog == nil || len(gitLog) == 0 {
		return
	}
	commandLine += "| awk '{ add += $1 ; subs += $2 ; loc += $1 - $2 } END { printf \"added lines: %s removed lines : %s total lines: %s\\n\",add,subs,loc }'"

	//fmt.Println(commandLine)

	paths := strings.Split(workDir, "/")
	path := paths[len(paths)-1]
	cmd = exec.Command("/bin/sh", "-c", commandLine)

	if gitLog, err = cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO regex
	// TODO log format with color
	fmt.Println("project:", path, ", stat:", string(gitLog))
}

func scan(path string) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo)
		} else {
			files = append(files, fileInfo)
		}
	}

	for _, dir := range dirs {
		if dir.Name() == ".git" {
			execScript(path, gitUserName)
			return nil
		}
		scan(path + "/" + dir.Name())
	}

	return nil
}
