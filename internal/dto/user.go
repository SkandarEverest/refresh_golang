package dto

type UserRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email,min=1,max=255"`
}

type UserResponse struct {
	AccessToken string `json:"token"`
}
