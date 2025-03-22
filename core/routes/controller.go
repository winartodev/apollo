package routes

type ControllerDependency struct {
	Repository *Repository
}

type Controller struct {
}

func NewController(dependency ControllerDependency) *Controller {
	_ = dependency.Repository

	return &Controller{}
}
