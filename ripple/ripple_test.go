package ripple

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSplitPath(t *testing.T) {
	type SplitPathTest struct {
		input string
		expected []string
	}

	var splitPathTests = []SplitPathTest{
		{ "/users/123",  []string{ "users", "123" } },
		{ "users/123",   []string{ "users", "123" } },
		{ "users/123/",  []string{ "users", "123" } },
		{ "users//123/", []string{ "users", "123" } },
		{ "users", []string{ "users" } },
		{ "", []string{} },
	}
	
	for _, d := range splitPathTests {
		output := strings.Join(splitPath(d.input), ",")
		expected := strings.Join(d.expected, ",")
		if output != expected {
			t.Errorf("Expected %s; Got %s", expected, output)
		}
	}	
}

func TestMakeMethodName(t *testing.T) {
	type MakeMethodNameTest struct {
		method string
		action string
		expected string
	}
	var makeMethodNameTests = []MakeMethodNameTest{
		{ "GET", "user", "GetUser" },
		{ "POST", "", "Post" },
		{ "DELETE", "image", "DeleteImage" },
	}
	for _, d := range makeMethodNameTests {
		output := makeMethodName(d.method, d.action)
		if output != d.expected {
			t.Errorf("Expected %s; Got %s", d.expected, output)
		}
	}	
}

func TestAddRoutePanic(t *testing.T) {
	app := NewApplication()
	defer func() { recover() }()
	app.AddRoute(Route{ Pattern: ":_controller", Controller: "nope" })
	t.Error("Added invalid controller but AddRoute did not panic.")
}

type ControllerTesters struct {}
func (this *ControllerTesters) Get(ctx *Context) {}
func (this *ControllerTesters) Post(ctx *Context) {}
func (this *ControllerTesters) Patch(ctx *Context) {}
func (this *ControllerTesters) GetTasks(ctx *Context) {}

type ControllerTesters3 struct {}
func (this *ControllerTesters3) GetNew(ctx *Context) {}
func (this *ControllerTesters3) PostCustom(ctx *Context) {}

func TestMatchRequest(t *testing.T) {		
	type MatchRequestTest struct {
		method string
		url string
		success bool
		controller string
		action string
		params map[string]string
	}
	var matchRequestTests = []MatchRequestTest{
		{ "GET", "/testers", true, "testers", "", map[string]string{} },
		{ "GET", "/testers/123/tasks", true, "testers", "tasks", map[string]string{ "id": "123" } },
		{ "POST", "/testers", true,  "testers", "", map[string]string{} },
		{ "GET", "/controllernotthere", false, "", "", map[string]string{} },
		{ "GET", "/testers/123/oops", false, "", "", map[string]string{} },
		{ "DELETE", "/testers/123/tasks", false, "", "", map[string]string{} },
		{ "GET", "/testers3/new", true, "testers3", "new", map[string]string{} },
		{ "GET", "/testers3/nothere", false, "", "", map[string]string{} },
		{ "POST", "/testers3/custom/something", true, "testers3", "custom", map[string]string{} },
		{ "POST", "/testers3/custom/123/456/789", true, "testers3", "custom", map[string]string{"one":"123","two":"456","three":"789"} },
	}	
	app := NewApplication()
	app.RegisterController("testers", &ControllerTesters{})
	app.RegisterController("testers3", &ControllerTesters3{})
	app.AddRoute(Route{ Pattern: "testers3/:_action", Controller: "testers3" })
	app.AddRoute(Route{ Pattern: "testers3/:_action/:one/:two/:three", Controller: "testers3" })
	app.AddRoute(Route{ Pattern: "testers3/custom/something", Controller: "testers3", Action: "custom" })
	app.AddRoute(Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(Route{ Pattern: ":_controller" })
	var reader io.Reader
	for _, d := range matchRequestTests {
		request, _ := http.NewRequest(d.method, d.url, reader)
		result := app.matchRequest(request)
		if result.Success != d.success {
			t.Errorf("%s %s: Expected success = '%b', got '%b'", d.method, d.url, d.success, result.Success)
		}
		if result.ControllerName != d.controller {
			t.Errorf("%s %s: Expected controller '%s', got '%s'", d.method, d.url, d.controller, result.ControllerName)
		}
		if result.ActionName != d.action {
			t.Errorf("%s %s: Expected action '%s', got '%s'", d.method, d.url, d.action, result.ActionName)
		}
		paramOk := true
		if len(result.Params) != len(d.params) {
			paramOk = false
		} else {
			for key, value := range result.Params {
				eValue, eExists := d.params[key]
				if !eExists || eValue != value {
					paramOk = false
					break
				}
			}
		}
		if !paramOk {
			t.Errorf("%s %s: Param mismatch: Expected '%s', got '%s'", d.method, d.url, d.params, result.Params)
		}
	}
}

type ControllerTesters2 struct {
	GetContext *Context
}
func (this *ControllerTesters2) Get(ctx *Context) {
	this.GetContext = ctx	
}

func TestDispatch(t *testing.T) {
	var controller ControllerTesters2
	
	app := NewApplication()
	app.RegisterController("testers2", &controller)
	app.AddRoute(Route{ Pattern: ":_controller/:id" })
	var reader io.Reader
	request, _ := http.NewRequest("GET", "/testers2/abcd", reader)
	app.Dispatch(request)
	
	paramId, ok := controller.GetContext.Params["id"]
	if !ok || paramId != "abcd" {
		t.Errorf("Controller action did not get correct parameter: %t/'%s' (Expected: 'abcd')", ok, paramId)
	}
	if controller.GetContext.Request == nil {
		t.Errorf("Controller action got a nil request.")
	}
}