package validator

type IpBlack struct{}

var IpBlackMessage = map[string]string{}

func (this IpBlack) Message() map[string]string {
	return IpBlackMessage
}

func (this IpBlack) Struct() any {
	return this
}
