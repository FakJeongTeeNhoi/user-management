package model

import (
	"cmp"
	"gorm.io/gorm"
)

type User struct {
	Account   Account `json:"account" gorm:"foreignKey:AccountId"`
	AccountId uint    `json:"account_id" gorm:"not null"`
	DeletedAt gorm.DeletedAt
	Role      string `json:"role" gorm:"not null"`
	UserId    string `json:"user_id" gorm:"unique;not null"`
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
				ID:        u.Account.ID,
				CreatedAt: u.Account.CreatedAt,
				UpdatedAt: u.Account.UpdatedAt,
				DeletedAt: u.Account.DeletedAt,
			},
			Email:    u.Account.Email,
			Password: u.Account.Password,
			Name:     cmp.Or(uur.Name, u.Account.Name),
			Faculty:  cmp.Or(uur.Faculty, u.Account.Faculty),
			Type:     cmp.Or(uur.Type, u.Account.Type),
		},
		AccountId: u.AccountId,
		Role:      cmp.Or(uur.Role, u.Role),
		UserId:    cmp.Or(uur.UserId, u.UserId),
	}
}

func (u *User) Create() error {
	result := MainDB.Model(&User{}).Create(u)
	return result.Error
}

func (u *Users) GetAll(filter interface{}) error {
	result := MainDB.Model(&User{}).Where(filter).Preload("Account").Find(u)
	return result.Error
}

func (u *User) GetOne(filter interface{}) error {
	result := MainDB.Model(&User{}).Where(filter).Preload("Account").First(u)
	return result.Error
}

func (u *User) Update() error {
	return MainDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("account_id = ?", u.AccountId).Updates(u).Error; err != nil {
			return err
		}
		if err := tx.Model(&Account{}).Where("id = ?", u.AccountId).Updates(u.Account).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *User) Delete() error {
	return MainDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("account_id = ?", u.AccountId).Delete(User{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&Account{}).Where("id = ?", u.AccountId).Delete(
			map[string]interface{}{
				"id": u.AccountId,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}
