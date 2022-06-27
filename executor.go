package L

import (
	"bufio"
	"fmt"
	"os"

	"github.com/btvoidx/L/internal/color"
	"github.com/btvoidx/L/internal/logger"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type Executor struct {
	FilePath      string
	Logger        *logger.Logger
	fnproto       *lua.FunctionProto
	tasknameCache []string
}

func noop(l *lua.LState) int { return 0 }

func getTaskFn(L *lua.LState, task string) (*lua.LFunction, error) {
	lv := L.G.Global.RawGetString("task")
	tasks, ok := lv.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("global 'task' is not a table")
	}

	lv = tasks.RawGetString(task)
	fn, ok := lv.(*lua.LFunction)
	if !ok {
		if lv.Type() == lua.LTNil {
			return nil, fmt.Errorf("task %s is not defined", task)
		}
		return nil, fmt.Errorf("global task.%s is not a function", task)
	}

	return fn, nil
}

func getTasks(L *lua.LState) (tasks []string, err error) {
	lv := L.G.Global.RawGetString("task")
	tasksTable, ok := lv.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("global 'task' is of type %s; expected table", lv.Type().String())
	}

	tasks = make([]string, 0, tasksTable.Len())
	tasksTable.ForEach(func(k, v lua.LValue) {
		if k.Type() != lua.LTString || v.Type() != lua.LTFunction {
			return
		}

		tasks = append(tasks, k.String())
	})

	return tasks, nil
}

func (e *Executor) Compile() error {
	file, err := os.Open(e.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, e.FilePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	proto, err := lua.Compile(chunk, e.FilePath)
	if err != nil {
		return fmt.Errorf("compile error: %w", err)
	}

	e.fnproto = proto
	return nil
}

func (e *Executor) loadScript(L *lua.LState) error {
	e.Logger.WriteEphemeral("L: loading script")

	L.G.Global.RawSetString("task", &lua.LTable{})
	L.G.Global.RawSetString("description", L.NewFunction(noop))
	L.G.Global.RawSetString("depends", L.NewFunction(noop))
	L.G.Global.RawSetString("defer", L.NewFunction(noop))
	L.G.Global.RawSetString("sources", L.NewFunction(noop))

	L.Push(L.NewFunctionFromProto(e.fnproto))
	err := L.PCall(0, lua.MultRet, nil)
	if err != nil {
		return err
	}

	return nil
}

func (e *Executor) Run(task string) (code int, err error) {
	L := lua.NewState()

	if err := e.loadScript(L); err != nil {
		return 1, err
	}

	fn, err := getTaskFn(L, task)
	if err != nil {
		return 1, err
	}

	e.Logger.Write("L: running '%s'", color.Magenta(task))

	if err := L.CallByParam(lua.P{
		Fn:      fn,
		Protect: true,
	}); err != nil {
		return 1, err
	}

	return 0, nil
}

func (e *Executor) List() ([]string, error) {
	if e.tasknameCache != nil && len(e.tasknameCache) != 0 {
		return e.tasknameCache, nil
	}

	L := lua.NewState()

	err := e.loadScript(L)
	if err != nil {
		return []string{}, err
	}

	tasks, err := getTasks(L)
	if err != nil {
		return []string{}, err
	}

	return tasks, nil
}
