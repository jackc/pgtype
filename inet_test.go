package pgtype_test

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInetTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "inet", []interface{}{
		&pgtype.Inet{IPNet: mustParseInet(t, "0.0.0.0/32"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "127.0.0.1/8"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "12.34.56.65/32"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "192.168.1.16/24"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "255.0.0.0/8"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "255.255.255.255/32"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "10.0.0.1"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "2607:f8b0:4009:80b::200e"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "::1/64"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "::/0"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "::1/128"), Status: pgtype.Present},
		&pgtype.Inet{IPNet: mustParseInet(t, "2607:f8b0:4009:80b::200e/64"), Status: pgtype.Present},
		&pgtype.Inet{Status: pgtype.Null},
	})
}

func TestCidrTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "cidr", []interface{}{
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

type textMarshaler struct {
	Text string
}

func (t textMarshaler) MarshalText() (text []byte, err error) {
	return []byte(t.Text), err
}

func TestInetSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.Inet
	}{
		{source: mustParseCIDR(t, "127.0.0.1/32"), result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: mustParseCIDR(t, "127.0.0.1/32").IP, result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: "127.0.0.1/32", result: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}},
		{source: "1.2.3.4/24", result: pgtype.Inet{IPNet: &net.IPNet{IP: net.ParseIP("1.2.3.4").To4(), Mask: net.CIDRMask(24, 32)}, Status: pgtype.Present}},
		{source: "10.0.0.1", result: pgtype.Inet{IPNet: mustParseInet(t, "10.0.0.1"), Status: pgtype.Present}},
		{source: "2607:f8b0:4009:80b::200e", result: pgtype.Inet{IPNet: mustParseInet(t, "2607:f8b0:4009:80b::200e"), Status: pgtype.Present}},
		{source: net.ParseIP(""), result: pgtype.Inet{Status: pgtype.Null}},
		{source: "0.0.0.0/8", result: pgtype.Inet{IPNet: mustParseInet(t, "0.0.0.0/8"), Status: pgtype.Present}},
		{source: "::ffff:0.0.0.0/104", result: pgtype.Inet{IPNet: &net.IPNet{IP: net.ParseIP("::ffff:0.0.0.0"), Mask: net.CIDRMask(104, 128)}, Status: pgtype.Present}},
		{source: textMarshaler{"127.0.0.1"}, result: pgtype.Inet{IPNet: mustParseInet(t, "127.0.0.1"), Status: pgtype.Present}},
		{source: func(s string) fmt.Stringer {
			var b strings.Builder
			b.WriteString(s)
			return &b
		}("127.0.0.1"), result: pgtype.Inet{IPNet: mustParseInet(t, "127.0.0.1"), Status: pgtype.Present}},
	}

	for i, tt := range successfulTests {
		var r pgtype.Inet
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		assert.Equalf(t, tt.result.Status, r.Status, "%d: Status", i)
		if tt.result.Status == pgtype.Present {
			assert.Equalf(t, tt.result.IPNet.Mask, r.IPNet.Mask, "%d: Mask", i)
			assert.Truef(t, tt.result.IPNet.IP.Equal(r.IPNet.IP), "%d: IP", i)
			assert.Equalf(t, len(tt.result.IPNet.IP), len(r.IPNet.IP), "%d: IP length", i)
		}
	}
}

type textUnmarshaler struct {
	Text string
}

func (u *textUnmarshaler) UnmarshalText(text []byte) error {
	u.Text = string(text)
	return nil
}

func TestInetAssignTo(t *testing.T) {
	var ipnet net.IPNet
	var pipnet *net.IPNet
	var ip net.IP
	var pip *net.IP
	var um textUnmarshaler
	var pum *textUnmarshaler

	simpleTests := []struct {
		src      pgtype.Inet
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &ipnet, expected: *mustParseCIDR(t, "127.0.0.1/32")},
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &ip, expected: mustParseCIDR(t, "127.0.0.1/32").IP},
		{src: pgtype.Inet{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: pgtype.Present}, dst: &um, expected: textUnmarshaler{"127.0.0.1/32"}},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &pipnet, expected: ((*net.IPNet)(nil))},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &pip, expected: ((*net.IP)(nil))},
		{src: pgtype.Inet{Status: pgtype.Null}, dst: &pum, expected: ((*textUnmarshaler)(nil))},
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
