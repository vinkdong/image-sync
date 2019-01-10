package docker

import (
	"testing"
	"gopkg.in/yaml.v2"
	"fmt"
)

func TestReadConfig(t *testing.T) {
	cfgData := []byte(`
apiVersion: v1
sync:
  from:
    registry: registry.a.com
    username: admin
  to: registry.b.com
  names:
  - app_one
  - app_two
  - app_three
  rules:
  - name: release
    value: "^v?(\\d+.)*\\d+$"
`)
	s := &NamedSync{}
	yaml.Unmarshal(cfgData, s)
	if s.Sync.Rules[0].Value != "^v?(\\d+.)*\\d+$" {
		t.Error("unmarshall config error")
	}
	if s.Sync.From.Username != "admin" {
		t.Error("umarshall config error")
	}
}

func TestSync(t *testing.T) {
	cfgData := []byte(`
apiVersion: v1
sync:
  from:
    registry: registry.a.com
    username: admin
    password: xxxx
  to: 
    registry: registry.b.com
  names:
  - "tools/nginx"
  - app_two
  - app_three
  rules:
  - name: release
    value: "^v?(\\d+.)*\\d+$"
`)
	s := &NamedSync{}
	yaml.Unmarshal(cfgData, s)
	err := s.Sync.Do()
	fmt.Println(err)
}
