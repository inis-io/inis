package validator

type Article struct{
	Top int `json:"top" rule:"number,min:0,max:1"`
}

var ArticleMessage = map[string]string{
	"top.number":    "top 只能是数字！",
	"top.min":       "top 只能是0或1！",
	"top.max":       "top 只能是0或1！",
}

func (this Article) Message() map[string]string {
	return ArticleMessage
}

func (this Article) Struct() any {
	return this
}

