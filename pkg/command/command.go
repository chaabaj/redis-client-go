package command

type Command interface {
	Encode() string
}
