package model

import (
	"bytes"
	"fmt"
)

type ProcessErrors []ProcessError

func (p ProcessErrors) String() string {
	var buf bytes.Buffer
	for _, err := range p {
		buf.WriteString(fmt.Sprintf("%s\n", err.String()))
	}
	return buf.String()
}

func (p ProcessErrors) Error() string {
	return p.String()
}

type ProcessError struct {
	Name string
	Err  error
}

func (e *ProcessError) String() string {
	return fmt.Sprintf("%s:%v", e.Name, e.Err)
}
