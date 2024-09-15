package model

type User struct {
	account
	Role   string `json:"role" gorm:"not null"`
	UserId string `json:"user_id" gorm:"unique;not null"`
}

type UserCreateRequest struct {
	accountCreateRequest
	Role   string `json:"role"`
	UserId string `json:"user_id"`
}

type Users []User

func (ucr *UserCreateRequest) ToUser() User {
	return User{
		account: account{
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

func (u *User) Create() (User, error) {
	result := MainDB.Model(&User{}).Create(u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return *u, nil
}

func (u *Users) GetAll(filter interface{}) (Users, error) {
	result := MainDB.Model(&User{}).Where(filter).Preload("Account").Find(u)
	if result.Error != nil {
		return Users{}, result.Error
	}
	return *u, nil
}

func (u *User) GetOne(filter interface{}) (User, error) {
	result := MainDB.Model(&User{}).Where(filter).First(u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return *u, nil
}

func (u *User) Update(filter interface{}) (User, error) {
	result := MainDB.Model(&User{}).Where(filter).Updates(u)
	if result.Error != nil {
		return User{}, result.Error
	}
	return *u, nil
}

func (u *User) Delete(filter interface{}) error {
	result := MainDB.Model(&User{}).Where(filter).Delete(u)
	return result.Error
}
