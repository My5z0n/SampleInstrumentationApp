package model

type GetUserInfoModelInput struct {
	User string `uri:"user" binding:"required"`
}
