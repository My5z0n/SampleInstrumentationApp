package model

type ProductDetailsModel struct {
	ProductName string `uri:"productname" binding:"required"`
}
