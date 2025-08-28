package models

type LoginRequest struct {
	Username string `json:"username" example:"alice@example.com"`
	Password string `json:"password" example:"s3cr3t"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"  example:"<jwt>"`
	RefreshToken string `json:"refresh_token" example:"<jwt>"`
	ExpiresAt    int64  `json:"expires_at"    example:"1756069590"`
}

type RefreshRequest struct {
	UserID       string `json:"user_id"       example:"User_65b80522-50be-4012-9964-550369cdcff7"`
	RefreshToken string `json:"refresh_token" example:"<jwt>"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"<jwt>"`
}
