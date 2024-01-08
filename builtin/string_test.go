package builtin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"unicode/utf8"
)

func TestStringLen(t *testing.T) {
	s := "Hello, 万里"
	if got, want := len(s), 13; got != want {
		t.Errorf("len(%q) got: %d, want: %d", s, got, want)
	}

	num := 0
	for range s {
		num++
	}
	if got, want := num, 9; got != want {
		t.Errorf("count range(%q) got: %d, want: %d", s, got, want)
	}
}

// 修改 []byte(string)，不会影响原 string，原 string 保持不变。
func TestStringModify(t *testing.T) {
	s := "Hello"
	data := []byte(s)
	data[0] = 'h'
	if got, want := data[0], byte('h'); got != want {
		t.Errorf("got: %c, want: %c", got, want)
	}
	if got, want := s[0], byte('H'); got != want {
		t.Errorf("got: %c, want: %c", got, want)
	}
	if got, want := s, "Hello"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestStringRune(t *testing.T) {
	s := ":中国:"
	b := []byte(s)
	want := []bool{true, true, false, false, true, false, false, true}
	for i, _ := range b {
		assert.Equal(t, want[i], utf8.Valid(b[i:]))
	}
}

func TestStringNonPrintable(t *testing.T) {
	want := "\x01"
	got := fmt.Sprintf("%s", want)
	if len(got) != 1 {
		t.Errorf("got: %d, want: %d", len(got), len(want))
	}
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestStringFromInt(t *testing.T) {
	want := "123"
	got := strconv.Itoa(123)
	require.Equal(t, want, got)
}
