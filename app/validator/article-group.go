package validator

type ArticleGroup struct{
	Pid int `json:"pid" rule:"number"`
}

var ArticleGroupMessage = map[string]string{
	"pid.number":    "pid 只能是数字！",
}

func (this ArticleGroup) Message() map[string]string {
	return ArticleGroupMessage
}

func (this ArticleGroup) Struct() any {
	return this
}