package handler

import (
	"strconv"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
)

// @Summary      Get a user by ID
// @Description  Retrieves a single user by their ID.
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} User "User found"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /users/{id} [get]
func GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, exists := model.Users[userId]
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "tidak dapat menemukan user",
		})
		return
	}

	ctx.JSON(200, user)
}