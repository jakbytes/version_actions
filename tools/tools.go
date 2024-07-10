package tools

import (
	"fmt"
	"github.com/jakbytes/version_actions/internal/utility"
	"html"
	"os"
	"strings"
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
	if _, err := o.WriteString(fmt.Sprintf("%s=%s\n", key, escapeGitHubActionsOutput(*value))); err != nil {
		panic(err)
	}
}

func escapeGitHubActionsOutput(text string) string {
	// HTML escape the text first
	text = html.EscapeString(text)

	// Escape GitHub Actions specific characters
	replacements := []struct {
		old string
		new string
	}{
		{"%", "%25"},
		{"\n", "%0A"},
		{"\r", "%0D"},
		{"]", "%5D"},
	}
	for _, r := range replacements {
		text = strings.ReplaceAll(text, r.old, r.new)
	}
	return text
}
