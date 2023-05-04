package route

import (
	"github.com/opensaucerer/barf"
	"github.com/opensaucerer/barf/app/controller"
)

func RegisterHomeRoutes() {

	barf.Get("/", controller.Home)
}
