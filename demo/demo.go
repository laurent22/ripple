package main

import (
	"../ripple"
	"./controllers"
	"net/http"
)

func main() {
	userController := rippledemo.NewUserController()
	
	// Build and run the REST application
	app := ripple.NewApplication()
	app.RegisterController("users", userController)
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	http.ListenAndServe(":8080", app)
}