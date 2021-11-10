package controller

import (
	"errors"
	"net/mail"
	"unicode"

	"menu_admin/authorization"
	"menu_admin/model"
	"menu_admin/storage"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyInUsed = errors.New("email already in used")
	ErrInvalidEmail       = errors.New("invalid email")
)

func hashAndSalt(pwd []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// GetUser return an user by ID
func GetUser(id uint) (model.User, error) {
	p := model.User{}
	err := storage.DB().First(&p, id).Error
	return p, err
}

// GetAllUser return all users
func GetAllUser() ([]model.User, error) {
	ps := make([]model.User, 0)
	r := storage.DB().Find(&ps)
	return ps, r.Error
}

// CreateUser create a new user, encrypt the password and send a confirmation code to the email
func CreateUser(m *model.User) error {
	var err error
	if err = isEmailAndPasswordValid(m.Email, m.Password); err != nil {
		return err
	}

	u := &model.User{}
	result := storage.DB().Where("email = ?", m.Email).First(u)
	if result.RowsAffected != 0 {
		return ErrEmailAlreadyInUsed
	}
	err = result.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	pwd, err := hashAndSalt([]byte(m.Password))
	if err != nil {
		return err
	}

	m.Password = string(pwd)
	r := storage.DB().Create(m)
	return r.Error
}

// ValidateUser receives a token, validates it and updates the user status to confirmed
func ValidateUser(token string) error {
	claim, err := authorization.ValidateCodeVerification(token)
	if err != nil {
		return err
	}
	return storage.DB().Model(&model.User{}).Where("email = ?", claim.Email).Update("is_confirmated", true).Error
}

// DeleteUser use soft delete to remove an user
func DeleteUser(id uint) error {
	r := storage.DB().Delete(&model.User{}, id)
	return r.Error
}

// Login Receive the username and password of a user, confirm that the credentials are correct and return a user
func Login(m *model.Login) (model.User, error) {
	user := model.User{}
	var err error
	if err = isEmailAndPasswordValid(m.Email, m.Password); err != nil {
		return user, err
	}

	err = storage.DB().First(&user,
		&model.User{
			Email: m.Email,
		}).Error
	if err != nil {
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(m.Password)); err != nil {
		return model.User{}, err
	}
	return user, nil
}

// isEmailValid return true if the email is valid, else return false
func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// isEmailAndPasswordValid return an error if the password or emails are invalid
func isEmailAndPasswordValid(email, password string) error {
	if !isEmailValid(email) {
		return ErrInvalidEmail
	}
	if !isPasswordValid(password) {
		return ErrInvalidPassword
	}
	return nil
}

func isPasswordValid(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}

	var (
		hasUpperCase bool
		hasSpecial   bool
		hasNumber    bool
		hasLower     bool
	)

	for _, v := range pwd {
		if hasLower && hasNumber && hasSpecial && hasUpperCase {
			return true
		}
		switch {
		case unicode.IsLower(v):
			hasLower = true
		case unicode.IsUpper(v):
			hasUpperCase = true
		case unicode.IsNumber(v):
			hasNumber = true
		case unicode.IsPunct(v) || unicode.IsSymbol(v):
			hasSpecial = true
		}
	}

	return hasLower && hasNumber && hasSpecial && hasUpperCase
}

// UpdateUserPassword update a user's password by ID
func UpdateUserPassword(id uint, password string) error {
	if !isPasswordValid(password) {
		return ErrInvalidPassword
	}
	m := &model.User{}
	m.ID = id

	pwd, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	return storage.DB().Model(m).Updates(model.User{
		Password: string(pwd),
	}).Error
}
