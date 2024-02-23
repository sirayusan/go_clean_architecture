package entity

type MessageResponse struct {
	UserID uint32     `gorm:"column:user_id"`
	List   []Messages `json:"list"`
}

// Message -.
type Message struct {
	LastName  string `gorm:"column:last_name"`
	FirstName string `gorm:"column:first_name"`
}

// Messages -.
type Messages struct {
	List []Messages `json:"list"`
}
