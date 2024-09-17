package model

import (
	"cmp"
	"gorm.io/gorm"
    "github.com/google/uuid"
)


type Staff struct {
	Account   Account `json:"account" gorm:"foreignKey:AccountId"`
	AccountId uint    `json:"account_id" gorm:"not null"`
	Space_list []uuid.UUID `json:"space_list" gorm:"type:uuid[]"`
}

type StaffCreateRequest struct {
    accountCreateRequest
}

type StaffUpdateRequest struct {
	accountUpdateRequest
}

type Staffs []Staff

func (stf *StaffCreateRequest) ToStaff() Staff {
	return Staff{
        Account: Account{
            Email:    stf.Email,
            Password: stf.Password,
            Name:     stf.Name,
            Faculty:  stf.Faculty,
            Type:     stf.Type,
        },
        Space_list: []uuid.UUID{},
    }
}

func (sur *StaffUpdateRequest) ToStaff(s Staff) Staff {
	return Staff{
		Account: Account{
			Model: gorm.Model{
				ID:        s.Account.ID,
				CreatedAt: s.Account.CreatedAt,
				UpdatedAt: s.Account.UpdatedAt,
				DeletedAt: s.Account.DeletedAt,
			},
			Email:    s.Account.Email,
			Password: s.Account.Password,
			Name:     cmp.Or(sur.Name, s.Account.Name),
			Faculty:  cmp.Or(sur.Faculty, s.Account.Faculty),
			Type:     cmp.Or(sur.Type, s.Account.Type),
		},
		AccountId: s.AccountId,
		Space_list: s.Space_list,
	}
}

func (s *Staff) Create() error {
	result := MainDB.Model(&Staff{}).Create(s)
	return result.Error
}

func (s *Staffs) GetAll(filter interface{}) error {
    result := MainDB.Model(&Staff{}).Where(filter).Preload("Account").Find(s)
    return result.Error
}

func (s *Staff) GetOne(filter interface{}) error {
	result := MainDB.Model(&Staff{}).Where(filter).Preload("Account").First(s)
	return result.Error
}

func (s *Staff) Update() error {
	return MainDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Staff{}).Where("account_id = ?", s.AccountId).Updates(s).Error; err != nil {
			return err
		}
		if err := tx.Model(&Account{}).Where("id = ?", s.AccountId).Updates(s.Account).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *Staff) Delete() error {
	return MainDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Staff{}).Where("account_id = ?", s.AccountId).Delete(Staff{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&Account{}).Where("id = ?", s.AccountId).Unscoped().Delete(
			map[string]interface{}{
				"id": s.AccountId,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}