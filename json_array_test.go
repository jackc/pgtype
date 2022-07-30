package pgtype_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestJSONArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "json[]", []interface{}{
		&pgtype.JSONArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     pgtype.Present,
		},
		&pgtype.JSONArray{
			Elements: []pgtype.JSON{
				{Bytes: []byte(`"foo"`), Status: pgtype.Present},
				{Status: pgtype.Null},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		},
		&pgtype.JSONArray{Status: pgtype.Null},
		&pgtype.JSONArray{
			Elements: []pgtype.JSON{
				{Bytes: []byte(`"foo"`), Status: pgtype.Present},
				{Bytes: []byte("null"), Status: pgtype.Present},
				{Bytes: []byte("42"), Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 3, LowerBound: 1}},
			Status:     pgtype.Present,
		},
	})
}

func TestJSONArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.JSONArray
	}{
		{source: []string{"{}"}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte("{}"), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 1, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
		{source: [][]byte{[]byte("{}")}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte("{}"), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 1, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
		{source: [][]byte{[]byte(`{"foo":1}`), []byte(`{"bar":2}`)}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte(`{"foo":1}`), Status: pgtype.Present}, {Bytes: []byte(`{"bar":2}`), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
		{source: []json.RawMessage{json.RawMessage(`{"foo":1}`), json.RawMessage(`{"bar":2}`)}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte(`{"foo":1}`), Status: pgtype.Present}, {Bytes: []byte(`{"bar":2}`), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
		{source: []json.RawMessage{json.RawMessage(`{"foo":12}`), json.RawMessage(`{"bar":2}`)}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte(`{"foo":12}`), Status: pgtype.Present}, {Bytes: []byte(`{"bar":2}`), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
		{source: []json.RawMessage{json.RawMessage(`{"foo":1}`), json.RawMessage(`{"bar":{"x":2}}`)}, result: pgtype.JSONArray{
			Elements:   []pgtype.JSON{{Bytes: []byte(`{"foo":1}`), Status: pgtype.Present}, {Bytes: []byte(`{"bar":{"x":2}}`), Status: pgtype.Present}},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		}},
	}

	for i, tt := range successfulTests {
		var d pgtype.JSONArray
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(d, tt.result) {
			t.Errorf("%d: expected %+v to convert to %+v, but it was %+v", i, tt.source, tt.result, d)
		}
	}
}
