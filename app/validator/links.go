package validator

type Links struct{}

var LinksMessage = map[string]string{}

func (this Links) Message() map[string]string {
	return LinksMessage
}

func (this Links) Struct() any {
	return this
}
