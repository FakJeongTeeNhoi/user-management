package service

import (
	"os"
	"time"

	"github.com/FakJeongTeeNhoi/user-management/model"

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
		"is_verify":    account.Is_verify,
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

	ttl := ParseToInt(os.Getenv("JWT_TTL"))
	claims["exp"] = time.Now().Add(time.Second * time.Duration(ttl)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))

}

func GetInfoFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
