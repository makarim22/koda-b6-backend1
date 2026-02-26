package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

var users = map[int]User{
	1: {ID: 1, Name: "Budi", Email: "budi@email.com", Password: "hashed123"},
	2: {ID: 2, Name: "Siti", Email: "siti@email.com", Password: "hashed456"},
}

var nextID = 3


func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("Hello!"))
	})

	r.GET("/user", func(ctx *gin.Context) {
       
		if len (users) == 0 {
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
	})

	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Invalid user ID",
			})
			return
		}

		user , exists := users[userId]
		if !exists {
          ctx.JSON(400, gin.H{
			"error" : "tidak dapat menemukan user",
		  })
		  return
		}

		
	    ctx.JSON(200, user)

	})

	r.POST("/users", func(ctx *gin.Context) {
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

		for _ , existingUser := range users {
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
		users[newUser.ID]= newUser

		ctx.JSON(201, gin.H{ 
			"success": true,
			"message": "Berhasil membuat user",
			"data":    newUser, 
		})

	})

	r.Run("localhost:8888")
}