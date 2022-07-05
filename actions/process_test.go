package actions

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	testKey := "TEST_USER"
	testValue := "user"

	os.Setenv(testKey, testValue)

	testCases := []struct {
		name        string
		pattern     []string
		template    string
		envName     string
		checkResult func(expected, received string)
		checkError  func(err error)
	}{
		{
			name:     "Matched pattern - Simple template",
			template: fmt.Sprintf(`{{env "%s"}}`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				if envValue != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", envValue, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Matched pattern - Template contains 2 begin parts",
			template: fmt.Sprintf(`{{env {{env "%s"}}`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf("{{env %s", envValue)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Matched pattern - Template contains 2 end parts",
			template: fmt.Sprintf(`{{env "%s"}}"}}`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf(`%s"}}`, envValue)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Matched pattern - Template contains begin part at then end",
			template: fmt.Sprintf(`{{env "%s"}}{{env`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf(`%s{{env`, envValue)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Matched pattern - Template contains another parts except template",
			template: fmt.Sprintf(`abc {{env "%s"}} def`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf(`abc %s def`, envValue)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Not match pattern - Simple template",
			template: fmt.Sprintf(`{{environment "%s"}}`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf(`{{environment "%s"}}`, testKey)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:     "Not match pattern - Invalid env name in template",
			template: fmt.Sprintf(`{{env "'%s'"}}`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				expected := fmt.Sprintf(`{{env "'%s'"}}`, testKey)
				if expected != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", expected, received)
				}
			},
			checkError: func(err error) {},
		},
		{
			name:        "Not match pattern - Notfound env name in template",
			template:    fmt.Sprintf(`{{env "%s"}}`, "_INVALID"),
			envName:     testKey,
			checkResult: func(envValue, received string) {},
			checkError: func(err error) {
				if !strings.Contains(err.Error(), "open") {
					t.Errorf("Wrong error: \n Expected: %+v\nReceived: %+v", errors.Wrap(err, "open"), err)
				}
			},
		},
		{
			name:     "Matched custom pattern - Simple template",
			pattern:  []string{"%", "%"},
			template: fmt.Sprintf(`%%env "%s"%%`, testKey),
			envName:  testKey,
			checkResult: func(envValue, received string) {
				if envValue != received {
					t.Errorf("Wrong parse: \nExpected: %s\nReceived: %s", envValue, received)
				}
			},
			checkError: func(err error) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.template)
			in := bufio.NewReader(reader)
			parser := NewSeqParser(tc.pattern)

			writer := new(strings.Builder)
			out := bufio.NewWriter(writer)

			err := parseInputToOutput(parser, in, []*bufio.Writer{out})
			if err != nil {
				tc.checkError(err)
				return
			}

			envValue := os.Getenv(testKey)
			if err != nil {
				t.Errorf("Failed to read sample content: %v+", err.Error())
				return
			}
			tc.checkResult(envValue, writer.String())
		})
	}
}
