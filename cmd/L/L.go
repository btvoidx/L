package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/btvoidx/L"
)

func getArg(i int) string {
	if i >= len(os.Args) {
		return ""
	}
	return os.Args[i]
}

func main() {
	var (
		list bool
	)

	pflag.BoolVarP(&list, "list", "l", false, "lists all tasks")
	pflag.Parse()

	e := L.Executor{
		FilePath: "task.lua",
	}

	if err := e.Compile(); err != nil {
		println(err.Error())
		return
	}

	task := getArg(1)
	if task == "" {
		task = "default"
	}

	if _, err := e.Run(task); err != nil {
		println(err.Error())
		return
	}
}
