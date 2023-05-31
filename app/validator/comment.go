package validator

type Comment struct{
	Pid     int    `json:"pid" rule:"number"`
	Content string `json:"content" rule:"required"`
}

var CommentMessage = map[string]string{
	"pid.number"	  : "pid 必须为数字！",
	"content.required": "content 不能为空！",
}

func (this Comment) Message() map[string]string {
	return CommentMessage
}

func (this Comment) Struct() any {
	return this
}