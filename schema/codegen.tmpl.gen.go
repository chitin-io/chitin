// AUTOMATICALLY GENERATED FILE. DO NOT EDIT.

package schema

var codegen = tmpl(asset{Name: "codegen.tmpl", Content: "" +
	"// {{.Warning}}\n\npackage {{.Package}}\n\nimport (\n\t\"encoding/binary\"\n\t\"io\"\n\t\"reflect\"\n\t\"unsafe\"\n\n\t\"github.com/dchest/varuint\"\n\t\"chitin.io/chitin\"\n)\n\n// use all packages to avoid errors\nvar (\n\t_ = io.ErrUnexpectedEOF\n\t_ reflect.StringHeader\n\t_ unsafe.Pointer\n\t_ = varuint.Uint64\n)\n\n{{range $ver, $msg := .Schema.Messages}}\n\n{{$type := printf \"%sV%d\" $ver.Name $ver.Version}}\nfunc New{{$type}}View(data []byte) (*{{$type}}View, error) {\n\tif len(data) < minLen{{$type}}View {\n\t\treturn nil, chitin.ErrWrongSize\n\t}\n\tview := &{{$type}}View{\n\t\tdata: data,\n\t}\n\treturn view, nil\n}\n\n{{$slots := $msg.GetSlots}}\n{{$fields := $msg.GetFields}}\n{{$slotOffsets := slotOffsets $slots}}\n{{$fieldsOffset := offsetEnd $slotOffsets}}\n\nconst (\n\tslotsLen{{$type}}View  = {{$fieldsOffset}}\n\tnumFields{{$type}}View = {{len $fields}}\n\tminLen{{$type}}View    = slotsLen{{$type}}View + 1*numFields{{$type}}View\n)\n\ntype {{$type}}View struct {\n\tdata []byte\n}\n\n{{range $slotIdx, $slot := $slots}}\nfunc (v *{{$type}}View) {{$slot.Name | ucFirst}}() {{$slot.Kind.GoType}} {\n\tdata := v.data[{{(index $slotOffsets $slotIdx).Start}}:{{(index $slotOffsets $slotIdx).Stop}}]\n\t{{$slot.Kind.GoSlotGetter}}\n}\n{{end}}\n\n{{if $fields}}\nfunc (v *{{$type}}View) Fields() (*{{$type}}ViewFields, error) {\n\tf := &{{$type}}ViewFields{}\n\toff := slotsLen{{$type}}View\n\n\t// TODO this only really implements length-prefixed fields\n\n\t{{range $fieldIdx, $field := $fields}}\nloop:\n\tl, n := varuint.Uint64(v.data[off:])\n\tif n < 0 {\n\t\treturn nil, io.ErrUnexpectedEOF\n\t}\n\tif l == 0 {\n\t\t// padding\n\t\toff += n\n\t\tgoto loop\n\t}\n\tl--\n\n\tconst maxInt = int(^uint(0) >> 1)\n\tif l > uint64(maxInt) {\n\t\t// technically, it has to be truncated because it wouldn't fit\n\t\t// in memory ;)\n\t\treturn nil, io.ErrUnexpectedEOF\n\t}\n\tli := int(l)\n\n\t// TODO prevent overflow here\n\tend := slotsLen{{$type}}View + n + li\n\tif end > len(v.data) {\n\t\treturn nil, io.ErrUnexpectedEOF\n\t}\n\n\tlow := off + n\n\thigh := low + li\n\tdata := v.data[low:high]\n\n\t{ {{$field.Kind.GoFieldPrep (printf \"f.field%s\" ($field.Name | ucFirst))}} }\n\toff = high\n\t{{end}}\n\n\treturn f, nil\n}\n\ntype {{$type}}ViewFields struct {\n\t{{range $fieldIdx, $field := $fields}}\n\tfield{{$field.Name | ucFirst}} {{$field.Kind.GoStateType}}\n\t{{end}}\n}\n\n{{range $fieldIdx, $field := $fields}}\nfunc (f *{{$type}}ViewFields) {{$field.Name | ucFirst}}() {{$field.Kind.GoType}} {\n\t{{$field.Kind.GoFieldGetter (printf \"f.field%s\" ($field.Name | ucFirst))}}\n}\n{{end}}\n\n{{end}}\n\n{{end}}\n" +
	"", etag: `"dN2YGztu9lM="`})