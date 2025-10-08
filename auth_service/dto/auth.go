package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EditRequest struct {
	ID       int64  `param:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteRequest struct {
	ID int64 `param:"id"`
}
