package mAuth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/S"
)

type User struct {
	*Pg.Adapter `json:"-"`
	Email       string `json:"email" db:"email"`
	Pass        string `json:"-"`
	Id          uint64

	DeletedAt   sql.NullTime `json:"deletedAt" db:"deletedAt"`
	LastLoginAt sql.NullTime `json:"lastLoginAt" db:"lastLoginAt"`
	UpdatedAt   time.Time
	CreatedAt   time.Time

	mutatedFields []string
	mutatedValues []any
}

var UserFields = []string{
	`email`,
	`pass`,
	`deletedAt`,
	`lastLoginAt`,
	`updatedAt`,
	`createdAt`,
}

func (u *User) queryParams() []any {
	return []any{
		u.Email,
		u.Pass,
		u.DeletedAt,
		u.LastLoginAt,
		u.UpdatedAt,
		u.CreatedAt,
	}
}

func (u *User) FindByEmail() bool {
	// TODO: change to all fields
	const query = `SELECT id, pass FROM users WHERE email = $1`
	u.QueryRow(context.Background(),
		query, &u.Pass)
}

func (u *User) CheckPassword(pass string) error {
	return S.CheckPassword(u.Pass, pass)
}

func (u *User) SetLastLoginAt(now int64) {
	u.LastLoginAt.Time = time.Unix(now, 0)
	u.LastLoginAt.Valid = true
	u.mutatedFields = append(u.mutatedFields, `lastLoginAt`)
	u.mutatedValues = append(u.mutatedValues, u.LastLoginAt.Time)
}

func (u *User) SetUpdatedAt(now int64) {
	u.UpdatedAt = time.Unix(now, 0)
	u.mutatedFields = append(u.mutatedFields, `updatedAt`)
	u.mutatedValues = append(u.mutatedValues, u.LastLoginAt.Time)
}

func generateIntDollar(n int) string {
	str := ``
	for z := range n {
		str += fmt.Sprintf(", $%d", z+1)
	}
	if len(str) > 0 {
		return str[1:]
	}
	return ``
}

func (u *User) conflictUpdateSql() string {
	str := ``
	for _, field := range UserFields {
		str += fmt.Sprintf(`, %s = EXCLUDED.%s`, field, field) + "\n"
	}
	if len(str) > 0 {
		return str[1:]
	}
	return ``
}

func (u *User) DoUpsert() bool {
	query := fmt.Sprintf(`-- User) DoUpsert
INSERT INTO users(%s) VALUES(%s)
ON CONFLICT DO UPDATE SET
	%s RETURNING id`,
		`"`+strings.Join(UserFields, `","`)+`"`,
		generateIntDollar(len(UserFields)),
		u.conflictUpdateSql(),
	)

	params := u.queryParams()
	row := u.Adapter.QueryRow(context.Background(), query, params...)
	err := row.Scan(&u.Id)
	if err != nil {
		log.Errorf(`failed query: %s %#v`, query, params)
		return false
	}
	return true
}

func (u *User) CensorFields() {
	u.Pass = ``
}

func NewUsersMutator(postgres *Pg.Adapter) *User {
	return &User{
		Adapter: postgres,
	}
}
