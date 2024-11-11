package model

type User struct {
	UUID string `json:"id"`
	Role string `json:"role"`
	UserCredentials
	Address     string `json:"address"`
	DateOfBirth string `json:"dateOfBirth"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Phone       string `json:"phone"`
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Username struct {
	Username string `json:"username"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserUUID     string `json:"user_uuid"`
}
