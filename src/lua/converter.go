package lua

import (
	glua "github.com/yuin/gopher-lua"
)

func luaTableToMap(table *glua.LTable) map[string]any {
	result := make(map[string]any)

	table.ForEach(func(key, value glua.LValue) {
		if k, ok := key.(glua.LString); ok {
			result[string(k)] = luaValueToGo(value)
		}
	})

	return result
}

func luaValueToGo(value glua.LValue) any {
	switch v := value.(type) {
	case *glua.LTable:
		if n := v.Len(); n > 0 && isArray(v, n) {
			result := make([]any, n)
			for i := range result {
				result[i] = luaValueToGo(v.RawGetInt(i + 1))
			}
			return result
		}
		return luaTableToMap(v)

	case glua.LString:
		return string(v)
	case glua.LNumber:
		return float64(v)
	case glua.LBool:
		return bool(v)
	case *glua.LNilType:
		return nil
	default:
		return v.String()
	}
}

func isArray(table *glua.LTable, n int) bool {
	count := 0

	for k, _ := table.Next(glua.LNil); k != glua.LNil; k, _ = table.Next(k) {
		num, ok := k.(glua.LNumber)
		if !ok {
			return false
		}

		i := int(num)
		if float64(i) != float64(num) || i < 1 || i > n {
			return false
		}
		count++
	}

	return count == n
}
