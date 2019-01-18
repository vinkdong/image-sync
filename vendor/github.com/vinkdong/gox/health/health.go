package health

type Health struct {
	HttpGet []HttpGetCheck `yaml:"httpGet"`
	Socket  []SocketCheck
}

type SocketCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

type HttpGetCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

func (h *Health) IsHealth() bool {
	return true
}