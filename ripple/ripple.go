package ripple

import (
	"log"
	"strings"
	"reflect"
	"net/http"
)

type Context struct {
	Params map[string] string;
	Request *http.Request;
}

type Application struct {
	controllers map[string] interface{}
	routes []Route
}

type Route struct {
	Pattern string
	Controller string
	Action string
}

func NewApplication() Application {
	var output Application
	output.controllers = make(map[string] interface{})
	return output
}

func (this *Application) checkRoute(route Route) {
	if route.Controller != "" {
		_, exists := this.controllers[route.Controller]
		if !exists {
			log.Panicf("\"%s\" controller does not exist.")
		}
	}
}

func (this *Application) RegisterController(name string, controller interface{}) {
	this.controllers[name] = controller
}

func (this *Application) AddRoute(route Route) {
	this.checkRoute(route)
	this.routes = append(this.routes, route)
}

func splitPath(path string) []string {
	var output []string
	if len(path) == 0 { return output }
	if path[0] == '/' { path = path[1:] }
	pathTokens := strings.Split(path, "/")
	for i := 0; i < len(pathTokens); i++ {
		e := pathTokens[i]
		if len(e) > 0 { output = append(output, e) }
	} 
	return output
}

func makeMethodName(requestMethod string, actionName string) string {
	return strings.Title(strings.ToLower(requestMethod)) + strings.Title(actionName)	
}

type MatchRequestResult struct {
	Success bool
	ControllerName string
	ActionName string
	ControllerValue reflect.Value
	ControllerMethod reflect.Value
	Params map[string] string
}

func (this *Application) matchRequest(request *http.Request) MatchRequestResult {
	var output MatchRequestResult
	output.Success = false
	
	path := request.URL.Path
	pathTokens := splitPath(path)
		
	for routeIndex := 0; routeIndex < len(this.routes); routeIndex++ {
		route := this.routes[routeIndex]
		patternTokens := splitPath(route.Pattern)
			
		if len(patternTokens) != len(pathTokens) { continue }
		
		var controller interface{}
		var exists bool
			
		controllerName := ""
		actionName := ""
		notMached := false
		params := make(map[string] string)
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
				notMached = true
				break
			}
		}
		
		if notMached { continue }
		
		if controllerName == "" {
			controllerName = route.Controller
		}
		
		if actionName == "" {
			actionName = route.Action
		}
		
		controller, exists = this.controllers[controllerName]
		if !exists { continue }
		
		methodName := makeMethodName(request.Method, actionName)
		controllerVal := reflect.ValueOf(controller)
				
		controllerMethod := controllerVal.MethodByName(methodName)
		if !controllerMethod.IsValid() { continue }
		
		output.Success = true
		output.ControllerName = controllerName
		output.ActionName = actionName
		output.ControllerValue = controllerVal
		output.ControllerMethod = controllerMethod
		output.Params = params
	}
	
	return output
}

func (this *Application) Dispatch(request *http.Request) {
	r := this.matchRequest(request)
	if !r.Success { return }
	
	ctx := new(Context)
	var args []reflect.Value
	ctx.Request = request
	ctx.Params = r.Params
	args = append(args, reflect.ValueOf(ctx))
	
	r.ControllerMethod.Call(args)
}