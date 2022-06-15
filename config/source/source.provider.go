package source

import (
	"github.com/asim/go-micro/v3/config/source"
	"github.com/asim/go-micro/v3/config/source/file"
)

func NewSourceProvider() source.Source {
	return file.NewSource(file.WithPath("./config.json"))
}

