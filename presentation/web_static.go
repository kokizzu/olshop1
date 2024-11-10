package presentation

import (
	"olshop1/domain"

	"github.com/gofiber/fiber/v2"
	//"github.com/kokizzu/gotro/C"
	"github.com/kokizzu/gotro/M"
)

func (w *WebServer) WebStatic(fw *fiber.App, d *domain.Domain) {

	fw.Get(`/`, func(c *fiber.Ctx) error {
		//	in, user, segments := userInfoFromContext(c, d)
		//	google := d.GuestExternalAuth(&domain.GuestExternalAuthIn{
		//		RequestCommon: in.RequestCommon,
		//		Provider:      domain.OauthGoogle,
		//	})
		return views.RenderIndex(c, M.SX{
			`title`: `Olshop1`,
			//`user`:   user,
			//`segments`:          segments,
		})
	})

	//fw.Get(`/debug`, func(ctx *fiber.Ctx) error {
	//	return views.RenderDebug(ctx, M.SX{})
	//})
}

//func notLogin(ctx *fiber.Ctx, d *domain.Domain, in domain.RequestCommon) bool {
//	var check domain.ResponseCommon
//	sess := d.MustLogin(in, &check)
//	if sess == nil {
//		_ = views.RenderError(ctx, M.SX{
//			`error`: check.Error,
//		})
//		return true
//	}
//	return false
//}
//
//func notAdmin(ctx *fiber.Ctx, d *domain.Domain, in domain.RequestCommon) bool {
//	var check domain.ResponseCommon
//	sess := d.MustLogin(in, &check)
//	if sess == nil {
//		_ = views.RenderError(ctx, M.SX{
//			`error`: check.Error,
//		})
//		return true
//	}
//	return false
//}
//
//func userInfoFromContext(c *fiber.Ctx, d *domain.Domain) (domain.UserProfileIn, *rqAuth.Users, M.SB) {
//	var in domain.UserProfileIn
//	err := webApiParseInput(c, &in.RequestCommon, &in, domain.UserProfileAction)
//	var user *rqAuth.Users
//	segments := M.SB{}
//	if err == nil {
//		out := d.UserProfile(&in)
//		user = out.User
//		segments = out.Segments
//	}
//	return in, user, segments
//}
//
//func userInfoFromRequest(rc domain.RequestCommon, d *domain.Domain) (*rqAuth.Users, M.SB) {
//	var user *rqAuth.Users
//	segments := M.SB{}
//	out := d.UserProfile(&domain.UserProfileIn{
//		RequestCommon: rc,
//	})
//	user = out.User
//	segments = out.Segments
//	return user, segments
//}
