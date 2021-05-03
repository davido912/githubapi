package main

import (
	"fmt"
	"os"
	"os/exec"
)


func main() {
	v, _ := exec.LookPath("vim")
	fmt.Println(v)
	cmd :=exec.Cmd{Path: v, Args: []string{"vi", "test.txt"}, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	cmd.Run()
}

const (
	UrlGit string = "https://api.github.com/repos/davido912/data-proc/issues"
	// UrlGit string = "https://api.github.com/repos/apache/airflow/issues"

)