package validator

type ApiKeys struct{}

var ApiKeysMessage = map[string]string{}

func (this ApiKeys) Message() map[string]string {
	return ApiKeysMessage
}

func (this ApiKeys) Struct() any {
	return this
}
