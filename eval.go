package stone

import (
	"github.com/mohae/deepcopy"
	"os"
)

type (
	VarMap     map[string]interface{}
	SectionMap map[string]*VarMap
)

func (stone *Stone) eval() *SectionMap {
	sections := SectionMap{}
	for _, s := range stone.Sections {
		sections[s.Name] = s.eval(&sections)
	}
	return &sections
}

func (section *Section) eval(s *SectionMap) *VarMap {
	vm := VarMap{}
	if section.ParentName != nil {
		vm = *(deepcopy.Copy((*s)[*section.ParentName]).(*VarMap))
	}
	for _, stmt := range section.Stmts {
		stmt.eval(&vm)
	}
	return &vm
}

func (stmt *Stmt) eval(vm *VarMap) {
	if stmt.Deletion != nil {
		stmt.Deletion.delete(vm)
		return
	}
	stmt.Assignment.eval(vm)
}

func (assign *Assignment) eval(vm *VarMap) {
	value := assign.Right.eval(vm)
	assign.Left.assign(vm, value)
}

func (lv *LeftValue) onlyIdent() bool {
	return lv.KeyString == nil && lv.KeyInt == nil
}

func (lv *LeftValue) delete(vm *VarMap) {
	if lv.onlyIdent() {
		delete(*vm, lv.Identifier)
		return
	}
	realV := (*vm)[lv.Identifier]
	switch v := realV.(type) {
	case map[string]interface{}:
		delete(v, *lv.KeyString)
	case []interface{}:
		i := *lv.KeyInt
		(*vm)[lv.Identifier] = append(v[:i], v[i+1:]...)
	default:
		panic("????")
	}
}

func (lv *LeftValue) assign(vm *VarMap, value interface{}) {
	if lv.onlyIdent() {
		(*vm)[lv.Identifier] = value
		return
	}
	realV := (*vm)[lv.Identifier]
	switch v := realV.(type) {
	case map[string]interface{}:
		v[*lv.KeyString] = value
	case []interface{}:
		v[*lv.KeyInt] = value
	default:
		panic("???")
	}
}

func (lv *LeftValue) eval(vm *VarMap) interface{} {
	if lv.onlyIdent() {
		return (*vm)[lv.Identifier]
	}
	realV := (*vm)[lv.Identifier]
	switch v := realV.(type) {
	case map[string]interface{}:
		return v[*lv.KeyString]
	case []interface{}:
		return v[*lv.KeyInt]
	default:
		panic("??")
	}
}

func (rv *RightValue) eval(vm *VarMap) interface{} {
	if rv.String != nil {
		return *rv.String
	} else if rv.Int != nil {
		return *rv.Int
	} else if rv.Float != nil {
		return *rv.Float
	} else if rv.Bool != nil {
		return *rv.Bool
	} else if rv.LeftValue != nil {
		return rv.LeftValue.eval(vm)
	} else if rv.ArrayValue != nil {
		return rv.ArrayValue.eval(vm)
	} else if rv.MapValue != nil {
		return rv.MapValue.eval(vm)
	} else if rv.EnvValue != nil {
		return rv.EnvValue.eval()
	}
	return nil
}

func (av *ArrayValue) eval(vm *VarMap) []interface{} {
	result := []interface{}{}
	for _, rv := range av.Array {
		value := rv.eval(vm)
		result = append(result, value)
	}
	return result
}

func (mv *MapValue) eval(vm *VarMap) map[string]interface{} {
	result := map[string]interface{}{}
	for _, kv := range mv.KVs {
		key := kv.K
		value := kv.V.eval(vm)
		result[key] = value
	}
	return result
}

func (ev *EnvValue) eval() string {
	return os.Getenv(ev.Env)
}
