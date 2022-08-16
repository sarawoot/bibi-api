package handler

import "api/internal/model"

type UserSignupRequest struct {
	Email         string       `json:"email" binding:"required,email"`
	Password      string       `json:"password" binding:"required,min=4,max=32"`
	Gender        model.Gender `json:"gender"`
	Birthdate     model.Date   `json:"birthdate"`
	SkinTypeID    model.UUID   `json:"skin_type_id"`
	SkinProblemID model.UUID   `json:"skin_problem_id"`
}

type UserSignupResponse struct {
	AccessToken string `json:"access_token"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
}
