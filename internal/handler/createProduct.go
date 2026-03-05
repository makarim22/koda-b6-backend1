package handler

import (
	"fmt"
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
)

// @Summary      Create a new product
// @Description  Creates a new product with the provided details.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      Product  true  "Product object to be created"
// @Success      201 {object} map[string]interface{} "Product created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Router       /products [post]
func CreateProduct(ctx *gin.Context) {
	var newProduct model.Product

	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error()})
		return
	}

	newProduct.ID = model.NextProductId
	model.Products[newProduct.ID] = newProduct

	model.NextProductId++

	fmt.Printf("Added product: %+v\n", newProduct)
	fmt.Println("Current products map:", model.Products)

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Product created successfully",
		"data":    newProduct,
	})
}