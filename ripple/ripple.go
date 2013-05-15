package ripple

import (
	"fmt"
	"strings"
	"reflect"
	"net/http"
)

type Dependencies struct {
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

func (this *Application) RegisterController(name string, controller interface{}) {
	// TODO: check that controller satisfies interface
	this.controllers[name] = controller
}

func (this *Application) AddRoute(route Route) {
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

func (this *Application) Dispatch(request *http.Request) {
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
		if !exists {
			// warn?
			continue
		}
		
		methodName := strings.Title(strings.ToLower(request.Method)) + strings.Title(actionName)	
		controllerVal := reflect.ValueOf(controller)
				
		controllerMethod := controllerVal.MethodByName(methodName)
		if !controllerMethod.IsValid() {
			// warn?
			continue
		}
		
		_ = fmt.Println
		
		dep := controllerVal.Elem().Field(0).Addr().Interface().(*Dependencies)
		dep.Request = request
		dep.Params = params
		controllerMethod.Call(nil)
		return
	}
	
	fmt.Println("No match: " + request.URL.Path)
}