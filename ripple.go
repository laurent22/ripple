package ripple

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var responseDefaultStatus = http.StatusOK

// A context holds information about the current request and response.
// An object of type Context is passed to methods of a controller.
type Context struct {
	// The parameters matched in the current URL. A parameter is defined
	// in a route by prefixing it with a ":". For example, ":id", ":name".
	Params map[string]string
	// The actual HTTP request.
	Request *http.Request
	// The response object.
	Response *Response
}

// Build a new context object.
func NewContext() *Context {
	return new(Context)
}

// A Ripple application. Use NewApplication() to build it.
type Application struct {
	controllers map[string]interface{}
	routes      []Route
	contentType string
}

// Build a new application object.
func NewApplication() *Application {
	output := new(Application)
	output.controllers = make(map[string]interface{})
	output.contentType = "application/json"
	return output
}

// A route maps a URL pattern to a controller and action.
// See README.md for more information on the format of the pattern.
type Route struct {
	Pattern    string
	Controller string
	Action     string
}

// Holds information about the HTTP response.
type Response struct {
	// HTTP code (200, 404, etc.). Use the constants of the http package.
	Status int
	// The response body. It will be serialized automatically
	// by the Ripple application before being sent to the client.
	Body interface{}
}

// Build a new response object.
func NewResponse() *Response {
	output := new(Response)
	output.Status = responseDefaultStatus
	output.Body = nil
	return output
}

// Helper struct used by `prepareServeHttpResponseData()`
type serveHttpResponseData struct {
	Status int
	Body   string
}

// Helper function to prepare the response writter data for `ServeHTTP()`
func (this *Application) prepareServeHttpResponseData(context *Context) serveHttpResponseData {
	var statusCode int
	var body string
	var err error
	if context == nil {
		statusCode = http.StatusNotFound
	} else {
		statusCode = context.Response.Status
	}
	if context != nil {
		body, err = this.serializeResponseBody(context.Response.Body)
		if err != nil {
			statusCode = http.StatusInternalServerError
		}
	}

	var output serveHttpResponseData
	output.Status = statusCode
	output.Body = body
	return output
}

// Serves an HTTP request - implementation of net.http.ServeHTTP
func (this *Application) ServeHTTP(writter http.ResponseWriter, request *http.Request) {
	context := this.Dispatch(request)
	r := this.prepareServeHttpResponseData(context)
	writter.Header().Set("Content-Type", this.contentType)
	writter.WriteHeader(r.Status)
	writter.Write([]byte(r.Body))
}

func (this *Application) serializeResponseBody(body interface{}) (string, error) {
	if body == nil {
		return "", nil
	}

	var output string
	var err error
	err = nil

	switch body.(type) {

	case string:

		output = body.(string)

	case int, int8, int16, int32, int64:

		output = strconv.Itoa(body.(int))

	case uint, uint8, uint16, uint32, uint64:

		output = strconv.FormatUint(body.(uint64), 10)

	case float32, float64:

		output = strconv.FormatFloat(body.(float64), 'f', -1, 64)

	case bool:

		if body.(bool) {
			output = "true"
		} else {
			output = "false"
		}

	default:

		if this.contentType == "application/json" {
			var b []byte
			b, err = json.Marshal(body)
			output = string(b)
		} else {
			log.Panicf("Unsupported content type: %s", this.contentType)
		}

	}

	return output, err
}

func (this *Application) checkRoute(route Route) {
	if route.Controller != "" {
		_, exists := this.controllers[route.Controller]
		if !exists {
			log.Panicf("\"%s\" controller does not exist.\n", route.Controller)
		}
	}
}

// Registers a controller. The name should be the same as in the URL path. For example
// if the URL is "users/1", the name should be "users". The controller itself can be
// any struct that implements HTTP method handlers. See README.md and the demo for more
// details on the structure of a controller.
func (this *Application) RegisterController(name string, controller interface{}) {
	this.controllers[name] = controller
}

// Add a route to the application.
func (this *Application) AddRoute(route Route) {
	this.checkRoute(route)
	this.routes = append(this.routes, route)
}

func splitPath(path string) []string {
	var output []string
	if len(path) == 0 {
		return output
	}
	if path[0] == '/' {
		path = path[1:]
	}
	pathTokens := strings.Split(path, "/")
	for i := 0; i < len(pathTokens); i++ {
		e := pathTokens[i]
		if len(e) > 0 {
			output = append(output, e)
		}
	}
	return output
}

func makeMethodName(requestMethod string, actionName string) string {
	return strings.Title(strings.ToLower(requestMethod)) + strings.Title(actionName)
}

// Provided for debugging/testing purposes only.
type MatchRequestResult struct {
	Success          bool
	ControllerName   string
	ActionName       string
	ControllerValue  reflect.Value
	ControllerMethod reflect.Value
	MatchedRoute     Route
	Params           map[string]string
}

func (this *Application) matchRequest(request *http.Request) MatchRequestResult {
	var output MatchRequestResult
	output.Success = false

	path := request.URL.Path
	pathTokens := splitPath(path)

	for routeIndex := 0; routeIndex < len(this.routes); routeIndex++ {
		route := this.routes[routeIndex]
		patternTokens := splitPath(route.Pattern)

		if len(patternTokens) != len(pathTokens) {
			continue
		}

		var controller interface{}
		var exists bool

		controllerName := ""
		actionName := ""
		notMatched := false
		params := make(map[string]string)
		for i := 0; i < len(patternTokens); i++ {
			patternToken := patternTokens[i]
			pathToken := pathTokens[i]
			if patternToken == ":_controller" {
				controllerName = pathToken
			} else if patternToken == ":_action" {
				actionName = pathToken
			} else if patternToken == pathToken {

			} else if patternToken[0] == ':' {
				params[patternToken[1:]] = pathToken
			} else {
				notMatched = true
				break
			}
		}

		if notMatched {
			continue
		}

		if controllerName == "" {
			controllerName = route.Controller
		}

		if actionName == "" {
			actionName = route.Action
		}

		controller, exists = this.controllers[controllerName]
		if !exists {
			continue
		}

		methodName := makeMethodName(request.Method, actionName)
		controllerVal := reflect.ValueOf(controller)

		controllerMethod := controllerVal.MethodByName(methodName)
		if !controllerMethod.IsValid() {
			continue
		}

		output.Success = true
		output.ControllerName = controllerName
		output.ActionName = actionName
		output.ControllerValue = controllerVal
		output.ControllerMethod = controllerMethod
		output.MatchedRoute = route
		output.Params = params
	}

	return output
}

// Provided for debugging/testing purposes only.
func (this *Application) Dispatch(request *http.Request) *Context {
	r := this.matchRequest(request)
	if !r.Success {
		log.Printf("No match for: %s %s\n", request.Method, request.URL)
		return nil
	}

	ctx := NewContext()
	ctx.Request = request
	ctx.Response = NewResponse()
	ctx.Params = r.Params
	var args []reflect.Value
	args = append(args, reflect.ValueOf(ctx))

	r.ControllerMethod.Call(args)
	return ctx
}
