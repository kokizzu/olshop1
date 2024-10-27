package mAuth

import (
	"context"
	"fmt"

	"olshop1/model"

	"github.com/gofiber/fiber/v2/log"
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

func (s *Session) FindBySessionToken() bool {
	// TODO: continue this
}

func (s *Session) SetExpiredAt(future int64) {
	s.ExpiredAt = future

}

func (s *Session) DoUpdateBySessionToken() bool {
	updateFields, updateValues := s.Mutator.ToQueryString()
	query := fmt.Sprintf(`-- Session) DoUpdateBySessionToken
UPDATE sessions
SET %s
WHERE sessionToken = $%d`,
		updateFields,
		len(updateValues))

	vals := append(updateValues, s.SessionToken)
	ct, err := s.Exec(context.Background(), query, vals...)
	if err != nil {
		log.Errorf("query error: %s %#v", query, vals)
		return false
	}

	return ct.RowsAffected() > 0
}

func (s *Session) DoInsert() bool {
	query :=fmt.Sprintf(`-- Session) DoInsert
INSERT INTO sessions(%s)
VALUES(%s)`


}

func NewSessionsMutator(postgres *Pg.Adapter) *Session {
	return &Session{
		Adapter: postgres,
	}
}
