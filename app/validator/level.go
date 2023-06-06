package validator

type Level struct{}

var LevelMessage = map[string]string{}

func (this Level) Message() map[string]string {
	return LevelMessage
}

func (this Level) Struct() any {
	return this
}
