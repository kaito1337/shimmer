package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"shimmer/controllers"
	"shimmer/models"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if DBData, err := gorm.Open(postgres.Open(os.Getenv("DATA_SOURCE_NAME")), &gorm.Config{}); err != nil {
		log.Fatal(err)
	} else {
		models.DB = DBData
	}
	if err := models.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	var router *gin.Engine = gin.Default()
	router.GET("/signIn", controllers.SignInHandler)
	router.POST("/signUp", controllers.SignUpHandler)
	router.POST("/updateUser", controllers.UpdateHandler)
	router.POST("/deleteUser", controllers.DeleteHandler)
	_ = router.Run(":8080")
}
