package auth

type AccoutOutDTO struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	ProfileName    string `json:"profileName"`
	ProfileSurname string `json:"profileSurname"`
}
