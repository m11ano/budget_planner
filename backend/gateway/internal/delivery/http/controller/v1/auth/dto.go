package auth

type AccoutOutDTO struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	ProfileName    string `json:"profileName"`
	ProfileSurname string `json:"profileSurname"`
}

type TokensOutDTO struct {
	RefreshJWT string `json:"refreshJWT"`
	AccessJWT  string `json:"accessJWT"`
}
