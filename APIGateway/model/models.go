package model

type GetUserInfoModelInput struct {
	User string `uri:"user" binding:"required"`
}
type CreateOrderModel struct {
	ProductName string `json:"ProductName" binding:"required"`
	Coupon      string `json:"Coupon"`
}
