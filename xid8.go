package pgtype

import (
	"database/sql/driver"
)

// XID8 is PostgreSQL's 64 bit Transaction ID type.
//
// Another identifier type used by the system is xid, or transaction (abbreviated xact) identifier. This is the data
// type of the system columns xmin and xmax. Transaction identifiers are 32-bit quantities.
// In some contexts, a 64-bit variant xid8 is used. Unlike xid values, xid8 values increase strictly monotonically and
// cannot be reused in the lifetime of a database cluster. See Section 66.1 for more details.
type XID8 pguint64

// Set converts from src to dst. Note that as XID8 is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *XID8) Set(src interface{}) error {
	return (*pguint64)(dst).Set(src)
}

func (dst XID8) Get() interface{} {
	return (pguint64)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as XID8 is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *XID8) AssignTo(dst interface{}) error {
	return (*pguint64)(src).AssignTo(dst)
}

func (dst *XID8) DecodeText(ci *ConnInfo, src []byte) error {
	return (*pguint64)(dst).DecodeText(ci, src)
}

func (dst *XID8) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*pguint64)(dst).DecodeBinary(ci, src)
}

func (src XID8) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (pguint64)(src).EncodeText(ci, buf)
}

func (src XID8) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (pguint64)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *XID8) Scan(src interface{}) error {
	return (*pguint64)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src XID8) Value() (driver.Value, error) {
	return (pguint64)(src).Value()
}
