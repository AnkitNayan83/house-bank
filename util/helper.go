package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func NewPgTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
