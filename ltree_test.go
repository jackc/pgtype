package pgtype_test

import (
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestLtreeTranscode(t *testing.T) {
	values := []interface{}{
		&pgtype.Ltree{String: "", Status: pgtype.Present},
		&pgtype.Ltree{String: "All.foo.one", Status: pgtype.Present},
		&pgtype.Ltree{Status: pgtype.Null},
	}

	testutil.TestSuccessfulTranscodeEqFunc(
		t, "ltree", values, func(ai, bi interface{}) bool {
			a := ai.(pgtype.Ltree)
			b := bi.(pgtype.Ltree)

			if a.String != b.String || a.Status != b.Status {
				return false
			}
			return true
		},
	)

}

func TestLtreeSet(t *testing.T) {
	successfulTests := []struct {
		src    interface{}
		result pgtype.Ltree
	}{
		{src: "All.foo.bar", result: pgtype.Ltree{String: "All.foo.bar", Status: pgtype.Present}},
		{src: (*string)(nil), result: pgtype.Ltree{Status: pgtype.Null}},
	}
	for i, tt := range successfulTests {
		var dst pgtype.Ltree
		err := dst.Set(tt.src)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !reflect.DeepEqual(dst, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.src, tt.result, dst)
		}
	}
}
