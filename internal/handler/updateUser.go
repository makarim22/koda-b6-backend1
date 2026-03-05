package handler

import (
	"fmt"
	"strings"
	"strconv"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)


// @Summary      Update an existing user
// @Description  Updates an existing user's details by ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        user body      User  true  "User object with updated fields"
// @Success      200 {object} map[string]interface{} "User updated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid user ID or request body"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /users/{id} [put]
func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	userID, err := strconv.Atoi(id)
	fmt.Println("usernya", userID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, exists := model.Users[userID]
	if !exists {
		ctx.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	}

	var updateData model.User
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	name := strings.TrimSpace(updateData.Name)
	email := strings.ToLower(strings.TrimSpace(updateData.Email))

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	password := string(HashedPassword)

	fmt.Println("name", name)
	fmt.Println("email", email)
	fmt.Println("password", password)

	user.Name = name
	user.Email = email
	user.Password = password

	model.Users[userID] = user

	ctx.JSON(200, gin.H{
		"message": "berhasil mengupdate user",
		"data":    user,
	})
}