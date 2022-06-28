package L

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/muesli/termenv"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"

	"github.com/btvoidx/L/internal/logger"
)

type Executor struct {
	Entrypoint    string
	Logger        *logger.Logger
	fnproto       *lua.FunctionProto
	taskinfoCache []TaskMeta
}

type TaskMeta struct {
	Name         string
	Description  string
	Dependencies []string
	Sources      []string
}

// This timeout is highly debatable.
// To be fair, anything under about 80ms will feel snappy for most people.
const initTimeout = 60 * time.Millisecond

func noop(l *lua.LState) int { return 0 }

// Parses tasks script and compiles it for future use.
func (e *Executor) Compile() error {
	file, err := os.Open(e.Entrypoint)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, e.Entrypoint)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	proto, err := lua.Compile(chunk, e.Entrypoint)
	if err != nil {
		return fmt.Errorf("compile error: %w", err)
	}

	e.fnproto = proto
	return nil
}

func (e *Executor) loadScript(L *lua.LState, trimUnsafePackages bool) error {
	e.Logger.WriteEphemeral("L: loading script")

	L.G.Global.RawSetString("task", &lua.LTable{})

	// todo!
	// if trimUnsafePackages {
	// 	var unsafePackages = []string{"assert"}

	// 	L.G.Global.ForEach(func(k, v lua.LValue) {
	// 		// if v.Type() == lua.LTTable {

	// 		// }

	// 		for _, v := range safePackages {
	// 			if k.String()
	// 		}
	// 		L.G.Global.RawSetString()
	// 	})
	// }

	L.Push(L.NewFunctionFromProto(e.fnproto))
	err := L.PCall(0, lua.MultRet, nil)
	if err != nil {
		return err
	}

	return nil
}

// Runs a given task
func (e *Executor) Run(taskName string) (code int, err error) {
	taskList, err := e.List()
	if err != nil {
		return 1, err
	}

	// If taskName is not in taskList
	if func() bool {
		for _, t := range taskList {
			if taskName == t.Name {
				return false
			}
		}
		return true
	}() {
		e.Logger.Err("L: task %s not found", termenv.String(taskName))
		return 1, nil
	}

	L := lua.NewState()
	defer L.Close()

	if err := e.loadScript(L, false); err != nil {
		return 1, err
	}

	e.Logger.Write("L: running %s", termenv.String(taskName))

	L.SetFuncs(L.G.Global, map[string]lua.LGFunction{
		"description": noop,
		"depends":     noop,
		"defer":       noop,
		"sources":     noop,
	})

	// Guaranteed correct data types, checked by call to List() above
	tasksTable, _ := L.G.Global.RawGetString("task").(*lua.LTable)
	fn, _ := tasksTable.RawGetString(taskName).(*lua.LFunction)

	if err := L.CallByParam(lua.P{
		Fn:      fn,
		Protect: true,
	}); err != nil {
		return 1, err
	}

	return 0, nil
}

// Returns all tasks from loaded script, by running it in safe-ish mode with harsh-ish timeout.
func (e *Executor) List() ([]TaskMeta, error) {
	if e.taskinfoCache != nil && len(e.taskinfoCache) != 0 {
		return e.taskinfoCache, nil
	}

	L := lua.NewState()
	defer L.Close()

	ctx, cancel := context.WithTimeout(context.Background(), initTimeout)
	defer cancel()
	L.SetContext(ctx)

	err := e.loadScript(L, true)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			err = fmt.Errorf("script took too long to run (>%s); make sure it's not doing heavy computations outside of functions",
				initTimeout)
		}
		return []TaskMeta{}, err
	}

	var tasksTable *lua.LTable
	var ok bool
	if tasksTable, ok = L.G.Global.RawGetString("task").(*lua.LTable); !ok {
		return nil, fmt.Errorf("global 'task' is of type %s; expected table", L.G.Global.RawGetString("task").Type().String())
	}

	tasks := make([]TaskMeta, 0, tasksTable.Len())
	tasksTable.ForEach(func(k, v lua.LValue) {
		if k.Type() != lua.LTString || v.Type() != lua.LTFunction {
			return
		}

		tasks = append(tasks, TaskMeta{
			Name: k.String(),
		})
	})

	e.taskinfoCache = tasks

	return tasks, nil
}
