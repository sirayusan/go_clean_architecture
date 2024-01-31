// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// User -.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserListResponse struct {
	UserList []User `json:"history"`
}
