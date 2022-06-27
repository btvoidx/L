package main

import (
	"os"
	"strings"

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
		timeout     uint
	)

	pflag.BoolVarP(&listFlag, "list", "l", false, "lists all tasks")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows L usage")
	pflag.BoolVar(&silentFlag, "silent", false, "disables output from L")
	pflag.BoolVar(&verboseFlag, "verbose", false, "enables verbose mode")
	pflag.StringVarP(&entrypoint, "taskfile", "f", "tasks.lua", "choose tasks file")
	pflag.UintVarP(&timeout, "timeout", "t", 0, "sets a limit on how long a task can run; 0 means no limit")
	pflag.Parse()

	log := logger.Logger{
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Verbose: verboseFlag,
		Silent:  silentFlag,
	}

	pflag.Usage = func() {
		var usage = `L 0.0.0
			
			Usage: L --list --silent [task1 task2 ...]
			Runs the specified task(s). Falls back to the "` + logger.Magenta("default") + `" task if no task name was specified.

			Example: "` + logger.Magenta("L hello") + `" with the following "` + logger.Magenta("tasks.lua") + `" file will generate an "` + logger.Magenta("output.txt") + `" file with the content "Hello World!".

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

		log.Write(strings.ReplaceAll(usage, "\t", ""))
	}

	if helpFlag {
		pflag.Usage()
		return
	}

	e := L.Executor{
		FilePath: entrypoint,
		Logger:   &log,
	}

	if err := e.Compile(); err != nil {
		s := err.Error()
		if strings.HasPrefix(s, "parse error") {
			s = strings.ReplaceAll(s, "parse error: ", "")
			s = strings.ReplaceAll(s, ":   parse error", "")
			log.Err("L | Parse error:\n%s", s)
			return
		}

		log.Err("L | Error:\n%s", s)
		return
	}

	if listFlag {
		tasks, err := e.List()
		if err != nil {
			println(err.Error()) // todo
			return
		}

		if len(tasks) == 0 {
			log.Write("L: no tasks available")
			return
		}

		log.Write("L: all available tasks:\n- %s", strings.Join(tasks, "\n- "))
		return
	}

	taskNames := pflag.Args()
	if len(taskNames) == 0 {
		taskNames = []string{"default"}
	}

	for _, tn := range taskNames {
		if _, err := e.Run(tn); err != nil {
			log.Err("L: task '%s': not found", tn)
		}
	}
}
