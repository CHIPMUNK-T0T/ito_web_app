package entity

type Role string

const (
	Host  Role = "host"
	Guest Role = "guest"
)

type RoomStatus string

const (
	RoomStatusWaiting RoomStatus = "waiting"
	RoomStatusPlaying RoomStatus = "playing"
)
