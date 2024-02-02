package entity

// User -.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserListResponse struct {
	UserList []User `json:"user_list"`
}
