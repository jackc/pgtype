package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgio"
	errors "golang.org/x/xerrors"
)

type Int4Array struct {
	Elements   []Int4
	Dimensions []ArrayDimension
	Status     Status
}

func (dst *Int4Array) Set(src interface{}) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = Int4Array{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {

	case []int16:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint16:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int32:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint32:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int64:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint64:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []int:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []uint:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			elements := make([]Int4, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = Int4Array{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     Present,
			}
		}

	case []Int4:
		if value == nil {
			*dst = Int4Array{Status: Null}
		} else if len(value) == 0 {
			*dst = Int4Array{Status: Present}
		} else {
			*dst = Int4Array{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Status:     Present,
			}
		}
	default:
		if originalSrc, ok := underlyingSliceType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Int4Array", value)
	}

	return nil
}

func (dst Int4Array) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Int4Array) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {

		case *[]int16:
			*v = make([]int16, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]uint16:
			*v = make([]uint16, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]int32:
			*v = make([]int32, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]uint32:
			*v = make([]uint32, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]int64:
			*v = make([]int64, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]uint64:
			*v = make([]uint64, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]int:
			*v = make([]int, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]uint:
			*v = make([]uint, len(src.Elements))
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

func (dst *Int4Array) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int4Array{Status: Null}
		return nil
	}

	uta, err := ParseUntypedTextArray(string(src))
	if err != nil {
		return err
	}

	var elements []Int4

	if len(uta.Elements) > 0 {
		elements = make([]Int4, len(uta.Elements))

		for i, s := range uta.Elements {
			var elem Int4
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

	*dst = Int4Array{Elements: elements, Dimensions: uta.Dimensions, Status: Present}

	return nil
}

func (dst *Int4Array) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int4Array{Status: Null}
		return nil
	}

	var arrayHeader ArrayHeader
	rp, err := arrayHeader.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	if len(arrayHeader.Dimensions) == 0 {
		*dst = Int4Array{Dimensions: arrayHeader.Dimensions, Status: Present}
		return nil
	}

	elementCount := arrayHeader.Dimensions[0].Length
	for _, d := range arrayHeader.Dimensions[1:] {
		elementCount *= d.Length
	}

	elements := make([]Int4, elementCount)

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

	*dst = Int4Array{Elements: elements, Dimensions: arrayHeader.Dimensions, Status: Present}
	return nil
}

func (src Int4Array) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
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

func (src Int4Array) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	arrayHeader := ArrayHeader{
		Dimensions: src.Dimensions,
	}

	if dt, ok := ci.DataTypeForName("int4"); ok {
		arrayHeader.ElementOID = int32(dt.OID)
	} else {
		return nil, errors.Errorf("unable to find oid for type name %v", "int4")
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
func (dst *Int4Array) Scan(src interface{}) error {
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
func (src Int4Array) Value() (driver.Value, error) {
	buf, err := src.EncodeText(nil, nil)
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}

	return string(buf), nil
}

func (src Int4Array) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		if src.Status == Null {
			return []byte("null"), nil
		} else {
			// each dimensions, marshal json and insert into result array
			if len(src.Dimensions) == 1 {
				// avoid default json.Marshal behavior for nil Element slices, this behavior is defined by src.Status
				if src.Elements == nil {
					return []byte("[]"), nil
				}
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
