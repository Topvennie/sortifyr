package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func fromString(s pgtype.Text) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func fromInt(i pgtype.Int4) int {
	if i.Valid {
		return int(i.Int32)
	}

	return 0
}

func fromBool(b pgtype.Bool) *bool {
	if b.Valid {
		return &b.Bool
	}

	return nil
}

func fromTime(t pgtype.Timestamptz) time.Time {
	if t.Valid {
		return t.Time
	}

	return time.Time{}
}
