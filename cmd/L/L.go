package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/pflag"

	"github.com/btvoidx/L"
)

var (
	mainStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("48"))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("239"))
	errSymbol    = errStyle.Bold(true).SetString("L [error]").String()
	lSymbol      = mainStyle.Bold(true).SetString("L").String()
)

func main() {
	var (
		entrypointPath string
		helpFlag       bool
		listFlag       bool
		initFlag       bool
	)

	pflag.StringVarP(&entrypointPath, "taskfile", "f", "tasks.lua", "choose tasks file")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows L usage")
	pflag.BoolVarP(&listFlag, "list", "l", false, "lists all tasks")
	// pflag.BoolVar(&silentFlag, "silent", false, "disables output from L")
	// pflag.BoolVar(&verboseFlag, "verbose", false, "enables verbose mode")
	pflag.BoolVar(&initFlag, "init", false, "creates a default tasks.lua file")
	pflag.Parse()

	if helpFlag {
		kw := lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
		fn := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
		str := lipgloss.NewStyle().Foreground(lipgloss.Color("71"))

		usage := strings.ReplaceAll(lSymbol+` 0.0.0
			Usage: `+mainStyle.Render("L --list --silent [task1 task2 ...]")+`
			Runs the specified task(s). Falls back to `+mainStyle.Render("default")+` if no task name was specified.
			Example: `+mainStyle.Render("L hello")+` with the following `+mainStyle.Render("tasks.lua")+` file will generate an `+mainStyle.Render("output.txt")+` file with the content "Hello World!".
			'''
			`+kw.Render("function")+` task.`+fn.Render("hello")+`()
			  `+fn.Render("print")+`(`+str.Render(`"Writing to a file named 'output.txt' now..."`)+`)
			  `+kw.Render("local")+` file = io.`+fn.Render("open")+`(`+str.Render(`"output.txt"`)+`, `+str.Render(`"w"`)+`)
			  file:`+fn.Render("write")+`(`+str.Render(`"Hello World!"`)+`)
			  file:`+fn.Render("close")+`()
			  `+fn.Render("print")+`(`+str.Render(`"Done writing!"`)+`)
			`+kw.Render("end")+`
			'''

			Options:
			`+pflag.CommandLine.FlagUsages(), "\t", "")

		fmt.Println(usage)
		return
	}

	if initFlag {
		dir, _ := os.Getwd()

		_, err := os.Open("tasks.lua")
		if err == nil {
			fmt.Printf("%s: %s already exists in %s\n", errSymbol, errStyle.Render(entrypointPath), errStyle.Render(dir))
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

		fmt.Printf("%s: created %s in %s\n", lSymbol, mainStyle.Render(entrypointPath), mainStyle.Render(dir))

		return
	}

	r := &L.Runner{}

	if err := r.Compile(entrypointPath); err != nil {
		s := err.Error()
		if strings.HasPrefix(s, "parse error") {
			s = strings.ReplaceAll(s, "parse error: ", "")
			s = strings.ReplaceAll(s, ":   parse error", "")
			fmt.Printf("%s: parse error:\n%s", errSymbol, s)
			return
		}

		if strings.HasPrefix(s, "open") {
			dir, _ := os.Getwd()
			fmt.Printf("%s: %s was not found in %s: use %s to create a new one\n",
				errSymbol,
				mainStyle.Render(entrypointPath),
				mainStyle.Render(dir),
				mainStyle.Render("L --init"))

			return
		}

		fmt.Printf("%s: error:\n%s", errSymbol, s)
		return
	}

	if listFlag {
		tasks, err := r.List()
		if err != nil {
			panic(err) // todo!
		}

		if len(tasks) == 0 {
			fmt.Printf("%s: no tasks available\n", errSymbol)
			return
		}

		sort.SliceStable(tasks, func(i, j int) bool {
			return strings.Compare(tasks[i].Name, tasks[j].Name) == -1
		})

		fmt.Printf("%s: all available tasks:", lSymbol)
		for _, task := range tasks {
			if task.Description != "" {
				fmt.Printf("\n- %s: %s", mainStyle.Render(task.Name), task.Description)
			} else {
				fmt.Printf("\n- %s %s", mainStyle.Render(task.Name), subtleStyle.Render("(description not found)"))
			}

			if len(task.Sources) != 0 {
				fmt.Printf("\n  tracks: %s", strings.Join(task.Sources, ", "))
			}
		}

		fmt.Println()
		return
	}

	taskQueue := pflag.Args()
	if len(taskQueue) == 0 {
		taskQueue = []string{"default"}
	}

	if err := tea.NewProgram(model{
		TaskQueue: taskQueue,
		Runner:    r,
	}).Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
