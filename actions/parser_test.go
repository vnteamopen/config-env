package actions

import (
	"os"
	"testing"
)

func Test_fillEnvs(t *testing.T) {
	content := `this is test file.
expected {{env "ENV_NAME"}} with content {{env "ENV_NAME2"}}
with empty env name {{env ""}}, show nothing.
with inexistent env name {{env "ENV_INEXISTENT"}}, show nothing.`
	expected := `this is test file.
expected Hello with content everyone
with empty env name , show nothing.
with inexistent env name , show nothing.`

	os.Setenv("ENV_NAME", "Hello")
	os.Setenv("ENV_NAME2", "everyone")
	if got := fillEnvs(content); got != expected {
		t.Errorf("fillEnvs() = %v, want %v", got, expected)
	}
}
