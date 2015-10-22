package schema

import "fmt"

type String struct{}

var _ FieldType = String{}

func (String) isField() {}

func (String) GoType() string {
	return `string`
}

func (String) GoStateType() string {
	return `string`
}

func (String) GoFieldPrep(stateVar string) string {
	return fmt.Sprintf(`
	p := (*reflect.StringHeader)(unsafe.Pointer(&%s))
	p.Data = uintptr(unsafe.Pointer(&data[0]))
	p.Len = len(data)
`, stateVar)
}

func (String) GoFieldGetter(stateVar string) string {
	return fmt.Sprintf(`return %s`, stateVar)
}
