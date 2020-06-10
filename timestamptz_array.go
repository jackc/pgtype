package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/jackc/pgio"
	errors "golang.org/x/xerrors"
)

type TimestamptzArray struct {
	Elements   []Timestamptz
	Dimensions []ArrayDimension
	Status     Status
}

func (dst *TimestamptzArray) Set(src interface{}) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = TimestamptzArray{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {

	case []time.Time:
		if value == nil {
			*dst = TimestamptzArray{Status: Null}
		} else if len(value) == 0 {
			*dst = TimestamptzArray{Status: Present}
		} else {
			elements := make([]Timestamptz, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = TimestamptzArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []Timestamptz:
		if value == nil {
			*dst = TimestamptzArray{Status: Null}
		} else if len(value) == 0 {
			*dst = TimestamptzArray{Status: Present}
		} else {
			*dst = TimestamptzArray{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Status:     Present,
			}
		}
	default:
		if originalSrc, ok := underlyingSliceType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to TimestamptzArray", value)
	}

	return nil
}

func (dst TimestamptzArray) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *TimestamptzArray) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {

		case *[]time.Time:
			*v = make([]time.Time, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return errors.Errorf("unable to assign to %T", dst)
		}
	case Null:
		return NullAssignTo(dst)
	}

	return errors.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *TimestamptzArray) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TimestamptzArray{Status: Null}
		return nil
	}

	uta, err := ParseUntypedTextArray(string(src))
	if err != nil {
		return err
	}

	var elements []Timestamptz

	if len(uta.Elements) > 0 {
		elements = make([]Timestamptz, len(uta.Elements))

		for i, s := range uta.Elements {
			var elem Timestamptz
			var elemSrc []byte
			if s != "NULL" {
				elemSrc = []byte(s)
			}
			err = elem.DecodeText(ci, elemSrc)
			if err != nil {
				return err
			}

			elements[i] = elem
		}
	}

	*dst = TimestamptzArray{Elements: elements, Dimensions: uta.Dimensions, Status: Present}

	return nil
}

func (dst *TimestamptzArray) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = TimestamptzArray{Status: Null}
		return nil
	}

	var arrayHeader ArrayHeader
	rp, err := arrayHeader.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	if len(arrayHeader.Dimensions) == 0 {
		*dst = TimestamptzArray{Dimensions: arrayHeader.Dimensions, Status: Present}
		return nil
	}

	elementCount := arrayHeader.Dimensions[0].Length
	for _, d := range arrayHeader.Dimensions[1:] {
		elementCount *= d.Length
	}

	elements := make([]Timestamptz, elementCount)

	for i := range elements {
		elemLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4
		var elemSrc []byte
		if elemLen >= 0 {
			elemSrc = src[rp : rp+elemLen]
			rp += elemLen
		}
		err = elements[i].DecodeBinary(ci, elemSrc)
		if err != nil {
			return err
		}
	}

	*dst = TimestamptzArray{Elements: elements, Dimensions: arrayHeader.Dimensions, Status: Present}
	return nil
}

func (src TimestamptzArray) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	if len(src.Dimensions) == 0 {
		return append(buf, '{', '}'), nil
	}

	buf = EncodeTextArrayDimensions(buf, src.Dimensions)

	dimElemCounts := getDimElemCounts(src.Dimensions)

	inElemBuf := make([]byte, 0, 32)
	for i, elem := range src.Elements {
		if i > 0 {
			buf = append(buf, ',')
		}

		for _, dec := range dimElemCounts {
			if i%dec == 0 {
				buf = append(buf, '{')
			}
		}

		elemBuf, err := elem.EncodeText(ci, inElemBuf)
		if err != nil {
			return nil, err
		}
		if elemBuf == nil {
			buf = append(buf, `NULL`...)
		} else {
			buf = append(buf, QuoteArrayElementIfNeeded(string(elemBuf))...)
		}

		for _, dec := range dimElemCounts {
			if (i+1)%dec == 0 {
				buf = append(buf, '}')
			}
		}
	}

	return buf, nil
}

func (src TimestamptzArray) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	arrayHeader := ArrayHeader{
		Dimensions: src.Dimensions,
	}

	if dt, ok := ci.DataTypeForName("timestamptz"); ok {
		arrayHeader.ElementOID = int32(dt.OID)
	} else {
		return nil, errors.Errorf("unable to find oid for type name %v", "timestamptz")
	}

	for i := range src.Elements {
		if src.Elements[i].Status == Null {
			arrayHeader.ContainsNull = true
			break
		}
	}

	buf = arrayHeader.EncodeBinary(ci, buf)

	for i := range src.Elements {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		elemBuf, err := src.Elements[i].EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if elemBuf != nil {
			buf = elemBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *TimestamptzArray) Scan(src interface{}) error {
	if src == nil {
		return dst.DecodeText(nil, nil)
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src TimestamptzArray) Value() (driver.Value, error) {
	buf, err := src.EncodeText(nil, nil)
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}

	return string(buf), nil
}

func (src TimestamptzArray) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		if src.Status == Null {
			return []byte("null"), nil
		} else {
			// each dimensions, marshal json and insert into result array
			if len(src.Dimensions) == 1 {
				// not nested
				bytes, err := json.Marshal(src.Elements)
				if err != nil {
					return nil, err
				}
				return bytes, nil
			}

			dimElemCounts := getDimElemCounts(src.Dimensions)
			var bytes []byte
			for i, elem := range src.Elements {
				if i > 0 {
					bytes = append(bytes, ',')
				}
				for _, dec := range dimElemCounts {
					if i%dec == 0 {
						bytes = append(bytes, '[')
					}
				}

				eb, err := json.Marshal(elem)
				if err != nil {
					return nil, err
				}
				bytes = append(bytes, eb...)

				for _, dec := range dimElemCounts {
					if (i+1)%dec == 0 {
						bytes = append(bytes, ']')
					}
				}
			}
			return bytes, nil
		}
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	return nil, errBadStatus
}
