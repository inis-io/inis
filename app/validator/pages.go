package validator

type Pages struct {
	Key string `json:"key" rule:"required,alphaDash"`
}

var PagesMessage = map[string]string{
	"key.required":  "唯一键不能为空",
	"key.alphaDash": "唯一键只能由字母、数字、破折号（ - ）以及下划线（ _ ）组成",
}

func (this Pages) Message() map[string]string {
	return PagesMessage
}

func (this Pages) Struct() any {
	return this
}
