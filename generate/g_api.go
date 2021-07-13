package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

func GenerateYamlAPI(yamlpath string) {
	fd, err := os.Open(yamlpath)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	var out Swagger
	err = yaml.NewDecoder(fd).Decode(&out)
	if err != nil {
		panic(err)
	}
	api := &API{}
	opts := make(map[string][]OptItem)
	for path, items := range out.Paths {
		rankOperation(opts, path, http.MethodPost, items.Post)
		rankOperation(opts, path, http.MethodGet, items.Get)
		rankOperation(opts, path, http.MethodHead, items.Head)
		rankOperation(opts, path, http.MethodDelete, items.Delete)
		rankOperation(opts, path, http.MethodPut, items.Put)
		rankOperation(opts, path, http.MethodOptions, items.Options)
	}
	var buf bytes.Buffer
	for _, opt := range opts {
		buf.Reset()
		api.actions = []Action{}
		api.name = getStructName(opt...)
		for _, v := range opt {
			actionName := FirstToUpper(strings.ToLower(v.method)) + getActionName(v.path)
			paramName := actionName + "Req"
			parameCode := ParameStructCode(paramName, v.Parameters...)
			resp := v.Responses["200"]
			responseName := actionName + "Resp"
			responseCode := ""
			for _, ctx := range resp.Content {
				raw, err := json.Marshal(ctx.Example)
				if err != nil {
					panic(err)
				}
				responseCode = Json2StructCode(raw, responseName)
				break
			}
			action := Action{
				ApiName:     api.name,
				Name:        actionName,
				Method:      v.method,
				Path:        v.path,
				Summary:     v.Summary,
				OperationID: v.OperationID,
				Parameters:  []ActionParam{{Name: paramName, Code: parameCode}},
				Responses:   []ActionParam{{Name: responseName, Code: responseCode}},
			}
			api.AddAction(action)
		}
		err = tplAPI.Execute(&buf, api)
		if err != nil {
			panic(err)
		}
		code := api.packageGoCode() + api.importsGoCode() + strings.TrimSpace(buf.String())
		os.WriteFile("../"+FirstToLower(api.name)+".go", []byte(code), 0777)
	}
}

type API struct {
	name        string
	summary     string
	packageName string
	imports     map[string]bool
	actions     []Action
}

type Action struct {
	ApiName     string
	Name        string
	Method      string
	Path        string
	Summary     string
	OperationID string
	Parameters  []ActionParam
	Responses   []ActionParam
}

type ActionParam struct {
	Code string
	Name string
}

func (a *API) PackageName() string {
	if len(a.packageName) == 0 {
		a.packageName = "agorago"
	}
	return strings.ToLower(a.packageName)
}

func (a *API) Summary() string {
	return a.summary
}

func (a *API) StructName() string {
	return a.name
}

func (a *API) AddImport(v string) error {
	a.imports[v] = true
	return nil
}

func (a *API) resetImports() {
	a.imports = map[string]bool{}
}

func (a *API) AddAction(v ...Action) {
	a.actions = append(a.actions, v...)
}

func (a *API) Actions() []Action {
	return a.actions
}

func (a *API) importsGoCode() string {
	if len(a.imports) == 0 {
		return ""
	}
	corePkgs, extPkgs := []string{}, []string{}
	for i := range a.imports {
		if strings.Contains(i, ".") {
			extPkgs = append(extPkgs, i)
		} else {
			corePkgs = append(corePkgs, i)
		}
	}
	sort.Strings(corePkgs)
	sort.Strings(extPkgs)

	code := "import (\n"
	for _, i := range corePkgs {
		code += fmt.Sprintf("\t%q\n", i)
	}
	if len(corePkgs) > 0 {
		code += "\n"
	}
	for _, i := range extPkgs {
		code += fmt.Sprintf("\t%q\n", i)
	}
	code += ")\n\n"
	return code
}

func (a *API) packageGoCode() string {
	code := "package " + a.PackageName() + "\n\n"
	return code
}

type OptItem struct {
	*Operation
	method string
	path   string
}

func rankOperation(opts map[string][]OptItem, path, method string, item *Operation) {
	if item == nil {
		return
	}
	for _, tag := range item.Tags {
		opt, ok := opts[tag]
		if !ok {
			if opt == nil {
				opts[tag] = make([]OptItem, 0, 1)
			}
		}
		opts[tag] = append(opts[tag], OptItem{
			Operation: item,
			method:    method,
			path:      path,
		})
	}
}

