# agorago
根据声网restfulapi接口的go版本实现


```golang
func (d *decodeState) objects() error {
	oi := d.objectInterface()
	l := list.New()
	l.PushBack(oi)
	typename := "Obj"
	for e := l.Front(); e != nil; e = e.Next() {
		if oiface, ok := e.Value.(map[string]interface{}); ok {
			d.WriteString("\ntype " + typename + " struct{\n")
			for k, v := range oiface {
				vr := reflect.ValueOf(v)
				kind := vr.Kind().String()
				switch vr.Kind() {
				case reflect.Map:
					l.PushBack(v)
					kind = getFieldsName(k, "_", FirstToUpper)
					typename = kind
				case reflect.Slice:
					kind = "[]"
					// ai := d.arrayInterface()
				default:
				}
				upper := getFieldsName(k, "_", FirstToUpper)
				lower := getFieldsName(k, "_", FirstToLower)
				line := "\t" + upper + "\t" + kind + "\t`" + `json:"` + lower + `"` + "`\n"
				d.WriteString(line)
			}
			d.WriteString("}\n")
		}
	}
	fmt.Println(d.String())
	return nil
}
```