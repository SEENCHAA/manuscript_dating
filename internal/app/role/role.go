package role

import "fmt"

// Role - тип для ролей пользователя
type Role int

const (
	Buyer   Role = iota // 0 - Обычный пользователь
	Manager             // 1 - Модератор
	Admin               // 2 - Администратор
)

// String для удобного вывода роли в JWT
func (r Role) String() string {
	switch r {
	case Buyer:
		return "Buyer"
	case Manager:
		return "Manager"
	case Admin:
		return "Admin"
	default:
		return fmt.Sprintf("Role(%d)", r)
	}
}
