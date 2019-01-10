package docker

import (
	"testing"
	"fmt"
)

func TestPullImage(t *testing.T) {
	d := Docker{
		Registry: "docker.io/library",
	}
	fmt.Println(d.pullImage("centos", "latest"))
}

func TestListTags(t *testing.T) {
	d := Docker{
		Registry: "registry.cn-shanghai.aliyuncs.com",
	}
	image, err:= d.listTags("vinkdong/nfs-provisioner")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(image.Tags)
}
