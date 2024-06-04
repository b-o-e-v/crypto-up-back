package server

type User struct {
	Id          string `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name,omitempty"`
	ImageUrl    string `json:"image_url"`
	Password    string `json:"password,omitempty"`
	Token       string `json:"token,omitempty"`
}
