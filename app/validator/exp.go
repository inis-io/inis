package validator

type EXP struct{}

var EXPMessage = map[string]string{}

func (this EXP) Message() map[string]string {
	return EXPMessage
}

func (this EXP) Struct() any {
	return this
}