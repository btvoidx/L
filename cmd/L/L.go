package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
	"github.com/spf13/pflag"

	"github.com/btvoidx/L"
	"github.com/btvoidx/L/internal/logger"
)

func main() {
	var (
		listFlag    bool
		helpFlag    bool
		silentFlag  bool
		verboseFlag bool
		entrypoint  string
		initFlag    bool
	)

	pflag.BoolVarP(&listFlag, "list", "l", false, "lists all tasks")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows L usage")
	pflag.BoolVar(&silentFlag, "silent", false, "disables output from L")
	pflag.BoolVar(&verboseFlag, "verbose", false, "enables verbose mode")
	pflag.StringVarP(&entrypoint, "taskfile", "f", "tasks.lua", "choose tasks file")
	pflag.BoolVar(&initFlag, "init", false, "creates a default tasks.lua file")
	pflag.Parse()

	log := logger.Logger{
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Verbose: verboseFlag,
		Silent:  silentFlag,
	}

	pflag.Usage = func() {
		var usage = `L 0.0.0

		Usage: %s
		Runs the specified task(s). Falls back to the %s if no task name was specified.

		Example: %s with the following %s file will generate an %s file with the content "Hello World!".

		'''
		function task.hello()
			print("Writing to a file named 'output.txt' now...")
			file = io.open("output.txt", "w")
			file:write("Hello World!")
			file:close()
			print("Done writing!")
		end
		'''

		Options:
		` + pflag.CommandLine.FlagUsages()

		log.Write(strings.ReplaceAll(usage, "\t", ""),
			termenv.String("L --list --silent [task1 task2 ...]"),
			termenv.String("default"),
			termenv.String("L hello"),
			termenv.String("tasks.lua"),
			termenv.String("output.txt"),
		)
	}

	if helpFlag {
		pflag.Usage()
		return
	}

	if initFlag {
		dir, _ := os.Getwd()

		_, err := os.Open("tasks.lua")
		if err == nil {
			log.Err("L: %s already exists in %s", termenv.String("tasks.lua"), termenv.String(dir))
			return
		}

		f, err := os.Create("tasks.lua")
		if err != nil {
			panic(err) // todo!
		}
		defer f.Close()

		var contents = `-- https://github.com/btvoidx/L

			function task.default()
			  description 'Says "Hello World"'
			  print("Hello World")
			end
		`

		_, err = fmt.Fprint(f, strings.ReplaceAll(contents, "\t", ""))
		if err != nil {
			panic(err) // todo!
		}

		log.Write("L: created %s in %s", termenv.String("tasks.lua"), termenv.String(dir))

		return
	}

	e := L.Executor{
		Entrypoint: entrypoint,
		Logger:     &log,
	}

	if err := e.Compile(); err != nil {
		s := err.Error()
		if strings.HasPrefix(s, "parse error") {
			s = strings.ReplaceAll(s, "parse error: ", "")
			s = strings.ReplaceAll(s, ":   parse error", "")
			log.Err("L: parse error:\n%s", s)
			return
		}

		if strings.HasPrefix(s, "open") && strings.HasSuffix(s, "no such file or directory") {
			dir, _ := os.Getwd()
			log.Err("L: %s was not found in %s: use %s to create a new one",
				termenv.String(entrypoint),
				termenv.String(dir),
				termenv.String("L --init"))
			return
		}

		log.Err("L: error:\n%s", s)
		return
	}

	if listFlag {
		tasks, err := e.List()
		if err != nil {
			log.Err("L: %s", err)
			return
		}

		if len(tasks) == 0 {
			log.Write("L: no tasks available")
			return
		}

		taskNames := make([]string, 0, len(tasks))
		for _, t := range tasks {
			taskNames = append(taskNames, t.Name)
		}

		log.Write("L: all available tasks:\n- %s", strings.Join(taskNames, "\n- "))
		return
	}

	taskNames := pflag.Args()
	if len(taskNames) == 0 {
		taskNames = []string{"default"}
	}

	for _, tn := range taskNames {
		if _, err := e.Run(tn); err != nil {
			log.Err("L: %s", err.Error())
		}
	}
}
