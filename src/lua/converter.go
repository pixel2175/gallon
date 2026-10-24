package lua

import	glua "github.com/yuin/gopher-lua"

func luaTableToMap(table *glua.LTable) map[string]any {
	result := make(map[string]any)

	table.ForEach(func(key, value glua.LValue) {
		if key.Type() != glua.LTString {
			return
		}

		result[key.String()] = luaValueToGo(value)
	})

	return result
}

func luaValueToGo(value glua.LValue) any {
	switch value := value.(type) {
	case *glua.LTable:
		return luaTableToMap(value)
	case glua.LString:
		return string(value)
	case glua.LNumber:
		return float64(value)
	case glua.LBool:
		return bool(value)
	case *glua.LNilType:
		return nil
	default:
		return value.String()
	}
}
