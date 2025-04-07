package enums

type ApplicationEnum string

const (
	ApolloPublic    ApplicationEnum = "apollo"
	ApolloInternal  ApplicationEnum = "apollo-internal"
	ApolloProtected ApplicationEnum = "apollo-protected"
)

var applicationEnumMap = map[ApplicationEnum]string{
	ApolloPublic:    "Apollo",
	ApolloInternal:  "Apollo Internal",
	ApolloProtected: "Apollo Protected",
}

func (i ApplicationEnum) IsValid() bool {
	for k, _ := range applicationEnumMap {
		if i == k {
			return true
		}
	}

	return false
}

func (i ApplicationEnum) String() string {
	for k, v := range applicationEnumMap {
		if i == k {
			return v
		}
	}

	return ""
}

func (i ApplicationEnum) ToSlug() string {
	for k, _ := range applicationEnumMap {
		if i == k {
			return string(k)
		}
	}

	return ""
}
