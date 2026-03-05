package handler 

import (
	"fmt"
	"strings"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
)


// @Summary      Create a new user
// @Description  Creates a new user with the provided details.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User object to be created"
// @Success      201 {object} map[string]interface{} "User created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid input or validation error"
// @Failure      409 {object} map[string]interface{} "Email already registered"
// @Router       /users [post]
func CreateUser(ctx *gin.Context) {
	var newUser model.User
	fmt.Println("newUser", newUser)

	err := ctx.ShouldBindJSON(&newUser)

	if err != nil {
		ctx.JSON(400, gin.H{
			"Success": false,
			"Message": "Gagal membuat user",
		})
		return
	}

	name := strings.TrimSpace(newUser.Name)
	email := strings.ToLower(strings.TrimSpace(newUser.Email))
	password := strings.TrimSpace(newUser.Password)

	if name == "" || email == "" || password == "" {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Nama, email, dan password tidak boleh kosong",
			"error":   "validation_error",
		})
		return
	}

	for _, existingUser := range model.Users {
		if existingUser.Email == email {
			ctx.JSON(409, gin.H{
				"success": false,
				"message": "Email sudah terdaftar",
				"error":   "duplicate_email",
			})
			return
		}
	}

	newUser.ID = model.NextID
	model.NextID++
	model.Users[newUser.ID] = newUser

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Berhasil membuat user",
		"data":    newUser,
	})
}