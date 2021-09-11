package pgtype_test

import (
	"context"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
	"github.com/stretchr/testify/require"
)

func TestV5ArrayV5TypeValue(t *testing.T) {
	V5ArrayV5Type := pgtype.NewV5ArrayV5Type("_text", pgtype.TextOID, func() pgtype.ValueTranscoder { return &pgtype.Text{} })

	err := V5ArrayV5Type.Set(nil)
	require.NoError(t, err)

	gotValue := V5ArrayV5Type.Get()
	require.Nil(t, gotValue)

	slice := []string{"foo", "bar"}
	err = V5ArrayV5Type.AssignTo(&slice)
	require.NoError(t, err)
	require.Nil(t, slice)

	err = V5ArrayV5Type.Set([]string{})
	require.NoError(t, err)

	gotValue = V5ArrayV5Type.Get()
	require.Len(t, gotValue, 0)

	err = V5ArrayV5Type.AssignTo(&slice)
	require.NoError(t, err)
	require.EqualValues(t, []string{}, slice)

	err = V5ArrayV5Type.Set([]string{"baz", "quz"})
	require.NoError(t, err)

	gotValue = V5ArrayV5Type.Get()
	require.Len(t, gotValue, 2)

	err = V5ArrayV5Type.AssignTo(&slice)
	require.NoError(t, err)
	require.EqualValues(t, []string{"baz", "quz"}, slice)
}

func TestV5ArrayV5TypeTranscode(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustCloseContext(t, conn)

	conn.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: pgtype.NewV5ArrayV5Type("_text", pgtype.TextOID, func() pgtype.ValueTranscoder { return &pgtype.Text{} }),
		Name:  "_text",
		OID:   pgtype.TextArrayOID,
	})

	var dstStrings []string
	err := conn.QueryRow(context.Background(), "select $1::text[]", []string{"red", "green", "blue"}).Scan(&dstStrings)
	require.NoError(t, err)

	require.EqualValues(t, []string{"red", "green", "blue"}, dstStrings)
}

func TestV5ArrayV5TypeEmptyArrayDoesNotBreakV5ArrayV5Type(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustCloseContext(t, conn)

	conn.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: pgtype.NewV5ArrayV5Type("_text", pgtype.TextOID, func() pgtype.ValueTranscoder { return &pgtype.Text{} }),
		Name:  "_text",
		OID:   pgtype.TextArrayOID,
	})

	var dstStrings []string
	err := conn.QueryRow(context.Background(), "select '{}'::text[]").Scan(&dstStrings)
	require.NoError(t, err)

	require.EqualValues(t, []string{}, dstStrings)

	err = conn.QueryRow(context.Background(), "select $1::text[]", []string{"red", "green", "blue"}).Scan(&dstStrings)
	require.NoError(t, err)

	require.EqualValues(t, []string{"red", "green", "blue"}, dstStrings)
}
