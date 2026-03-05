package handler

import (
    "github.com/gin-gonic/gin"
    "golang-backend/internal/model"
)


// @Summary      Get all products
// @Description  Retrieves a list of all products.
// @Tags         products
// @Produce      json
// @Success      200 {object} map[string]interface{} "Successfully retrieved products"
// @Failure      404 {object} map[string]interface{} "No products found"
// @Router       /products [get]
func GetProducts(ctx *gin.Context) {
	if len(model.Products) == 0 {
		ctx.JSON(404, gin.H{
			"error": "No products found",
		})
		return
	}

	productList := make([]model.Product, 0, len(model.Products))
	for _, product := range model.Products {
		productList = append(productList, product)
	}

	ctx.JSON(200, gin.H{
		"data":  productList,
		"count": len(productList),
	})
}