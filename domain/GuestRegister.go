package domain

import (
	"olshop1/model/mAuth"

	"github.com/kokizzu/gotro/S"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file GuestRegister.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type GuestRegister.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type GuestRegister.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type GuestRegister.go
//go:generate farify doublequote --file GuestRegister.go

type (
	GuestRegisterIn struct {
		RequestCommon
		Email    string `json:"email" form:"email" query:"email" long:"email" msg:"email"`
		Password string `json:"password" form:"password" query:"password" long:"password" msg:"password"`
	}
	GuestRegisterOut struct {
		ResponseCommon
		User mAuth.User `json:"user" form:"user" query:"user" long:"user" msg:"user"`

		verifyEmailUrl string
	}
)

const (
	GuestRegisterAction = `guest/register`

	ErrGuestRegisterEmailInvalid       = `email must be valid`
	ErrGuestRegisterPasswordTooShort   = `password must be at least 12 characters`
	ErrGuestRegisterEmailUsed          = `email already used`
	ErrGuestRegisterUserCreationFailed = `user creation failed`

	minPassLength = 12
)

func (d *Domain) GuestRegister(in *GuestRegisterIn) (out GuestRegisterOut) {
	in.Email = S.Trim(S.ValidateEmail(in.Email))
	if in.Email == `` {
		out.SetError(400, ErrGuestRegisterEmailInvalid)
		return
	}
	if len(in.Password) < minPassLength {
		out.SetErrorf(400, ErrGuestRegisterPasswordTooShort)
		return
	}
	exists := mAuth.NewUsersMutator(d.Postgres)
	exists.Email = in.Email
	if exists.FindByEmail() {
		out.SetError(400, ErrGuestRegisterEmailUsed)
		return
	}
	user := mAuth.NewUsersMutator(d.Postgres)
	user.Email = in.Email
	user.SetEncryptedPassword(in.Password)

	user.SetUpdatedAt(in.UnixNow())
	user.SetCreatedAt(in.UnixNow())
	if !user.DoInsert() {
		out.SetError(500, ErrGuestRegisterUserCreationFailed)
		return
	}
	out.actor = user.Id
	out.refId = user.Id

	user.CensorFields()
	out.User = *user

	return
}
