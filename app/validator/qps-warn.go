package validator

type QpsWarn struct{}

var QpsWarnMessage = map[string]string{}

func (this QpsWarn) Message() map[string]string {
	return QpsWarnMessage
}

func (this QpsWarn) Struct() any {
	return this
}

