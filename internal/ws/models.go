package ws

type DirectMessage struct {
	UserID uint
	Msg    []byte
}

type User struct {
	ID   uint
	Name string
}

// `gorm:"primaryKey"`

type MessagePayload struct {
	SenderID   uint
	ReceiverID *uint
	Content    string
}
