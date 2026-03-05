package handler

import (
	"github.com/gin-gonic/gin"
	"golang-backend/internal/model"
)

// @Summary      Get all users
// @Description  Retrieves a list of all registered users.
// @Tags         users
// @Produce      json
// @Success      200 {object} map[string]interface{} "Successfully retrieved users"
// @Failure      404 {object} map[string]interface{} "No users found"
// @Router       /users [get]
func GetUsers(ctx *gin.Context) {
	if len(model.Users) == 0 {
		ctx.JSON(404, gin.H{
			"error": "No users found",
		})
		return
	}

	userList := make([]model.User, 0, len(model.Users))
	for _, user := range model.Users {
		userList = append(userList, user)
	}

	ctx.JSON(200, gin.H{
		"data":  userList,
		"count": len(userList),
	})
}