func getStructName(opts ...OptItem) string {
	if opts == nil {
		return ""
	}
	paths := make([]string, 0, len(opts))
	for _, opt := range opts {
		paths = append(paths, opt.path)
	}
	if len(paths) == 1 {
		path := paths[0]
		fold := EqualFold(path, path, "/", 2)
		sp := strings.Split(fold, "-")
		name := ""
		for _, s := range sp {
			name += FirstToUpper(s)
		}
		return name
	}
	dest := make([]string, 0, len(opts))
	folds := make(map[string]int, len(opts))
	for k, v := range paths {
		dest = dest[:0]
		dest = append(dest, paths[0:k]...)
		dest = append(dest, paths[k+1:]...)
		for _, d := range dest {
			fold := EqualFold(d, v, "/", 2)
			if len(fold) > 0 {
				if f, ok := folds[fold]; ok {
					f++
					folds[fold] = f
					continue
				}
				folds[fold] = 1
			}
		}
	}
	if len(folds) > 0 {
		max := 0
		btter := ""
		for k, v := range folds {
			if v > max {
				max = v
				btter = k
			}
		}
		sp := strings.Split(btter, "-")
		name := ""
		for _, s := range sp {
			name += FirstToUpper(s)
		}
		return name
	}
	return ""
}

func getActionName(path string) string {
	if len(path) <= 0 {
		return "AotuGenerate"
	}
	fold := EqualFold(path, path, "/", 2)
	if strings.Index(fold, "/{") > 0 {
		fold = fold[:strings.Index(fold, "/{")]
	} else if strings.Index(fold, "/") > 0 {
		fold = fold[:strings.Index(fold, "/")]
	}
	fold = strings.ReplaceAll(fold, "/", "-")
	return getFieldsName(getFieldsName(fold, "-", FirstToUpper), "_", FirstToUpper)
}

func getFieldsName(fileName, split string, fn func(s string) string) string {
	sp := strings.Split(fileName, split)
	var name string
	for _, s := range sp {
		name += fn(s)
	}
	return name
}

func FirstToUpper(s string) string {
	if s != "" {
		var sr rune
		var ss string
		if s[0] < utf8.RuneSelf {
			sr, ss = rune(s[0]), s[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s)
			sr, ss = r, s[size:]
		}
		r := unicode.SimpleFold(sr)
		for r != sr {
			r = unicode.SimpleFold(r)
		}
		return strings.ToUpper(string(sr)) + ss
	}
	return s
}

func FirstToLower(s string) string {
	if s != "" {
		var sr rune
		var ss string
		if s[0] < utf8.RuneSelf {
			sr, ss = rune(s[0]), s[1:]
		} else {
			r, size := utf8.DecodeRuneInString(s)
			sr, ss = r, s[size:]
		}
		r := unicode.SimpleFold(sr)
		for r != sr {
			r = unicode.SimpleFold(r)
		}
		return strings.ToLower(string(sr)) + ss
	}
	return s
}

func EqualFold(s, t, split string, step int) (fold string) {
	ss := strings.Split(strings.Trim(s, split), split)
	st := strings.Split(strings.Trim(t, split), split)
	src := ss
	dst := st
	if len(st) > len(ss) {
		dst = ss
		src = st
	}
	if len(dst) < step {
		return
	}
	dst = dst[step:]
	src = src[step:]
	for k, v := range dst {
		if !strings.EqualFold(src[k], v) {
			break
		}
		fold += split + v
	}
	if len(fold) > 0 {
		fold = strings.TrimPrefix(fold, split)
	}
	return
}

type Parame struct {
	Name        string
	Tag         string
	Description string
	Type        string
}

type ParameStruct struct {
	paramesName string
	parames     []Parame
}

func (p *ParameStruct) ParamesName() string {
	return p.paramesName
}

func (p *ParameStruct) Parame() []Parame {
	return p.parames
}

