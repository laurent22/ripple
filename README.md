# Ripple, a REST API framework for Go #

Ripple is a simple, yet flexible, REST API framework for the Go language (golang).

Since building REST APIs often involves a lot of boiler plate code (building routes, handling GET, POST, etc. for each route), the framework attempts to simplify this by making assumptions about the way REST APIs are usually structured. The framework is flexible though, so if the defaults get in the way, they can be changed.

# Installation #

To install the library, clone this git directory:

    git clone https://github.com/laurent22/ripple.git

and import it in your main Go file using `import ./ripple`

# Demo #

A simple todo demo is available there: [https://github.com/laurent22/ripple-demo](https://github.com/laurent22/ripple-demo)

# Usage #

Have a look at [demo/demo.go](demo/demo.go) and [demo/demo_client.go](demo/demo_client.go) for a simple example of a REST API application. 

In general, the steps to build a REST application are as follow:

``` go
package main

import (
	"./ripple"
	"./controllers"
	"net/http"
)

func main() {	
	// Build the REST application
	
	app := ripple.NewApplication()
	
	// Create a controller and register it. Any number of controllers
	// can be registered that way.
	
	userController := rippledemo.NewUserController()
	app.RegisterController("users", userController)
	
	// Setup the routes. The special patterns `_controller` will automatically match
	// an existing controller, as defined above. Likewise, `_action` will match any 
	// existing action.
	
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	
	// Start the server
	
	http.ListenAndServe(":8080", app)
}
```

## Application ##

A Ripple application implements the [net.http Handler interface](http://golang.org/pkg/net/http/#Handler) and thus can be used with any Go server function that accepts this interface, including `ListenAndServe` or `ListenAndServeTLS`. To build a new application, call:

``` go
app := ripple.NewApplication()
```

Then, after having setup the controllers and routes (see below), call:

``` go
http.ListenAndServe(":8080", app)
```
    
This will create a REST API on `http://localhost:8080`

Ripple can also be used along other HTTP servers using `http.HandleFunc()` and accessing the `ServeHTTP()` function directly. For instance, to serve an HTML5/JS app under "/app" and the REST API under "/api", the following could be done:

``` go
// Handle the HTML5 app
http.HandleFunc("/app/", someHandleFunc)

// Handle the REST API	
app.SetBaseUrl("/api/")
http.HandleFunc("/api/", app.ServeHTTP) // Use the ServeHTTP function directly

// Start the server
http.ListenAndServe(":8080", nil)
```

Then the REST API will be at `http://localhost/api/:8080`, while the web application will be at `http://localhost/app/:8080`.

## Controllers ##

A Ripple controller is a `struct` with functions that handle the GET, POST, PUT, etc. HTTP methods (custom HTTP methods are also supported). The mapping between URLs and controller functions is done via routes (see below). Each function must start with the method name, followed by the (optional) action name. Each function receives a `ripple.Context` object that provides access to the full HTTP request, as well as the optional parameters. It also allows responding to the request. The code below shows a very simple controller that handles a GET method:

``` go
type UserController struct {}

func (this *UserController) Get(ctx *ripple.Context) {
	// Get the user ID:
	userId, _ := strconv.Atoi(ctx.Params["id"])
	if userId > 0 {
		// If a user ID is provided, we return the user with this ID.
		ctx.Response.Body = this.userCollection.Get(userId)
	} else {
		// If no user ID is provided, we return all the users.
		ctx.Response.Body = this.userCollection.GetAll()
	}
}
```

In the above code, `ctx.Params["id"]` is used to retrieve the user ID, the response is provided by setting `ctx.Response.Body`. The body will automatically be serialized to JSON.

To handle the POST method, you would write something like this:

``` go
func (this *UserController) Post(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	var user rippledemo.UserModel
	json.Unmarshal(body, &user)
	ctx.Response.Body = this.userCollection.Add(user)
}
```

Finally, more complex actions can be created. For example, this kind of function can be created to handle a REST URL such as `/users/123/friends`:

``` go
func (this *UserController) GetFriends(ctx *ripple.Context) {
	userId, _ := strconv.Atoi(ctx.Params["id"])
	var output []rippledemo.UserModel
	for _, d := range this.friends {
		if d.UserId1 == userId {
			output = append(output, this.userCollection.Get(d.UserId2))
		} else if d.UserId2 == userId {
			output = append(output, this.userCollection.Get(d.UserId1))
		}
	} 
	ctx.Response.Body = output
}
```

## Routes ##

The routes map a given URL to a given controller / action. Before being used in a route, the controllers must first be registered:

``` go
// Create some example controllers:

userController = new(UserController)
imageController = new(ImageController)

// And register them:

app.RegisterController("users", userController)
app.RegisterController("images", imageController)
```

Then the routes can be created:

``` go
app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
app.AddRoute(ripple.Route{ Pattern: ":_controller" })
```

Parameters can be defined by prefixing them with `:`; they are then accessible from the context object via `ctx.Params["id"]`.

Route patterns also accept two special parameters:

* `_controller`: Match any registered controller.
* `_action`: Match any existing controller action.

For example, the routes above would match URLs such as "users/123", "images/7", "images/456/metadata", etc. You do not need to specify the supported HTTP methods - whether a method is supported or not is implied from the controller functions. For instance, if the controller has a GetMetadata method, then `GET images/456/metadata` is automatically supported. Likewise, if it does *not* have a `DeleteMetadata` method, `DELETE images/456/metadata` will *not* be supported.

Routing can be as flexible as needed. If the automatic mapping of `_controller` and `_action` doesn't do the job, it is possible to explicitly specify the controller and action. For example:

``` go
app.AddRoute(Route{ Pattern: "some/very/custom/url", Controller: "users", Action: "example" })
```
	
With the above route, doing `GET some/very/custom/url` would call `UserController::GetExample`

## Models? ##

Ripple does not have built-in support for models since data storage can vary a lot from one application to another. For an example on how to connect a controller to a model, see [demo/controllers/users.go](demo/controllers/users.go) and [demo/models/user.go](demo/models/user.go). Usually, you would inject a database connection or other data source into the controller then use that from the various actions.

# Testing ##

Ripple is built with testability in mind. The whole framework is fully unit tested, and applications built with the framework can also be easily unit tested. Each controller method takes a `ripple.Context` object as parameter, which can be mocked for unit testing. The framework also exposes the `Application::Dispatch` method, which can be used to test the response for a given HTTP request.

# Ripple API reference ##

See the [Ripple GoDoc reference](http://godoc.org/github.com/laurent22/ripple) for more information.

# License #

The MIT License (MIT)

Copyright (c) 2013-2015 Laurent Cozic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
