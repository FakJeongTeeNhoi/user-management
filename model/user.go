package model

import (
	"cmp"
	"gorm.io/gorm"
)

type User struct {
	Account
	Role   string `json:"role" gorm:"not null"`
	UserId string `json:"user_id" gorm:"unique;not null"`
}

type UserCreateRequest struct {
	accountCreateRequest
	Role   string `json:"role" gorm:"not null"`
	UserId string `json:"user_id" gorm:"unique;not null"`
}

type UserUpdateRequest struct {
	accountUpdateRequest
	Role   string `json:"role"`
	UserId string `json:"user_id"`
}

type Users []User

func (ucr *UserCreateRequest) ToUser() User {
	return User{
		Account: Account{
			Email:    ucr.Email,
			Password: ucr.Password,
			Name:     ucr.Name,
			Faculty:  ucr.Faculty,
			Type:     ucr.Type,
		},
		Role:   ucr.Role,
		UserId: ucr.UserId,
	}
}

func (uur *UserUpdateRequest) ToUser(u User) User {
	return User{
		Account: Account{
			Model: gorm.Model{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
				DeletedAt: u.DeletedAt,
			},
			Email:    u.Email,
			Password: u.Password,
			Name:     cmp.Or(uur.Name, u.Name),
			Faculty:  cmp.Or(uur.Faculty, u.Faculty),
			Type:     cmp.Or(uur.Type, u.Type),
		},
		Role:   cmp.Or(uur.Role, u.Role),
		UserId: cmp.Or(uur.UserId, u.UserId),
	}
}

func (u *User) Create() error {
	result := MainDB.Model(&User{}).Create(u)
	return result.Error
}

func (u *Users) GetAll(filter interface{}) error {
	result := MainDB.Model(&User{}).Where(filter).Find(u)
	return result.Error
}

func (u *User) GetOne(filter interface{}) error {
	result := MainDB.Model(&User{}).Where(filter).First(u)
	return result.Error
}

func (u *User) Update() error {
	result := MainDB.Model(&User{}).Where("id = ?", u.ID).Updates(u)
	return result.Error
}

func (u *User) Delete() error {
	result := MainDB.Model(&User{}).Delete(u)
	return result.Error
}
