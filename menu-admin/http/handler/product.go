package handler

import (
	"net/http"
	"strconv"

	"menu_admin/controller"
	"menu_admin/model"

	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	var err error
	m := &model.Product{}
	if err = c.Bind(m); err != nil {
		return err
	}

	err = controller.CreateProduct(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m)
}

func UpdateProduct(c echo.Context) error {
	var err error
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return err
	}

	m := model.Product{}

	if err = c.Bind(&m); err != nil {
		return err
	}
	m.ID = uint(id)
	err = controller.UpdateProduct(&m)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, m)
}

func DeleteProduct(c echo.Context) error {
	var err error
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return err
	}

	err = controller.DeleteProduct(uint(id))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
