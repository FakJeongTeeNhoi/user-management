package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	Name      string `json:"name" gorm:"not null"`
	Faculty   string `json:"faculty" gorm:"not null"`
	Type      string `json:"type" gorm:"not null"`
	Is_verify bool   `json:"is_verify" gorm:"default:false"`
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
	Name     string `json:"name" binding:"required"`
	Faculty  string `json:"faculty" binding:"required"`
	Type     string `json:"type" binding:"required"`
}

type accountUpdateRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Faculty string `json:"faculty"`
	Type    string `json:"type"`
}

func (a *Account) GetOne(filter interface{}) error {
	result := MainDB.Model(&Account{}).Where(filter).First(a)
	return result.Error
}

func (a *Account) Create() error {
	result := MainDB.Model(&Account{}).Create(a)
	return result.Error
}

func (a *Account) Update() error {
	result := MainDB.Model(&Account{}).Where("id = ?", a.ID).Updates(a)
	return result.Error
}

func (a *Account) Delete() error {
	result := MainDB.Model(&Account{}).Where("id = ?", a.ID).Unscoped().Delete(
		map[string]interface{}{
			"id": a.ID,
		})
	return result.Error
}
