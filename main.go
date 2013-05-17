package main

import (
	"fmt"
	//"io"
	"net/http"
	"./ripple"
)

type ControllerUsers struct {

}

type ControllerSessions struct {

}

type Essai struct {
	One string
	Two int	
}

func (this *ControllerUsers) Get(ctx *ripple.Context) {
	// var test Essai
	// test.One = "abcd";
	// test.Two = 123
	// ctx.Response.Body = test
	ctx.Response.Body = false
}

func (this *ControllerUsers) Post(ctx *ripple.Context) {
	
}

func (this *ControllerUsers) Patch(ctx *ripple.Context) {
	
}

func (this *ControllerUsers) GetFriends(ctx *ripple.Context) {
	fmt.Println("get friends")
	fmt.Println(ctx)
}

func (this *ControllerSessions) GetNew(ctx *ripple.Context) {
	fmt.Println("Get new session")
}

func main() {
	_ = fmt.Println
	
	app := ripple.NewApplication()
	app.RegisterController("users", &ControllerUsers{})
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	http.ListenAndServe(":8080", app)
	
	
	// _ = fmt.Println
	
	// var reader io.Reader
	// request, _ := http.NewRequest("GET", "http://localhost:8080/users/123/friends", reader)
	
	// app := ripple.NewApplication()
	// app.RegisterController("users", &ControllerUsers{})
	// app.RegisterController("sessions", &ControllerSessions{})
	// app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	// app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	// app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	// app.AddRoute(ripple.Route{ Pattern: "sessions/:_action", Controller: "sessions" })
	// app.Dispatch(request)
	
	// request, _ = http.NewRequest("GET", "http://localhost:8080/sessions/new", reader)
	// app.Dispatch(request)
	
	// test := ripple.NewResponse()
	
	// body := Essai{ One: "something", Two: 123456 }
	// test.Body = body
	
	// b, _ := json.Marshal(test.Body)
		
	// fmt.Println(string(b))
	
	
	
	// request, _ = http.NewRequest("GET", "http://localhost:8080/users/123", reader)
	// app.Dispatch(request)
	
	//request, _ = http.NewRequest("GET", "http://localhost:8080/users", reader)
	//app.Dispatch(request)
	
	//http.HandleFunc("/", app.RequestHandler())
	//http.ListenAndServe(":8080", app)
}