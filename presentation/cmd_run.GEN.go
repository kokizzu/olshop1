package presentation

import (
	"os"

	"olshop1/domain"
)


// Code generated by 1_codegen_test.go DO NOT EDIT.


func cmdRun(b *domain.Domain, action string, payload []byte) {
	switch action {
	case domain.GuestDebugAction:
		in := domain.GuestDebugIn{}
		if !in.RequestCommon.FromCli(action, payload, &in) {
			return
		}
		out := b.GuestDebug(&in)
		in.RequestCommon.ToCli(os.Stdout, out, out.ResponseCommon)

	case domain.GuestLoginAction:
		in := domain.GuestLoginIn{}
		if !in.RequestCommon.FromCli(action, payload, &in) {
			return
		}
		out := b.GuestLogin(&in)
		in.RequestCommon.ToCli(os.Stdout, out, out.ResponseCommon)

	case domain.GuestRegisterAction:
		in := domain.GuestRegisterIn{}
		if !in.RequestCommon.FromCli(action, payload, &in) {
			return
		}
		out := b.GuestRegister(&in)
		in.RequestCommon.ToCli(os.Stdout, out, out.ResponseCommon)

	case domain.UserLogoutAction:
		in := domain.UserLogoutIn{}
		if !in.RequestCommon.FromCli(action, payload, &in) {
			return
		}
		out := b.UserLogout(&in)
		in.RequestCommon.ToCli(os.Stdout, out, out.ResponseCommon)

	}
}

// Code generated by 1_codegen_test.go DO NOT EDIT.
