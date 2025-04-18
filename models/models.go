package models

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Guid string `json:"guid"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}
