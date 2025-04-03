package config

import (
	"errors"
)

var (
	errStat = errors.New("no file")
	errLoad = errors.New("cannot load file")
	errRead = errors.New("cannot read file")
)
