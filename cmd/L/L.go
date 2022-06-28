package main

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/btvoidx/L"
	"github.com/btvoidx/L/internal/color"
	"github.com/btvoidx/L/internal/logger"
)

func main() {
	var (
		listFlag    bool
		helpFlag    bool
		silentFlag  bool
		verboseFlag bool
		entrypoint  string
		timeout     time.Duration
		initFlag    bool
	)

	pflag.BoolVarP(&listFlag, "list", "l", false, "lists all tasks")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows L usage")
	pflag.BoolVar(&silentFlag, "silent", false, "disables output from L")
	pflag.BoolVar(&verboseFlag, "verbose", false, "enables verbose mode")
	pflag.StringVarP(&entrypoint, "taskfile", "f", "tasks.lua", "choose tasks file")
	pflag.DurationVarP(&timeout, "timeout", "t", 0, "sets a limit on how long each task can run; <=0 means no limit")
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
			
			Usage: ` + color.Magenta("L --list --silent [task1 task2 ...]") + `
			Runs the specified task(s). Falls back to the ` + color.Magenta("default") + ` if no task name was specified.

			Example: ` + color.Magenta("L hello") + ` with the following ` + color.Magenta("tasks.lua") + ` file will generate an ` + color.Magenta("output.txt") + ` file with the content "Hello World!".

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

	// if initFlag {
	// 	var contents = `
	// 		-- https://github.com/btvoidx/L

	// 		function task.default()
	// 		  local to_say = "Hello Loser (Lua+Cool+User)"
	// 		  description 'Says ' .. to_say
	// 		  print(to_say)
	// 		end
	// 		`

	// 	_, err := os.Open("tasks.lua")
	// 	if err == nil {
	// 		log.Err("L: %s already exists in %s",
	// 		color.Magenta("tasks.lua"),
	// 		color.Magenta(dir),)
	// 		return
	// 	}

	// 	f, err := os.Create("tasks.lua")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()

	// 	return
	// }

	if timeout <= 0 {
		// Seems like a fair deal.
		// It's not a max int64 value due to a possibility of overflow when adding L.initTimeout to this timeout.
		timeout = 24 * 31 * time.Hour
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
				color.Magenta(entrypoint),
				color.Magenta(dir),
				color.Magenta("L --init"))
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
		if _, err := e.Run(tn, timeout); err != nil {
			log.Err("L: %s", err.Error())
		}
	}
}
