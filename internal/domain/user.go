package domain

// Структура пользователя для понимания какую команду надо выполнять.
type User struct {
	Command          string
	Money            int
	TargetMap        map[string]int
	DbChangeTargetId int
}
