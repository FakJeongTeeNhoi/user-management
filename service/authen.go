package service

import (
	"github.com/FakJeongTeeNhoi/user-management/model"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ValidateCredential(request model.LoginRequest) (model.Account, error) {
	account := model.Account{}
	err := account.GetOne(map[string]interface{}{"email": request.Email})
	if err != nil {
		return account, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))

	return account, err
}

func GenerateToken(userType string, account model.Account) (string, error) {
	claims := jwt.MapClaims{
		"account_id":   account.ID,
		"account_type": userType,
		"email":        account.Email,
		"name":         account.Name,
		"faculty":      account.Faculty,
		"type":         account.Type,
	}

	if userType == "staff" {
		staff := model.Staff{}
		err := staff.GetOne(map[string]interface{}{"account_id": account.ID})
		if err != nil {
			return "", err
		}
	} else if userType == "user" {
		user := model.User{}
		err := user.GetOne(map[string]interface{}{"account_id": account.ID})
		if err != nil {
			return "", err
		}
		claims["role"] = user.Role
		claims["user_id"] = user.UserId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("JWT_SECRET"))
}
