package model

// Product represents a product in the system.
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	VariantId   int    `json:"variant_id"`
	SizeId      int    `json:"size_id"`
}


var Products = map[int]Product{}

var NextProductId = 1