package dto

// LoginInput входные данные для входа
type LoginInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginOutput выходные данные после входа
type LoginOutput struct {
	User  UserOutput `json:"user"`
	Token string     `json:"token"`
}

// RegisterInput входные данные для регистрации
type RegisterInput struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}

// RegisterOutput выходные данные после регистрации
type RegisterOutput struct {
	Message string     `json:"message"`
	User    UserOutput `json:"user"`
}

// ValidateTokenInput входные данные для валидации токена
type ValidateTokenInput struct {
	Token string `json:"token" validate:"required"`
}