func ParameStructCode(structName string, params ...Parameter) string {
	p := &ParameStruct{}
	p.paramesName = structName
	for _, param := range params {
		description := strings.ReplaceAll(param.Description, "\n", "")
		fieldName := getFieldsName(param.Name, "_", FirstToUpper)
		tag := fmt.Sprintf("`"+`json:"%s"`+"`", param.Name)
		if param.Schema.Type == "number" {
			param.Schema.Type = "int"
		}
		p.parames = append(p.parames, Parame{
			Name:        fieldName,
			Tag:         tag,
			Description: description,
			Type:        param.Schema.Type,
		})
	}
	var buf bytes.Buffer
	err := tplParames.Execute(&buf, p)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

type decodeState struct {
	bytes.Buffer
}

func (d *decodeState) Map2StructCode(m map[string]interface{}, depth int) {
	for k, v := range m {
		if v == nil {
			continue
		}
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Map: // 对象类型
			d.WriteString(strings.Repeat("\t", depth) + FirstToUpper(k) + "\tstruct{\t\n")
			d.Map2StructCode(reflect.ValueOf(v).Interface().(map[string]interface{}), depth+1)
			d.WriteString(fmt.Sprintf(strings.Repeat("\t", depth)+`} %sjson:"%s"%s`, "`", FirstToLower(k), "`") + "\n")
			continue
		case reflect.Slice: // 切片
			kind := "[]"
			array := reflect.ValueOf(v).Interface().([]interface{})
			if len(array) > 0 {
				vf := array[0]
				vfkind := reflect.TypeOf(vf).Kind()
				switch vfkind {
				case reflect.Map:
					d.WriteString(strings.Repeat("\t", depth) + FirstToUpper(k) + "\tstruct{\t\n")
					d.Map2StructCode(reflect.ValueOf(vf).Interface().(map[string]interface{}), depth+1)
					d.WriteString(fmt.Sprintf(strings.Repeat("\t", depth)+`} %sjson:"%s"%s`, "`", FirstToLower(k), "`") + "\n")
					continue
				default:
					kind += vfkind.String()
				}
			}
			d.WriteString(strings.Repeat("\t", depth) + FirstToUpper(k) + "\t" + kind + "\t" + fmt.Sprintf(`%sjson:"%s"%s`, "`", FirstToLower(k), "`") + "\n")
		default:
			d.WriteString(strings.Repeat("\t", depth) + FirstToUpper(k) + "\t" + t.Kind().String() + "\t" + fmt.Sprintf(`%sjson:"%s"%s`, "`", FirstToLower(k), "`") + "\n")
		}
	}
}

func Json2StructCode(js []byte, structName string) string {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(js), &m)
	if err != nil {
		panic(fmt.Sprintf("%s,%s", err.Error(), string(js)))
	}
	var d decodeState
	if len(structName) == 0 {
		structName = "Object"
	}
	d.WriteString("type " + structName + " struct{\n")
	d.Map2StructCode(m, 1)
	d.WriteString("}\n")
	return d.String()
}

var tplParames = template.Must(template.New("agorago").Parse(`
type {{ .ParamesName }}  struct {
	{{- range $_, $o := .Parame }}
	{{ $o.Name }}  {{ $o.Type }}  {{ $o.Tag }}   // {{ $o.Description }}
	{{- end }}
}
`))

var tplAPI = template.Must(template.New("agorago").Funcs(template.FuncMap{}).Parse(`
type {{ .StructName }} struct {
	*Request
}

type {{ .StructName }}Option func(c *{{ .StructName }})

func Add{{ .StructName }}Request(req *Request) {{ .StructName }}Option  {
	return func(c *{{ .StructName }}) {
		c.Request = req
	}
}

func New{{ .StructName }}(opts ...{{ .StructName }}Option) *{{ .StructName }} {
	r := &{{ .StructName }}{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}


{{- range $_, $o := .Actions }}

{{- range $_, $p := $o.Parameters }}
{{ $p.Code }}
{{- end }}

{{- range $_, $p := $o.Responses }}
{{ $p.Code }}
{{- end }}

// {{ $o.Summary }}
func (self *{{ $o.ApiName }}) {{ $o.Name }}({{- range $_, $p := $o.Parameters }}req {{ $p.Name }},{{- end }} {{- range $_, $p := $o.Responses }} ret *{{ $p.Name }}{{- end }}) error {
	err := self.Do("{{ $o.Path }}", "{{ $o.Method }}", req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}
{{- end }}

`))
