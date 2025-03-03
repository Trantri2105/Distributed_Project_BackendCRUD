package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	request "user_crud/dto"
	"user_crud/middleware"
	"user_crud/model"
	"user_crud/service"
	"user_crud/utils"
)

func NewUserController(service service.UserService) *gin.Engine {
	r := gin.Default()
	j := utils.NewJwtUtils()
	authMiddleware := middleware.NewAuthMiddleware(j)
	r.POST("/register", func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err.Error())
			return
		}
		err := service.RegisterUser(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println(err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
	})

	r.POST("/login", func(c *gin.Context) {
		var req request.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := service.Login(c, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	r.PATCH("/update", authMiddleware.ValidateAndExtractJwt(), func(c *gin.Context) {
		claims, _ := c.Get(middleware.JWTClaimsContextKey)
		userClaims := claims.(jwt.MapClaims)
		id := int(userClaims["userId"].(float64))
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = id
		err := service.UpdateUser(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully!"})
	})

	r.GET("/profile", authMiddleware.ValidateAndExtractJwt(), func(c *gin.Context) {
		claims, _ := c.Get(middleware.JWTClaimsContextKey)
		userClaims := claims.(jwt.MapClaims)
		id := int(userClaims["userId"].(float64))
		user, err := service.GetUserById(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":     user.Id,
			"name":   user.Name,
			"email":  user.Email,
			"phone":  user.Phone,
			"gender": user.Gender,
		})
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Chinh chong",
		})
	})
	return r
}
