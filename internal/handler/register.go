package handler


import (
	"strings"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Register a new user
// @Description  Registers a new user and hashes their password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User registration details"
// @Success      200 {object} map[string]interface{} "User registered successfully"
// @Failure      400 {object} map[string]interface{} "Invalid input or validation error"
// @Failure      409 {object} map[string]interface{} "Email already registered"
// @Router       /auth/register [post]
func RegisterUser(ctx *gin.Context) {
	var newUser model.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error()})
		return
	}

	Name := strings.TrimSpace(newUser.Name)
	Email := strings.ToLower(strings.TrimSpace(newUser.Email))

	if Name == "" || Email == "" || newUser.Password == "" {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Nama, email, dan password tidak boleh kosong",
			"error":   "validation_error"})
		return
	}

	if !model.EmailRegex.MatchString(Email) {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Format email tidak valid",
			"error":   "invalid_email_format"})
		return
	}

	if _, emailExists := model.UserEmails[newUser.Email]; emailExists {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Email sudah terdaftar",
			"error":   "duplicate_email"})
		return
	}

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}
	newUser.Password = string(HashedPassword)

	newUser.ID = model.NextID
	model.NextID++
	model.Users[newUser.ID] = newUser
	model.UserEmails[newUser.Email] = newUser.ID

	responseUser := newUser
	responseUser.Password = "" // hilangkan password
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data":    responseUser})
}