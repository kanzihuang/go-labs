package builtin

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_getFieldPointer(t *testing.T) {
	var (
		id   = 1024
		addr = "Beijing"
		buf  = []byte("buffer")
		time = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	)
	type testCase struct {
		name      string
		fieldName string
		id        int
		addr      string
		buf       []byte
		time      sql.NullTime
		want      any
		err       error
	}
	tests := []testCase{
		{
			fieldName: "id",
			want:      &id,
		},
		{
			fieldName: "addr",
			want:      &addr,
		},
		{
			fieldName: "buf",
			want:      &buf,
		},
		{
			fieldName: "time",
			want:      &time,
		},
		{
			fieldName: "invalid",
			err:       errors.New("field name is not found: invalid"),
		},
	}
	for _, tt := range tests {
		tt.name = fmt.Sprintf("get field %s", tt.fieldName)
		tt.id = id
		tt.buf = buf
		tt.addr = addr
		tt.time = time
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldPointer(&tt, tt.fieldName)
			require.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
