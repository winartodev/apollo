package enums

type ApplicationServiceEnum string

const (
	TestInternalServices1  ApplicationServiceEnum = "test-internal-services-1"
	TestInternalServices2  ApplicationServiceEnum = "test-internal-services-2"
	TestInternalServices3  ApplicationServiceEnum = "test-internal-services-3"
	TestProtectedServices3 ApplicationServiceEnum = "test-protected-services-3"
)

var applicationServiceEnumMap = map[ApplicationServiceEnum]string{
	TestInternalServices1:  "test-internal-services-1",
	TestInternalServices2:  "test-internal-services-2",
	TestInternalServices3:  "test-internal-services-3",
	TestProtectedServices3: "test-protected-services-3",
}

func (i ApplicationServiceEnum) IsValid() bool {
	for k, _ := range applicationServiceEnumMap {
		if i == k {
			return true
		}
	}

	return false
}

func (i ApplicationServiceEnum) String() string {
	for k, v := range applicationServiceEnumMap {
		if i == k {
			return v
		}
	}

	return ""
}

func (i ApplicationServiceEnum) ToSlug() string {
	for k, _ := range applicationServiceEnumMap {
		if i == k {
			return string(k)
		}
	}

	return ""
}
