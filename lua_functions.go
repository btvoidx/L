package L

import (
	"runtime"

	lua "github.com/yuin/gopher-lua"
)

func luaNoop(l *lua.LState) int { return 0 }

func luaSetTaskDescription(s *string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		*s = L.CheckString(1)
		return 0
	}
}

func luaSetTaskSources(arr *[]string) func(*lua.LState) int {
	return func(L *lua.LState) int {
		t := L.CheckTable(1)
		t.ForEach(func(k, v lua.LValue) {
			if k.Type() != lua.LTNumber {
				L.RaiseError("all keys must be of type number")
				return
			}

			if v.Type() != lua.LTString {
				L.RaiseError("all values must be of type string")
				return
			}

			*arr = append(*arr, v.String())
		})
		return 0
	}
}

func luaGetOs(L *lua.LState) int {
	L.Push(lua.LString(runtime.GOOS))
	return 1
}
