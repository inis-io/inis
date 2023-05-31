package validator

type Placard struct{}

var PlacardMessage = map[string]string{}

func (this Placard) Message() map[string]string {
	return PlacardMessage
}

func (this Placard) Struct() any {
	return this
}
