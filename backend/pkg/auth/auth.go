package auth

import (
	"errors"
	"fmt"

	"github.com/fabiokaelin/f-oauth/pkg/db"
	"github.com/fabiokaelin/f-oauth/pkg/notification"
	"github.com/fabiokaelin/f-oauth/pkg/user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

func RegisterUser(newUser user.User) (user.User, error) {

	ifExist, err := db.DoesUserExist(newUser.Email)
	if err != nil {
		fmt.Println("error, google, failed to get first time user", err)
		return user.User{}, errors.New("error appeared")
	}

	if ifExist {
		return user.User{}, errors.New("user with that email already exists")
	}

	dbuser := db.DatabaseUser{
		Name:      newUser.Name,
		Email:     newUser.Email,
		Password:  newUser.Password,
		Role:      newUser.Role,
		Photo:     newUser.Photo,
		Verified:  newUser.Verified,
		Provider:  newUser.Provider,
		CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	id, err := db.CreateUser(dbuser)
	if err != nil {
		fmt.Println("error, basic, failed to create user", err)
		return user.User{}, errors.New("error appeared")
	}
	newUser.ID, err = uuid.FromString(id)
	if err != nil {
		fmt.Println("error, basic, failed to create user", err)
		return user.User{}, errors.New("error appeared")
	}

	notification.NewUserNotification(newUser)

	return newUser, nil
}

func LoginUser(email string, password string) (user.User, error) {

	ifExist, err := db.DoesUserExist(email)
	if err != nil {
		fmt.Println(err)
		return user.User{}, errors.New("error appeared")
	}

	if !ifExist {
		return user.User{}, errors.New("invalid email or password")
	}

	currentUser, err := db.GetUserByEmail(email)
	if err != nil {
		fmt.Println(err)
		return user.User{}, errors.New("error appeared")
	}

	if currentUser.Provider != "local" {
		return user.User{}, errors.New("you have already signed up with a different method")
	}

	equal := ComparePasswords(currentUser.Password, []byte(password))
	if !equal {
		return user.User{}, errors.New("invalid email or password")
	}

	currentUserNew := user.User{
		ID:       uuid.FromStringOrNil(currentUser.ID),
		Name:     currentUser.Name,
		Email:    currentUser.Email,
		Password: currentUser.Password,
		Role:     currentUser.Role,
		Photo:    currentUser.Photo,
		Verified: currentUser.Verified,
		Provider: currentUser.Provider,
		// CreatedAt: currentUser.CreatedAt,
		// UpdatedAt: currentUser.UpdatedAt,
	}

	return currentUserNew, nil
}
