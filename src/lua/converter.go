package lua

import (
	glua "github.com/yuin/gopher-lua"
)

func toGo(value glua.LValue) any {
	switch v := value.(type) {
	case glua.LString:
		return string(v)
	case glua.LNumber:
		return float64(v)
	case glua.LBool:
		return bool(v)
	case *glua.LTable:
		return tableToGo(v)
	case *glua.LNilType:
		return nil
	default:
		return v.String()
	}
}

func toMap(t *glua.LTable) map[string]any {
	m := make(map[string]any)

	t.ForEach(func(k, v glua.LValue) {
		if s, ok := k.(glua.LString); ok {
			m[string(s)] = toGo(v)
		}
	})

	return m
}

func tableToGo(t *glua.LTable) any {
	n := t.Len()

	var arr []any
	if n > 0 {
		arr = make([]any, n)
	}

	var m map[string]any
	ints, mixed := 0, false

	t.ForEach(func(k, v glua.LValue) {
		switch k := k.(type) {
		case glua.LString:
			if m == nil {
				m = make(map[string]any)
			}
			m[string(k)] = toGo(v)
		case glua.LNumber:
			i := int(k)
			if float64(i) == float64(k) && i >= 1 && i <= n {
				arr[i-1] = toGo(v)
				ints++
			} else {
				mixed = true
			}
		default:
			mixed = true
		}
	})

	if n > 0 && ints == n && m == nil && !mixed {
		return arr
	}
	if m == nil {
		return map[string]any{}
	}
	return m
}
