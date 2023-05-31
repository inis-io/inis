package validator

type LinksGroup struct{}

var LinksGroupMessage = map[string]string{}

func (this LinksGroup) Message() map[string]string {
	return LinksGroupMessage
}

func (this LinksGroup) Struct() any {
	return this
}
