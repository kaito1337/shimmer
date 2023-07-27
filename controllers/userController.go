package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"shimmer/models"
	"shimmer/services"
)

func SignInHandler(c *gin.Context) {
	requestBody := struct {
		Login    string
		Password string
	}{}
	_ = c.BindJSON(&requestBody)
	user, err := services.FindByLogin(requestBody.Login)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success auth", "data": user})
}
func SignUpHandler(c *gin.Context) {
	requestBody := struct {
		Login    string
		Password string
		Name     string
		Age      uint
	}{}
	_ = c.BindJSON(&requestBody)
	user, _ := services.FindByLogin(requestBody.Login)
	if user != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Login already taken"})
		return
	}
	if len(requestBody.Password) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Password length must be > 0"})
		return
	}
	if err := services.CreateUser(requestBody.Login, requestBody.Password, requestBody.Name, requestBody.Age); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Success created user"})
}

func UpdateHandler(c *gin.Context) {
	var requestBody *models.User
	_ = c.BindJSON(&requestBody)
	data, err := services.UpdateUser(requestBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": data, "message": "Success update user"})
}

func DeleteHandler(c *gin.Context) {
	var requestBody *models.UserResponse
	_ = c.BindJSON(&requestBody)
	if err := services.DeleteUser(requestBody.Login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success delete user"})
}
