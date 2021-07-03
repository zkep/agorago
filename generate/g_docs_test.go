package generate

import (
	"os"
	"testing"

	yaml "gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

func TestParseYaml(t *testing.T) {
	fd, err := os.Open("./server_restfulapi_cn.yaml")
	assert.NoError(t, err)
	defer fd.Close()
	assert.NoError(t, err)
	var out Swagger
	err = yaml.NewDecoder(fd).Decode(&out)
	assert.NoError(t, err)
	t.Logf("%+v", out.Servers)
	for k, item := range out.Paths {
		t.Logf("%s", k)
		t.Logf("%+v", item.Post)
		t.Logf("%+v", item.Get)
	}
}
