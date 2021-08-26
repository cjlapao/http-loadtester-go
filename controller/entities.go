package controller

// LoginRequest entity
type LoginRequest struct {
	Username string `json:"username" bson:"username" yaml:"username"`
	Password string `json:"password" bson:"password" yaml:"username"`
}

// LoginErrorResponse entity
type LoginErrorResponse struct {
	Code    string `json:"code" yaml:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// LoginResponse entity
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expiring    string `json:"expiring"`
}

// User entity
type User struct {
	ID        string `json:"id" bson:"_id"`
	Email     string `json:"email" bson:"email"`
	Username  string `json:"username" bson:"username"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Password  string `json:"password" bson:"password"`
}
