package validator

type AuthRules struct {
	Common int `json:"common" rule:"number,min:0,max:1"`
}

var AuthRulesMessage = map[string]string{}

func (this AuthRules) Message() map[string]string {
	return AuthRulesMessage
}

func (this AuthRules) Struct() any {
	return this
}
