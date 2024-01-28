package builtin

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func chProjectDir(projectName string, directories ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	directories = append([]string{"/", projectName}, directories...)
	suffix := filepath.Join(directories...)
	for {
		if strings.HasSuffix(wd, suffix) {
			if err := os.Chdir(wd); err != nil {
				return "", err
			}
			return wd, nil
		}
		upper := filepath.Join(wd, "..")
		if upper == wd {
			return "", errors.New(fmt.Sprintf("not found directory: %s", suffix[1:]))
		}
		wd = upper
	}
}

func TestGetwd(t *testing.T) {
	wd, err := chProjectDir("go-labs", "builtin")
	require.NoError(t, err)
	want := filepath.Join(wd, "..", "builtin")
	require.Equal(t, want, wd)
}
