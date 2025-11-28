package ws

type Message struct {
	ID         string
	SenderID   string
	ReceiverID string // for direct user messaging
	RoomID     string
	Content    string
	Timestamp  string // TODO: change to time.Time type
}
