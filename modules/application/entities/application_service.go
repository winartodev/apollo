package entities

type ApplicationService struct {
	Application `json:"application"`
	Services    []Service `json:"services"`
}
