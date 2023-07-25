package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"os"
)

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
	Name     string
	Age      uint `gorm:"check:age >= 18"`
	IsActive bool `gorm:"default:true"`
}
type UserResponse struct {
	ID       uint `json:"-"`
	Login    string
	Name     string
	Age      uint
	IsActive bool
	Password string `json:"-"`
}

var db *gorm.DB

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if dbData, err := gorm.Open(postgres.Open(os.Getenv("DATA_SOURCE_NAME")), &gorm.Config{}); err != nil {
		log.Fatal(err)
	} else {
		db = dbData
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}
	var router *gin.Engine = gin.Default()
	router.GET("/signIn", signInHandler)
	router.POST("/signUp", signUpHandler)
	router.POST("/updateUser", updateHandler)
	router.POST("/deleteUser", deleteHandler)
	_ = router.Run(":8080")
}

func signInHandler(c *gin.Context) {
	requestBody := struct {
		Login    string
		Password string
		Name     string
	}{}
	var user UserResponse
	_ = c.BindJSON(&requestBody)
	result := db.Model(&User{}).Where("login = ?", requestBody.Login).First(&user)
	if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "zero results"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success auth", "data": user})
}
func signUpHandler(c *gin.Context) {
	requestBody := struct {
		Login    string
		Password string
		Name     string
		Age      uint
	}{}
	_ = c.BindJSON(&requestBody)
	if len(requestBody.Password) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User creation error"})
		return
	}
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 12)
	if err := createUser(requestBody.Login, string(hashedPass), requestBody.Name, requestBody.Age); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Success created user"})
}

func updateHandler(c *gin.Context) {
	var requestBody *User
	_ = c.BindJSON(&requestBody)
	data, err := updateUser(requestBody.Login, requestBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": data, "message": "Success update user"})
}

func deleteHandler(c *gin.Context) {
	var requestBody *UserResponse
	_ = c.BindJSON(&requestBody)
	if err := deleteUser(requestBody.Login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success delete user"})
}

func deleteUser(login string) error {
	_ = db.Model(&User{}).Where("login = ?", login).Updates(User{IsActive: false})
	result := db.Where("login = ?", login).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func updateUser(login string, data *User) (*UserResponse, error) {
	var hashedPass []byte
	var user *UserResponse
	if len(data.Password) > 0 {
		hashedPass, _ = bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	}

	result := db.Model(&User{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "login"}, {Name: "name"}, {Name: "age"}, {Name: "is_active"}}}).
		Where("login = ?", login).
		Updates(User{Name: data.Name, Age: data.Age, IsActive: data.IsActive, Password: string(hashedPass)}).
		First(&user)
	log.Println(user.Password)
	if result.RowsAffected == 0 {
		return nil, errors.New("wrong data")
	}
	return user, nil
}

func createUser(login string, password string, name string, age uint) error {
	user := User{
		Login:    login,
		Password: password,
		Name:     name,
		Age:      age,
		IsActive: true,
	}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
