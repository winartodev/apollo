package application

type ID int64

const (
	apollo ID = iota + 1
)

func (i ID) ToString() string {
	return [...]string{"Apollo"}[i-1]
}

func (i ID) ToInt64() int64 {
	return int64(i)
}

type Scope int64

const (
	Public Scope = iota + 1
	Internal
	Protected
)

func (s Scope) ToString() string {
	return [...]string{"Public", "Internal", "Protected"}[s-1]
}

func (s Scope) ToInt64() int64 {
	return int64(s)
}

type Access struct {
	ID    ID    `json:"id"`
	Scope Scope `json:"scope"`
}

var (
	ApolloPublic    = &Access{apollo, Public}
	ApolloInternal  = &Access{apollo, Internal}
	ApolloProtected = &Access{apollo, Protected}
)
