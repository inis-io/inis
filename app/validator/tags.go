package validator

type Tags struct{}

var TagsMessage = map[string]string{}

func (this Tags) Message() map[string]string {
	return TagsMessage
}

func (this Tags) Struct() any {
	return this
}
