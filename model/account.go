package model

import "gorm.io/gorm"

type account struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Faculty  string `json:"faculty"`
	Type     string `json:"type"`
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Faculty  string `json:"faculty"`
	Type     string `json:"type"`
}
