package enums

type ApplicationScope int

const (
	Public ApplicationScope = iota
	Protected
	Internal
)

func (s ApplicationScope) ToString() string {
	return [...]string{"public", "protected", "internal"}[s]
}

func (s ApplicationScope) ToInt() int {
	return int(s)
}
