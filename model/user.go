package model

type User struct {
	ID           uint   `json:"id"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"displayname"`
	UserName     string `json:"username"`
	PasswordHash string `json:"passwordhash"`
}
