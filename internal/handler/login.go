package handler


import (
	"strings"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)


// @Summary      Log in a user
// @Description  Authenticates a user with email and password, returning user details upon success.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginPayload  true  "User login credentials"
// @Success      200 {object} map[string]interface{} "Login successful"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      403 {object} map[string]interface{} "Invalid credentials"
// @Router       /auth/login [post]
func LoginUser(ctx *gin.Context) {
	var payload model.LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error()})
		return
	}

	Email := strings.ToLower(strings.TrimSpace(payload.Email))

	userID, emailExists := model.UserEmails[Email]
	if !emailExists {
		ctx.JSON(403, gin.H{
			"success": false,
			"message": "Email atau password salah",
			"error":   "invalid_credentials"})
		return
	}

	user := model.Users[userID]


	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		ctx.JSON(403, gin.H{
			"success": false,
			"message": "Email atau password salah",
			"error":   "invalid_credentials"})
		return
	}
	responseUser := user
	responseUser.Password = ""
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    responseUser})
}