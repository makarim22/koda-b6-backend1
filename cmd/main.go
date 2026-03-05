package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "golang-backend/docs"
)

// --- Models ---

// User represents a user in the system.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

// Product represents a product in the system.
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	VariantId   int    `json:"variant_id"`
	SizeId      int    `json:"size_id"`
}

// LoginPayload represents the request body for user login.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Global variables for simplicity
var users = map[int]User{
	1: {ID: 1, Name: "Budi", Email: "budi@email.com", Password: "hashed123"},
	2: {ID: 2, Name: "Siti", Email: "siti@email.com", Password: "hashed456"},
}

var products = map[int]Product{}

var nextID = 3

var userEmails = map[string]int{}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

var nextProductId = 1

// --- HANDLER FUNCTIONS ---

// @Summary      Get all users
// @Description  Retrieves a list of all registered users.
// @Tags         users
// @Produce      json
// @Success      200 {object} map[string]interface{} "Successfully retrieved users"
// @Failure      404 {object} map[string]interface{} "No users found"
// @Router       /users [get]
func getUsers(ctx *gin.Context) {
	if len(users) == 0 {
		ctx.JSON(404, gin.H{
			"error": "No users found",
		})
		return
	}

	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	ctx.JSON(200, gin.H{
		"data":  userList,
		"count": len(userList),
	})
}

// @Summary      Get a user by ID
// @Description  Retrieves a single user by their ID.
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} User "User found"
// @Failure      400 {object} map[string]interface{} "Invalid user ID"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /users/{id} [get]
func getUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, exists := users[userId]
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "tidak dapat menemukan user",
		})
		return
	}

	ctx.JSON(200, user)
}

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
func createUser(ctx *gin.Context) {
	var newUser User
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

	for _, existingUser := range users {
		if existingUser.Email == email {
			ctx.JSON(409, gin.H{
				"success": false,
				"message": "Email sudah terdaftar",
				"error":   "duplicate_email",
			})
			return
		}
	}

	newUser.ID = nextID
	nextID++
	users[newUser.ID] = newUser

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Berhasil membuat user",
		"data":    newUser,
	})
}

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
func updateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	userID, err := strconv.Atoi(id)
	fmt.Println("usernya", userID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, exists := users[userID]
	if !exists {
		ctx.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	}

	var updateData User
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

	users[userID] = user

	ctx.JSON(200, gin.H{
		"message": "berhasil mengupdate user",
		"data":    user,
	})
}

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
func registerUser(ctx *gin.Context) {
	var newUser User

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

	if !emailRegex.MatchString(Email) {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Format email tidak valid",
			"error":   "invalid_email_format"})
		return
	}

	if _, emailExists := userEmails[newUser.Email]; emailExists {
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

	newUser.ID = nextID
	nextID++
	users[newUser.ID] = newUser
	userEmails[newUser.Email] = newUser.ID

	responseUser := newUser
	responseUser.Password = "" // hilangkan password
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data":    responseUser})
}

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
func loginUser(ctx *gin.Context) {
	var payload LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error()})
		return
	}

	Email := strings.ToLower(strings.TrimSpace(payload.Email))

	userID, emailExists := userEmails[Email]
	if !emailExists {
		ctx.JSON(403, gin.H{
			"success": false,
			"message": "Email atau password salah",
			"error":   "invalid_credentials"})
		return
	}

	user := users[userID]
	fmt.Println("user", user)

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

// @Summary      Create a new product
// @Description  Creates a new product with the provided details.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      Product  true  "Product object to be created"
// @Success      201 {object} map[string]interface{} "Product created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Router       /products [post]
func createProduct(ctx *gin.Context) {
	var newProduct Product

	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error()})
		return
	}

	newProduct.ID = nextProductId
	products[newProduct.ID] = newProduct

	nextProductId++

	fmt.Printf("Added product: %+v\n", newProduct)
	fmt.Println("Current products map:", products)

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Product created successfully",
		"data":    newProduct,
	})
}

// @Summary      Get all products
// @Description  Retrieves a list of all products.
// @Tags         products
// @Produce      json
// @Success      200 {object} map[string]interface{} "Successfully retrieved products"
// @Failure      404 {object} map[string]interface{} "No products found"
// @Router       /products [get]
func getProducts(ctx *gin.Context) {
	if len(products) == 0 {
		ctx.JSON(404, gin.H{
			"error": "No products found",
		})
		return
	}

	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}

	ctx.JSON(200, gin.H{
		"data":  productList,
		"count": len(productList),
	})
}

// @title           Simple API
// @version         1.0
// @description     Simple API to learn Swaggo
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8888
// @BasePath  /
func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("Hello!"))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// User Endpoints
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)

	// Auth Endpoints
	r.POST("/auth/register", registerUser)
	r.POST("/auth/login", loginUser)

	// Product Endpoints
	r.POST("/products", createProduct)
	r.GET("/products", getProducts)

	r.Run("localhost:8888")
}