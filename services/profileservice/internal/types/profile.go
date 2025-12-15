package types

import "time"

type Profile struct {
	Id        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Page      JSONB     `db:"page"`
	CreatedAt time.Time `db:"created_at"`
}
