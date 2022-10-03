package model

type GetUserInfoModelInput struct {
	User string `uri:"user" binding:"required"`
}
type CreateOrderModel struct {
	ProductName string `json:"user" binding:"required"`
}
type ProductDetailsModel struct {
	ProductName string `uri:"productname" binding:"required"`
}
