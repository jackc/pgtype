package pgtype_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestInetTranscode(t *testing.T) {
	for _, pgTypeName := range []string{"inet", "cidr"} {
		testutil.TestSuccessfulTranscode(t, pgTypeName, []interface{}{
			&pgtype.Inet{IPNet: mustParseCIDR(t, "0.0.0.0/32"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "192.168.1.0/24"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "255.0.0.0/8"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "255.255.255.255/32"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "::/128"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "::/0"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "::1/128"), Status: pgtype.Present},
			&pgtype.Inet{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: pgtype.Present},
			&pgtype.Inet{Status: pgtype.Null},
		})
	}
}

func TestInetSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.Inet
	}{
		{source: mustParseCIDR(t, "127.0.0.1/32"), result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: mustParseCIDR(t, "127.0.0.1/32").IP, result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: "127.0.0.1/32", result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: net.ParseIP(""), result: pgtype.Inet{Status: pgtype.Null}},
	}

	for i, tt := range successfulTests {
		var r pgtype.Inet
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInetAssignTo(t *testing.T) {
	var ipnet net.IPNet
	var pipnet *net.IPNet
	var ip net.IP
	var pip *net.IP

	simpleTests := []struct {
		src      pgtype.Inet
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &ipnet, expected: *mustParseCIDR(t, "127.0.0.1/32")},
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &ip, expected: mustParseCIDR(t, "127.0.0.1/32").IP},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &pipnet, expected: ((*net.IPNet)(nil))},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &pip, expected: ((*net.IP)(nil))},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %#v, but result was %#v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      pgtype.Inet
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &pipnet, expected: *mustParseCIDR(t, "127.0.0.1/32")},
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &pip, expected: mustParseCIDR(t, "127.0.0.1/32").IP},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src pgtype.Inet
		dst interface{}
	}{
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "192.168.0.0/16"), Status: pgtype.Present}, dst: &ip},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &ipnet},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}

func TestInetMarshalJSON(t *testing.T) {
	successfulTests := []struct {
		json   string
		source pgtype.Inet
	}{
		{source: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, json: `"127.0.0.1/32"`},
		{source: pgtype.Inet{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: pgtype.Present}, json: `"2607:f8b0:4009:80b::200e/128"`},
		{source: pgtype.Inet{Status: pgtype.Null}, json: `""`},
		{source: pgtype.Inet{}, json: `""`},
	}

	for i, tt := range successfulTests {
		got, err := tt.source.MarshalJSON()
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !reflect.DeepEqual(got, []byte(tt.json)) {
			t.Errorf("%d: expected JSON `%s`, but it was %s", i, tt.json, string(got))
		}
	}
}

func TestInetUnmarshalJSON(t *testing.T) {
	successfulTests := []struct {
		json     string
		expected pgtype.Inet
	}{
		{expected: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, json: `"127.0.0.1/32"`},
		{expected: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, json: `"127.0.0.1"`},
		{expected: pgtype.Inet{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: pgtype.Present}, json: `"2607:f8b0:4009:80b::200e/128"`},
		{expected: pgtype.Inet{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: pgtype.Present}, json: `"2607:f8b0:4009:80b::200e"`},
		{expected: pgtype.Inet{Status: pgtype.Null}, json: `""`}, // empty is OK, equivalent to our null struct
	}
	badJSON := []string{
		`"127.0.0.1/"`,         // no network
		`"444.555.666.777/32"`, // bad addr
		`"nonsense"`,           // bad everything
	}

	for i, tt := range successfulTests {
		got := pgtype.Inet{}
		if err := got.UnmarshalJSON([]byte(tt.json)); err != nil {
			t.Errorf("%d: %v", i, err)
		}
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%d: expected %v from JSON `%s`, but it was %v", i, tt.expected, tt.json, got)
		}
	}

	for i, example := range badJSON {
		got := pgtype.Inet{}
		if err := got.UnmarshalJSON([]byte(example)); err == nil {
			t.Errorf("%d: Expected error for %s, but got none", i, example)
		}
	}
}
