package handler

import (
	"log"
	"net/http"
	"strconv"

	"menu_admin/authorization"
	"menu_admin/controller"
	"menu_admin/http/sessionsCookie"
	"menu_admin/model"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func createCookie(r *http.Request, w *http.ResponseWriter, token string) error {
	s, err := sessionsCookie.Cookie().Get(r, "session")
	if err != nil {
		return err
	}
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	s.Values["token"] = token
	err = s.Save(r, *w)
	return err
}

func GetUser(c echo.Context) error {
	var err error
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return err
	}
	ms, err := controller.GetUser(uint(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ms)
}

func GetAllUser(c echo.Context) error {
	ms, err := controller.GetAllUser()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ms)
}

func CreateUser(c echo.Context) error {
	var err error
	m := &model.User{}

	if err = c.Bind(m); err != nil {
		log.Fatal("error", m, err, c.Request())
		return err
	}
	err = controller.CreateUser(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func UpdateUserPassword(c echo.Context) error {
	claim, ok := c.Get("claim").(model.Claim)
	if !ok {
		return authorization.ErrCannotGetClaim
	}
	m := &model.User{}
	err := c.Bind(m)
	if err != nil {
		return err
	}
	err = controller.UpdateUserPassword(claim.UserID, m.Password)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func DeleteUser(c echo.Context) error {
	var err error
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return err
	}

	err = controller.DeleteUser(uint(id))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func LoginUser(c echo.Context) error {
	var err error
	m := &model.Login{}
	if err = c.Bind(m); err != nil {
		return err
	}

	user, err := controller.Login(m)
	if err != nil {
		return err
	}

	token, err := authorization.GenerateToken(&user)
	if err != nil {
		return err
	}

	err = createCookie(c.Request(), &c.Response().Writer, token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func GetMyUser(c echo.Context) error {
	claim, ok := c.Get("claim").(model.Claim)
	if !ok {
		return authorization.ErrCannotGetClaim
	}

	u, err := controller.GetUser(claim.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}
