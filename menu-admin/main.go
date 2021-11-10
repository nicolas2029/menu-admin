package main

import (
	"errors"
	"log"
	"menu_admin/authorization"
	"menu_admin/controller"
	"menu_admin/http/route"
	"menu_admin/http/sessionsCookie"
	"menu_admin/model"
	"menu_admin/storage"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newUser() error {
	email, isEnv := os.LookupEnv("MENU-OWNER-EMAIL")
	if !isEnv {
		return errors.New("environment variable (MENU-OWNER-EMAIL) not found")
	}
	password, isEnv := os.LookupEnv("MENU-OWNER-PASSWORD")
	if !isEnv {
		return errors.New("environment variable (MENU-OWNER-PASSWORD) not found")
	}
	user := model.User{
		Email:    email,
		Password: password,
	}
	return controller.CreateUser(&user)
}

func newDB() error {
	err := storage.DB().First(&model.User{}).Error
	if err != nil {
		err = newUser()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := authorization.LoadCertificates()
	if err != nil {
		log.Fatalf("no se pudo cargar los certificados: %v", err)
	}
	storage.New()
	err = sessionsCookie.NewCookieStore()
	if err != nil {
		log.Fatalf("no se pudo cargar los certificados para cookies: %v", err)
	}

	err = storage.DB().AutoMigrate(
		&model.Product{},
		&model.User{},
	)
	if err != nil {
		log.Fatalf("no se realizararon las migraciones: %v", err)
	}
	err = newDB()
	if err != nil {
		log.Fatalf("error al crear la base de datos: %v", err)
	}
	e := echo.New()
	e.Use(middleware.Recover())
	//e.Use(middleware.Logger())
	e.Use(session.Middleware(sessionsCookie.Cookie()))
	route.All(e)
	err = e.Start(":3000")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
