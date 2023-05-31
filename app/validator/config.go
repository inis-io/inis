package validator

type Config struct{
	Key     string `json:"key" rule:"required"`
}

var ConfigMessage = map[string]string{
	"key.required": "key 不能为空！",
}

func (this Config) Message() map[string]string {
	return ConfigMessage
}

func (this Config) Struct() any {
	return this
}
