package actions

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	pattern = `{{env *"([a-zA-Z_]+[a-zA-Z0-9_]*)*" *}}`
)

// TODO: please convert it to use rune Walk and state checking
func parse(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "open")
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errors.Wrap(err, "read file")
	}
	return fillEnvs(string(b)), nil
}

func fillEnvs(content string) string {
	re := regexp.MustCompile(pattern)
	output := re.ReplaceAllStringFunc(content, func(pattern string) string {
		firstDoubleQuote := strings.Index(pattern, "\"")
		lastDoubleQuote := strings.LastIndex(pattern, "\"")
		if lastDoubleQuote-firstDoubleQuote <= 1 { // Empty env name
			return ""
		}
		if firstDoubleQuote == -1 || lastDoubleQuote == -1 || firstDoubleQuote+1 >= lastDoubleQuote {
			return pattern
		}
		variableEnvironment := pattern[firstDoubleQuote+1 : lastDoubleQuote]
		return os.Getenv(variableEnvironment)
	})
	return output
}
