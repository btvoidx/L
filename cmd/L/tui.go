package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/btvoidx/L"
)

type taskRunResult string

var (
	taskRunComplete taskRunResult = "success"
	taskRunError    taskRunResult = "error"
	taskRunSkipped  taskRunResult = "skipped"
)

type taskCompleteMsg struct {
	result taskRunResult
}

type model struct {
	TaskQueue []string
	Runner    *L.Runner

	currentTask int
	taskResults []taskRunResult
	done        bool
}

func (m model) Init() tea.Cmd {
	return runNextTask()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case taskCompleteMsg:
		m.currentTask++
		m.taskResults = append(m.taskResults, msg.result)

		batch := make([]tea.Cmd, 0, 3)
		batch = append(batch, tea.Printf("%s: task complete: %s", lSymbol, mainStyle.Render(m.TaskQueue[m.currentTask-1])))

		if m.currentTask >= len(m.TaskQueue) {
			m.done = true

			batch = append(batch, tea.Quit)
			return m, tea.Batch(batch...)
		}

		batch = append(batch, runNextTask())
		return m, tea.Batch(batch...)
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return ""
	}

	var s strings.Builder

	s.WriteString(lSymbol + ":")

	for cn, task := range m.TaskQueue {
		s.WriteString(" ")

		if m.currentTask == cn {
			s.WriteString(mainStyle.Render(task))
			continue
		}

		if m.currentTask < cn {
			s.WriteString(subtleStyle.Render(task))
			continue
		}

		switch m.taskResults[cn] {
		case taskRunComplete:
			s.WriteString(successStyle.Render(task))
		case taskRunError:
			s.WriteString(errStyle.Render(task))
		default:
			s.WriteString(task)
		}
	}

	return s.String()
}

func runNextTask() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return taskCompleteMsg{result: taskRunComplete}
	})
}
