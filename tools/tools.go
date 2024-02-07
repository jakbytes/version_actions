package tools

import (
	"fmt"
	"github.com/jakbytes/version_actions/internal/utility"
	"os"
)

func String(input string) *string {
	return &input
}

type Output struct {
	*os.File
}

func OpenOutput(handler func(out Output)) {
	if os.Getenv("GITHUB_OUTPUT") == "" {
		return
	}
	err := utility.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644, func(file *os.File) error {
		handler(Output{file})
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (o *Output) Set(key string, value *string) {
	if value == nil {
		return
	}
	if _, err := o.WriteString(fmt.Sprintf("%s=%s\n", key, *value)); err != nil {
		panic(err)
	}
}
