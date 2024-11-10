package mAuth

import (
	"context"
	"fmt"

	"olshop1/model"

	"github.com/gofiber/fiber/v2/log"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Pg"
)

type Session struct {
	*Pg.Adapter  `json:"-"`
	UserId       uint64
	Device       string
	LoginIPs     string
	SessionToken string // unique
	ExpiredAt    int64

	model.Mutator `json:"-"`
}

var SessionColumns = []string{`userId`, `device`, `loginIPs`, `sessionToken`, `expiredAt`}

func (s *Session) FindBySessionToken() bool {
	query := fmt.Sprintf(`-- Session) FindBySessionToken
SELECT %s
FROM sessions
WHERE sessionToken = $1`, A.StrJoin(SessionColumns, `,`))
	row := s.QueryRow(context.Background(), query, s.SessionToken)
	err := row.Scan(&s.UserId, &s.Device, &s.LoginIPs, &s.SessionToken, &s.ExpiredAt)
	if err != nil {
		log.Errorf("query error: %v %s %#v", err, query, s.SessionToken)
		return false
	}
	return true
}

func (s *Session) SetExpiredAt(future int64) {
	s.ExpiredAt = future

}

func (s *Session) queryParams() []any {
	return []any{
		s.UserId, s.Device, s.LoginIPs, s.SessionToken, s.ExpiredAt,
	}
}

func (s *Session) DoUpdateBySessionToken() bool {
	updateFields, updateValues := s.Mutator.ToUpdateQueryString()
	query := fmt.Sprintf(`-- Session) DoUpdateBySessionToken
UPDATE sessions
SET %s
WHERE sessionToken = $%d`,
		updateFields,
		len(updateValues))

	vals := append(updateValues, s.SessionToken)
	ct, err := s.Exec(context.Background(), query, vals...)
	if err != nil {
		log.Errorf("query error: %v %s %#v", err, query, vals)
		return false
	}

	return ct.RowsAffected() > 0
}

func (s *Session) DoInsert() bool {
	query := fmt.Sprintf(`-- Session) DoInsert
INSERT INTO sessions(%s)
VALUES(%s)
`, A.StrJoin(SessionColumns, `,`),
		model.GenerateDollar(len(SessionColumns)))

	vals := s.queryParams()
	ins, err := s.Adapter.Exec(context.Background(), query, vals...)
	if err == nil {
		return ins.RowsAffected() > 0
	}
	log.Errorf("query error: %v %s %#v", err, query, vals)
	return false
}

func NewSessionsMutator(postgres *Pg.Adapter) *Session {
	return &Session{
		Adapter: postgres,
	}
}

func (s *Session) Migrate() bool {
	const query = `-- Session) Migrate
CREATE TABLE IF NOT EXISTS sessions (
	userId BIGINT NOT NULL,
	device VARCHAR(255) NOT NULL,
	loginIPs VARCHAR(255) NOT NULL,
	sessionToken VARCHAR(255) NOT NULL,
	expiredAt BIGINT NOT NULL,
	PRIMARY KEY (sessionToken),
	FOREIGN KEY (userId) REFERENCES users(id) 
)`

	_, err := s.Adapter.Exec(context.Background(), query)
	if err != nil {
		log.Errorf("query error: %v %s", err, query)
		return false
	}
	return true
}
