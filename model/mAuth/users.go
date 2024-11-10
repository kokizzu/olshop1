package mAuth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"olshop1/model"

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

	model.Mutator
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
	query := fmt.Sprintf(`-- User) FindByEmail
SELECT %s 
FROM users 
WHERE email = $1
`, strings.Join(UserFields, `,`))
	row := u.QueryRow(context.Background(),
		query, u.Email)
	err := row.Scan(&u.Id, &u.Pass)
	if err != nil {
		log.Errorf("query error: %v %s %#v", err, query, u.Email)
		return false
	}
	return true
}

func (u *User) CheckPassword(pass string) error {
	return S.CheckPassword(u.Pass, pass)
}

func (u *User) SetLastLoginAt(now int64) {
	u.LastLoginAt.Time = time.Unix(now, 0)
	u.LastLoginAt.Valid = true
	u.Mutator.Add(`lastLoginAt`, u.LastLoginAt.Time)
}

func (u *User) SetUpdatedAt(now int64) {
	u.UpdatedAt = time.Unix(now, 0)
	u.Mutator.Add(`updatedAt`, u.UpdatedAt)
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
		model.GenerateDollar(len(UserFields)),
		u.conflictUpdateSql(),
	)

	params := u.queryParams()
	row := u.Adapter.QueryRow(context.Background(), query, params...)
	err := row.Scan(&u.Id)
	if err != nil {
		log.Errorf(`failed query: %v %s %#v`, err, query, params)
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

func (u *User) Migrate() bool {
	const query = `-- User) Migrate
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL NOT NULL,
	email VARCHAR(255) NOT NULL,
	pass VARCHAR(255) NOT NULL,
	deletedAt BIGINT,
	lastLoginAt BIGINT,
	updatedAt BIGINT NOT NULL,
	createdAt BIGINT NOT NULL,
	PRIMARY KEY (id)
)`

	_, err := u.Adapter.Exec(context.Background(), query)
	if err != nil {
		log.Errorf("query error: %v %s", err, query)
		return false
	}
	return true
}

func (u *User) FindById() bool {
	query := fmt.Sprintf(`-- User) FindById
SELECT %s
FROM users
WHERE id = $1`, strings.Join(UserFields, `,`))
	row := u.QueryRow(context.Background(), query, u.Id)
	err := row.Scan(&u.Id, &u.Email, &u.Pass, &u.DeletedAt, &u.LastLoginAt, &u.UpdatedAt, &u.CreatedAt)
	if err != nil {
		log.Errorf("query error: %v %s %#v", err, query, u.Id)
		return false
	}
	return true
}

func (u *User) SetEncryptedPassword(password string) {
	u.SetPassword(S.EncryptPassword(password))
}

func (u *User) SetPassword(val string) bool {
	if val != u.Pass {
		u.Mutator.Add(`pass`, val)
		u.Pass = val
		return true
	}
	return false
}

func (u *User) SetCreatedAt(now int64) {
	u.CreatedAt = time.Unix(now, 0)
	u.Mutator.Add(`createdAt`, u.CreatedAt)
}

func (u *User) DoInsert() bool {
	query := fmt.Sprintf(`-- User) DoInsert
INSERT INTO users(%s) 
VALUES(%s) 
RETURNING id`,
		strings.Join(UserFields, `,`),
		model.GenerateDollar(len(UserFields)),
	)
	params := u.queryParams()
	row := u.Adapter.QueryRow(context.Background(), query, params...)
	err := row.Scan(&u.Id)
	if err != nil {
		log.Errorf(`failed query: %v %s %#v`, err, query, params)
		return false
	}
	return true
}
