package application

import (
	"github.com/winartodev/apollo/core/errors"
	"net/http"
)

const (
	SRV001 = "SRV_001"
)

var (
	ServiceNameIsEmptyErr = errors.New(http.StatusBadRequest, SRV001, "Service name is empty")
)
