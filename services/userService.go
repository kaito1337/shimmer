package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
	"shimmer/models"
)

func CreateUser(login string, password string, name string, age uint) error {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := models.User{
		Login:    login,
		Password: string(hashedPass),
		Name:     name,
		Age:      age,
		IsActive: true,
	}
	result := models.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FindByLogin(login string) (*models.UserResponse, error) {
	var user *models.UserResponse
	result := models.DB.Model(&models.User{}).Where("login = ?", login).First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New("wrong credentials")
	} else {
		return user, nil
	}
}

func UpdateUser(data *models.User) (*models.UserResponse, error) {
	var hashedPass []byte
	var user *models.UserResponse
	if len(data.Password) > 0 {
		hashedPass, _ = bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	}

	result := models.DB.Model(&models.User{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "login"}, {Name: "name"}, {Name: "age"}, {Name: "is_active"}}}).
		Where("login = ?", data.Login).
		Updates(models.User{Name: data.Name, Age: data.Age, IsActive: data.IsActive, Password: string(hashedPass)}).
		First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New("wrong data")
	}
	return user, nil
}

func DeleteUser(login string) error {
	_ = models.DB.Model(&models.User{}).Where("login = ?", login).Updates(models.User{IsActive: false})
	result := models.DB.Where("login = ?", login).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
