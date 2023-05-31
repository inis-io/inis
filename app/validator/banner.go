package validator

type Banner struct{}

var BannerMessage = map[string]string{}

func (this Banner) Message() map[string]string {
	return BannerMessage
}

func (this Banner) Struct() any {
	return this
}
