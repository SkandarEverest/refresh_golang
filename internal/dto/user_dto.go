package dto

type AuthRequest struct {
	Password string `json:"password" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email,min=1,max=255"`
	Action   string `json:"action" validate:"required,authaction"`
}

type AuthResponse struct {
	AccessToken string `json:"token"`
}

type UserResponse struct {
	Email           string `json:"email"`
	Username        string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}

type UpdateUserRequest struct {
	Email           *string `json:"email" validate:"omitempty,email,min=1,max=255"`
	Username        *string `json:"name"`
	UserImageUri    *string `json:"userImageUri"`
	CompanyName     *string `json:"companyName"`
	CompanyImageUri *string `json:"companyImageUri"`
}
