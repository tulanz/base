package nlp

import "github.com/tulanz/base/utils/simple"

type Summary interface {
	Default(title, text string, length uint64) (string, error)
}

type Default struct {
}

func (Default) Default(title, text string, length uint64) (string, error) {
	return simple.Substr(text, 0, 200), nil
}
