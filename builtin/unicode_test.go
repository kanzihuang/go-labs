package builtin

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUnicode_Equal(t *testing.T) {
	text := "新\r\n年\r快\n乐"
	rows := strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case '\r', '\n':
			return true
		default:
			return false
		}
	})
	require.Equal(t, []string{"新", "年", "快", "乐"}, rows)
	t.Log("rows: ", rows)
}
