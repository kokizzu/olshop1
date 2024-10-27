package domain

import (
	"olshop1/model/mAuth"

	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestLogin.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestLogin.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestLogin.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestLogin.go
//go:generate farify doublequote --file GuestLogin.go

type (
	GuestLoginIn struct {
		RequestCommon
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	GuestLoginOut struct {
		ResponseCommon
		User *mAuth.User `json:"user" form:"user" query:"user" long:"user" msg:"user"`

		Segments M.SB `json:"segments" form:"segments" query:"segments" long:"segments" msg:"segments"`
	}
)

const (
	GuestLoginAction = `guest/login`

	ErrGuestLoginEmailInvalid             = `email must be valid`
	ErrGuestLoginUserDeactivated          = `user deactivated`
	ErrGuestLoginEmailOrPasswordIncorrect = `incorrect email or password`
	ErrGuestLoginPasswordOrEmailIncorrect = `incorrect password or email`
	ErrGuestLoginFailedStoringSession     = `failed storing session for login`

	WarnFailedSetLastLoginAt = `failed setting lastLoginAt`
)

func (d *Domain) GuestLogin(in *GuestLoginIn) (out GuestLoginOut) {
	in.Email = S.Trim(S.ValidateEmail(in.Email))
	if in.Email == `` {
		out.SetError(400, ErrGuestLoginEmailInvalid)
		return
	}
	user := mAuth.NewUsersMutator(d.Postgres)
	user.Email = in.Email
	if !user.FindByEmail() {
		out.SetError(400, ErrGuestLoginEmailOrPasswordIncorrect)
		return
	}
	out.actor = user.Id
	out.refId = user.Id

	if user.DeletedAt.Valid {
		out.SetError(400, ErrGuestLoginUserDeactivated)
		return
	}

	if err := user.CheckPassword(in.Password); err != nil {
		out.SetError(400, ErrGuestLoginPasswordOrEmailIncorrect)
		return
	}
	user.SetLastLoginAt(in.UnixNow())
	user.SetUpdatedAt(in.UnixNow())
	if !user.DoUpsert() {
		out.AddTrace(WarnFailedSetLastLoginAt)
		return
	}
	user.CensorFields()
	out.User = user
	session, sess := d.CreateSession(user.Id, user.Email, in.UserAgent, in.IpAddress)

	// TODO: set list of roles in the session
	if !session.DoInsert() {
		out.SetError(500, ErrGuestLoginFailedStoringSession)
		return
	}
	out.SessionToken = session.SessionToken
	out.Segments = d.segmentsFromSession(sess)
	return
}
