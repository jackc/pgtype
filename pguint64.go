package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/jackc/pgio"
)

// pguint64 is the core type that is used to implement PostgreSQL types such as
// XID8.
type pguint64 struct {
	Uint   uint64
	Status Status
}

// Set converts from src to dst. Note that as pguint64 is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *pguint64) Set(src interface{}) error {
	switch value := src.(type) {
	case int64:
		if value < 0 {
			return fmt.Errorf("%d is less than minimum value for pguint64", value)
		}
		*dst = pguint64{Uint: uint64(value), Status: Present}
	case uint64:
		*dst = pguint64{Uint: value, Status: Present}
	default:
		return fmt.Errorf("cannot convert %v to pguint64", value)
	}

	return nil
}

func (dst pguint64) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Uint
	case Null:
		return nil
	default:
		return dst.Status
	}
}

// AssignTo assigns from src to dst. Note that as pguint64 is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *pguint64) AssignTo(dst interface{}) error {
	switch v := dst.(type) {
	case *uint64:
		if src.Status == Present {
			*v = src.Uint
		} else {
			return fmt.Errorf("cannot assign %v into %T", src, dst)
		}
	case **uint64:
		if src.Status == Present {
			n := src.Uint
			*v = &n
		} else {
			*v = nil
		}
	}

	return nil
}

func (dst *pguint64) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = pguint64{Status: Null}
		return nil
	}

	n, err := strconv.ParseUint(string(src), 10, 64)
	if err != nil {
		return err
	}

	*dst = pguint64{Uint: uint64(n), Status: Present}
	return nil
}

func (dst *pguint64) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = pguint64{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length: %v", len(src))
	}

	n := binary.BigEndian.Uint64(src)
	*dst = pguint64{Uint: n, Status: Present}
	return nil
}

func (src pguint64) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, strconv.FormatUint(src.Uint, 10)...), nil
}

func (src pguint64) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendUint64(buf, src.Uint), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *pguint64) Scan(src interface{}) error {
	if src == nil {
		*dst = pguint64{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case uint64:
		*dst = pguint64{Uint: src, Status: Present}
		return nil
	case int64:
		*dst = pguint64{Uint: uint64(src), Status: Present}
		return nil
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src pguint64) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return int64(src.Uint), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
