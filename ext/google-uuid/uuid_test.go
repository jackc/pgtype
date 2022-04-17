package uuid_test

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ggl "github.com/jackc/pgtype/ext/google-uuid"
)

func TestUUIDTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "uuid", []interface{}{
		&ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		&ggl.UUID{Status: pgtype.Null},
	})
}

func TestUUIDSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source interface{}
		result ggl.UUID
	}{
		{
			name:   "source uuid type",
			source: &ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
			result: ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		},
		{
			name:   "source fixed bytes",
			source: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		},
		{
			name:   "source variable bytes",
			source: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		},
		{
			name:   "source string",
			source: "00010203-0405-0607-0809-0a0b0c0d0e0f",
			result: ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r ggl.UUID
			err := r.Set(tt.source)
			require.NoError(t, err)
			assert.Equal(t, tt.result, r)
		})
	}
}

func TestUUIDAssignTo(t *testing.T) {
	t.Run("assign to fixed bytes", func(t *testing.T) {
		src := ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present}
		var dst [16]byte
		expected := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		err := src.AssignTo(&dst)
		require.NoError(t, err)
		assert.Equal(t, expected, dst)
	})

	t.Run("assign to variable bytes", func(t *testing.T) {
		src := ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present}
		var dst []byte
		expected := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		err := src.AssignTo(&dst)
		require.NoError(t, err)
		assert.Equal(t, expected, dst)
	})

	t.Run("assign to string", func(t *testing.T) {
		src := ggl.UUID{UUID: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present}
		var dst string
		expected := "00010203-0405-0607-0809-0a0b0c0d0e0f"

		err := src.AssignTo(&dst)
		require.NoError(t, err)
		assert.Equal(t, expected, dst)
	})
}
