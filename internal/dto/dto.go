package dto

type CreateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type UpdateProductInput struct {
	ID    string  `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GetJWTInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
