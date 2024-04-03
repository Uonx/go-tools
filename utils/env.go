package utils

import (
	"os"
	"strconv"
)

type entityEnv struct {
	name     string
	_default string
	value    string
	convert
}

func Env(name, _default string) *entityEnv {
	var e = entityEnv{}
	e._default = _default
	e.name = name
	e.value = os.Getenv(name)
	return &e
}

func (e *entityEnv) String() string {
	if len(e.value) == 0 {
		return e._default
	}
	return e.value
}

func (e *entityEnv) Int() int {
	s := e._default
	if len(e.value) != 0 {
		s = e.value
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (e *entityEnv) Bool() bool {
	if e.value == "" {
		return Convert().Bool(e._default)
	}
	return Convert().Bool(e.value)
}
