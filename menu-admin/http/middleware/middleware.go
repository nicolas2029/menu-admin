package middleware

import (
	"errors"

	"menu_admin/authorization"
	"menu_admin/http/sessionsCookie"
	"menu_admin/model"

	"github.com/labstack/echo/v4"
)

var (
	ErrUserNotLogin = errors.New("user not login")
)

func authorizeLogin(c echo.Context) (model.Claim, error) {
	var err error
	if v := c.Request().Header.Get("authorization"); v != "" {
		m, err := authorization.ValidateToken(v)
		if err != nil {
			return model.Claim{}, err
		}
		return m, nil
	}
	cookie := sessionsCookie.Cookie()
	sess, err := cookie.Get(c.Request(), "session")
	if err != nil {
		return model.Claim{}, err
	}
	v, f := sess.Values["token"]
	if !f {
		return model.Claim{}, ErrUserNotLogin
	}
	m, err := authorization.ValidateToken(v.(string))
	if err != nil {
		return model.Claim{}, err
	}
	c.Request().Header.Set("authorization", v.(string))
	return m, nil
}

func DeleteSession(c echo.Context) error {
	cookie := sessionsCookie.Cookie()
	sess, err := cookie.Get(c.Request(), "session")
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	return cookie.Save(c.Request(), c.Response(), sess)
}

func AuthorizeIsLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := authorizeLogin(c)
		if err != nil {
			return err
		}
		c.Set("claim", m)
		return next(c)
	}
}
