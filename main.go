package main

import (
	"fmt"
	"io"
	"net/http"
	"./ripple"
)

type ControllerUsers struct {
	D ripple.Dependencies;
}

type ControllerSessions struct {
	D ripple.Dependencies;
}

func (this *ControllerUsers) Get() {
	fmt.Println("GET")
}

func (this *ControllerUsers) Post() {
	
}

func (this *ControllerUsers) Patch() {
	
}

func (this *ControllerUsers) GetFriends() {
	fmt.Println("get friends")
	fmt.Println(this.D)
}

func (this *ControllerSessions) GetNew() {
	fmt.Println("Get new session")
}

func handler(writter http.ResponseWriter, request *http.Request) {	
	fmt.Fprintf(writter, "Hi  %s!", request.URL.Path[1:])
}

func main() {
	var reader io.Reader
	request, _ := http.NewRequest("GET", "http://localhost:8080/users/123/friends", reader)
	
	app := ripple.NewApplication()
	app.RegisterController("users", &ControllerUsers{})
	app.RegisterController("sessions", &ControllerSessions{})
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	app.AddRoute(ripple.Route{ Pattern: "sessions/:_action", Controller: "sessions" })
	app.Dispatch(request)
	
	request, _ = http.NewRequest("GET", "http://localhost:8080/sessions/new", reader)
	app.Dispatch(request)
	
	// request, _ = http.NewRequest("GET", "http://localhost:8080/users/123", reader)
	// app.Dispatch(request)
	
	//request, _ = http.NewRequest("GET", "http://localhost:8080/users", reader)
	//app.Dispatch(request)
	
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)
}