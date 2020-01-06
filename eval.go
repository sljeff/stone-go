package stone

import (
	"github.com/mohae/deepcopy"
)

type (
	varMap     map[string]interface{}
	sectionMap map[string]*varMap
)

func (stone *Stone) eval() *sectionMap {
	sections := sectionMap{}
	for _, s := range stone.Sections {
		sections[s.Name] = s.eval(&sections)
	}
	return &sections
}

func (section *Section) eval(s *sectionMap) *varMap {
	vm := varMap{}
	if section.ParentName != nil {
		vm = deepcopy.Copy((*s)[*section.ParentName]).(varMap)
	}
	for _, stmt := range section.Stmts {
		name, value := stmt.eval(&vm)
		if name == "" {
			continue
		}
		vm[name] = value
	}
	return &vm
}

func (stmt *Stmt) eval(vm *varMap) (string, interface{}) {
	if stmt.Deletion != nil {
		stmt.Deletion.delete(vm)
		return "", nil
	}
	return stmt.Assignment.eval(vm)
}

func (assign *Assignment) eval(vm *varMap) (string, interface{}) {
	return "", nil
}

func (lv *LeftValue) delete(vm *varMap) {
	if lv.KeyString == nil && lv.KeyInt == nil {
		delete(*vm, lv.Identifier)
	}
}

func (lv *LeftValue) assign(vm *varMap, value interface{}) {
	if lv.KeyString == nil && lv.KeyInt == nil {
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
