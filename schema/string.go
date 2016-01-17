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
	msg, next, err := chitinParseLengthPrefixed(data)
	if err != nil {
		return nil, err
	}
	p := (*reflect.StringHeader)(unsafe.Pointer(&%s))
	p.Data = uintptr(unsafe.Pointer(&msg[0]))
	p.Len = len(msg)
	data = next
`, stateVar)
}

func (String) GoFieldGetter(stateVar string) string {
	return fmt.Sprintf(`return %s`, stateVar)
}

func (String) GoFieldBytes(stateVar string) string {
	return fmt.Sprintf(`
	var lb [varuint.MaxUint64Len]byte
	var ll int

	ll = varuint.PutUint64(lb[:], uint64(len(%[1]s))+1)
	data = append(data, lb[:ll]...)
	data = append(data, %[1]s...)
`, stateVar)
}
