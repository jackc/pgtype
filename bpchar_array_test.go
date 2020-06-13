package pgtype_test

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestBPCharArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "char(8)[]", []interface{}{
		&pgtype.BPCharArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     pgtype.Present,
		},
		&pgtype.BPCharArray{
			Elements: []pgtype.BPChar{
				{String: "foo     ", Status: pgtype.Present},
				{Status: pgtype.Null},
			},
			Dimensions: []pgtype.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     pgtype.Present,
		},
		&pgtype.BPCharArray{Status: pgtype.Null},
		&pgtype.BPCharArray{
			Elements: []pgtype.BPChar{
				{String: "bar     ", Status: pgtype.Present},
				{String: "NuLL    ", Status: pgtype.Present},
				{String: `wow"quz\`, Status: pgtype.Present},
				{String: "1       ", Status: pgtype.Present},
				{String: "1       ", Status: pgtype.Present},
				{String: "null    ", Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{
				{Length: 3, LowerBound: 1},
				{Length: 2, LowerBound: 1},
			},
			Status: pgtype.Present,
		},
		&pgtype.BPCharArray{
			Elements: []pgtype.BPChar{
				{String: " bar    ", Status: pgtype.Present},
				{String: "    baz ", Status: pgtype.Present},
				{String: "    quz ", Status: pgtype.Present},
				{String: "foo     ", Status: pgtype.Present},
			},
			Dimensions: []pgtype.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: pgtype.Present,
		},
	})
}


func TestBPCharArrayMarshalJSON(t *testing.T) {
	successfulTests := []struct {
		source pgtype.BPCharArray
		result string
	}{
		{source: pgtype.BPCharArray{Status: pgtype.Null}, result: "null"},
		{source: pgtype.BPCharArray{Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 0}}, Status: pgtype.Present}, result: "[]"},
		{
			source: pgtype.BPCharArray{
				Elements: []pgtype.BPChar{
					{String: " asdf ... ", Status: pgtype.Present},
					{String: "...", Status: pgtype.Present},
				},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 2}},
				Status: pgtype.Present,
			},
			result: `[" asdf ... ","..."]`,
		},
		{
			source: pgtype.BPCharArray{
				Elements: []pgtype.BPChar{
					{String: "test123", Status: pgtype.Present},
					{Status: pgtype.Null},
					{Status: pgtype.Null},
				},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 3}},
				Status: pgtype.Present,
			},
			result: `["test123",null,null]`,
		},
		{
			source: pgtype.BPCharArray{
				Elements: []pgtype.BPChar{
					{String: "", Status: pgtype.Present},
				},
				Dimensions: []pgtype.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status: pgtype.Present,
			},
			result: `[""]`,
		},
	}
	for i, tt := range successfulTests {
		r, err := tt.source.MarshalJSON()
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if string(r) != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, string(r))
		}
	}
}