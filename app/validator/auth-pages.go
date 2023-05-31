package validator

type AuthPages struct {
	Name string `json:"name" rule:"required"`
	Path string `json:"path" rule:"required"`
}

var AuthPagesMessage = map[string]string{
	"name.required": "名称不能为空！",
	"path.required": "路径不能为空！",
}

func (this AuthPages) Message() map[string]string {
	return AuthPagesMessage
}

func (this AuthPages) Struct() any {
	return this
}
