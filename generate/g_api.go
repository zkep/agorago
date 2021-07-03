package generate

var tpl = `package  agorago

import ()

type {{projectName}} struct {
	req *Request
}

type {{projectName}}Option func(c *{{projectName}})

func Add{{projectName}}Request(req *Request) {{projectName}}Option {
	return func(c *{{projectName}}) {
		c.req = req
	}
}

type New{{projectName}}(opts ...{{projectName}}Option) struct {
	r := &{{projectName}}{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (self *{{projectName}}) {{action}}(req {{action}}Request, ret *{{action}}Response) error {
	uri := fmt.Sprintf(CLOUD_RECORDING_ACQUIRE_URL, self.req.appid)
	err := self.req.Do(uri, http.Method{{method}}, req, nil, nil, ret)
	if err != nil {
		return err
	}
	return nil
}


`
