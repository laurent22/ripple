package ripple

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSplitPath(t *testing.T) {
	type SplitPathTest struct {
		input    string
		expected []string
	}

	var splitPathTests = []SplitPathTest{
		{"/users/123", []string{"users", "123"}},
		{"users/123", []string{"users", "123"}},
		{"users/123/", []string{"users", "123"}},
		{"users//123/", []string{"users", "123"}},
		{"users", []string{"users"}},
		{"", []string{}},
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
		method   string
		action   string
		expected string
	}
	var makeMethodNameTests = []MakeMethodNameTest{
		{"GET", "user", "GetUser"},
		{"POST", "", "Post"},
		{"DELETE", "image", "DeleteImage"},
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
	app.AddRoute(Route{Pattern: ":_controller", Controller: "nope"})
	t.Error("Added invalid controller but AddRoute did not panic.")
}

type ControllerTesters struct{}

func (this *ControllerTesters) Get(ctx *Context)      {}
func (this *ControllerTesters) Post(ctx *Context)     {}
func (this *ControllerTesters) Patch(ctx *Context)    {}
func (this *ControllerTesters) GetTasks(ctx *Context) {}

type ControllerTesters3 struct{}

func (this *ControllerTesters3) GetNew(ctx *Context)     {}
func (this *ControllerTesters3) PostCustom(ctx *Context) {}

func TestMatchRequest(t *testing.T) {
	type MatchRequestTest struct {
		method     string
		url        string
		success    bool
		controller string
		action     string
		params     map[string]string
	}
	var matchRequestTests = []MatchRequestTest{
		{"GET", "/testers", true, "testers", "", map[string]string{}},
		{"GET", "/testers/123/tasks", true, "testers", "tasks", map[string]string{"id": "123"}},
		{"POST", "/testers", true, "testers", "", map[string]string{}},
		{"GET", "/controllernotthere", false, "", "", map[string]string{}},
		{"GET", "/testers/123/oops", false, "", "", map[string]string{}},
		{"DELETE", "/testers/123/tasks", false, "", "", map[string]string{}},
		{"GET", "/testers3/new", true, "testers3", "new", map[string]string{}},
		{"GET", "/testers3/nothere", false, "", "", map[string]string{}},
		{"POST", "/testers3/custom/something", true, "testers3", "custom", map[string]string{}},
		{"POST", "/testers3/custom/123/456/789", true, "testers3", "custom", map[string]string{"one": "123", "two": "456", "three": "789"}},
	}
	app := NewApplication()
	app.RegisterController("testers", &ControllerTesters{})
	app.RegisterController("testers3", &ControllerTesters3{})
	app.AddRoute(Route{Pattern: "testers3/:_action", Controller: "testers3"})
	app.AddRoute(Route{Pattern: "testers3/:_action/:one/:two/:three", Controller: "testers3"})
	app.AddRoute(Route{Pattern: "testers3/custom/something", Controller: "testers3", Action: "custom"})
	app.AddRoute(Route{Pattern: ":_controller/:id/:_action"})
	app.AddRoute(Route{Pattern: ":_controller/:id/"})
	app.AddRoute(Route{Pattern: ":_controller/"})
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

type ControllerTesters4 struct{}

func (this *ControllerTesters4) Get(ctx *Context)      {}
func (this *ControllerTesters4) GetOther(ctx *Context) {}

func TestHardCodedAction(t *testing.T) {
	var controller ControllerTesters4

	app := NewApplication()
	app.RegisterController("testers", &controller)
	app.AddRoute(Route{Pattern: ":_controller/:id"})
	app.AddRoute(Route{Pattern: ":_controller/:id/other", Action: "other"})
	var reader io.Reader
	request, _ := http.NewRequest("GET", "/testers/abcd/other", reader)
	r := app.matchRequest(request)
	if r.ControllerName != "testers" {
		t.Errorf("Expected %s, got %s", "testers", r.ControllerName)
	}
	if r.ActionName != "other" {
		t.Errorf("Expected %s, got %s", "other", r.ActionName)
	}
}

type ControllerTesters2 struct {
	GetContext *Context
}

func (this *ControllerTesters2) Get(ctx *Context) {
	this.GetContext = ctx
	ctx.Response.Status = 202
}
func (this *ControllerTesters2) GetOther(ctx *Context) {
	this.GetContext = ctx
}

func TestDispatch(t *testing.T) {
	var controller ControllerTesters2

	app := NewApplication()
	app.RegisterController("testers2", &controller)
	app.AddRoute(Route{Pattern: ":_controller/:id"})
	app.AddRoute(Route{Pattern: ":_controller/:id/other", Action: "other"})
	var reader io.Reader
	request, _ := http.NewRequest("GET", "/testers2/abcd", reader)
	context := app.Dispatch(request)

	paramId, ok := controller.GetContext.Params["id"]
	if !ok || paramId != "abcd" {
		t.Errorf("Controller action did not get correct parameter: %t/'%s' (Expected: 'abcd')", ok, paramId)
	}
	if controller.GetContext.Request == nil {
		t.Errorf("Controller action got a nil request.")
	}
	if controller.GetContext.Response.Status != 202 {
		t.Errorf("Controller response has not been modified. Got %d, expected %d", controller.GetContext.Response.Status, 202)
	}
	if context.Response.Status != 202 {
		t.Errorf("Controller response has not been modified. Got %d, expected %d", context.Response.Status, 202)
	}

	request, _ = http.NewRequest("GET", "/testers2/abcd/other", reader)
	context = app.Dispatch(request)
	if context.Response.Status != http.StatusOK {
		t.Errorf("Response status is not set to correct default. Expected %d, got %d", http.StatusOK, context.Response.Status)
	}
}

func TestSerializeResponseBody(t *testing.T) {
	type SerializeTest struct {
		input    interface{}
		expected string
		success  bool
	}
	type StructTest struct {
		Something string
	}
	var serializeTests = []SerializeTest{
		{"abcdef", "abcdef", true},
		{123456, "123456", true},
		{123.45000, "123.45", true},
		{nil, "", true},
		{-12, "-12", true},
		{true, "true", true},
		{false, "false", true},
		{StructTest{Something: "hello"}, "{\"Something\":\"hello\"}", true},
		{StructTest{Something: "hello"}, "{\"Something\":\"hello\"}", true},
		{map[string]StructTest{"one": {"123"}}, "{\"one\":{\"Something\":\"123\"}}", true},
	}

	app := NewApplication()

	for _, d := range serializeTests {
		s, err := app.serializeResponseBody(d.input)
		if err == nil && !d.success {
			t.Errorf("Serialization should have failed.")
		}
		if err != nil && d.success {
			t.Errorf("Serialization should have succeeded.")
		}
		if s != d.expected {
			t.Errorf("Bad serialization. Expected %s, Got %s", d.expected, s)
		}
	}
}

func TestContextIsFullyInitialized(t *testing.T) {
	ctx := NewContext()
	if ctx.Response == nil {
		t.Errorf("Context response not initialized.")
	}
	defer func() {
		if recover() != nil {
			t.Errorf("Context params not initialized.")
		}
	}() 
	ctx.Params["id"] = "123"
}

func TestPrepareServeHttpResponseData(t *testing.T) {
	app := NewApplication()
	type ResponseTest struct {
		Status         int
		Body           string
		ExpectedStatus int
		ExpectedBody   string
	}
	var responseTests = []ResponseTest{
		{200, "ok", 200, "ok"},
		{200, "", 200, ""},
		{202, "", 202, ""},
	}
	for _, d := range responseTests {
		c := NewContext()
		c.Response = NewResponse()
		c.Response.Status = d.Status
		c.Response.Body = d.Body
		r := app.prepareServeHttpResponseData(c)
		if r.Status != d.ExpectedStatus {
			t.Errorf("Expected %d, Got %d", d.ExpectedStatus, r.Status)
		}
		if r.Body != d.ExpectedBody {
			t.Errorf("Expected %d, Got %d", d.ExpectedBody, r.Body)
		}
	}
}

func TestBaseUrl(t *testing.T) {
	type MatchRequestTest struct {
		baseUrl    string
		url        string
		success    bool
		controller string
		action     string
	}
	var matchRequestTests = []MatchRequestTest{
		{"/", "/testers/123", true, "testers", ""},
		{"/", "/api/testers/123", false, "", ""},
		{"/api/", "/api/testers/123", true, "testers", ""},
		{"/api/", "/api/testers/123/tasks", true, "testers", "tasks"},
	}
	
	app := NewApplication()
	app.RegisterController("testers", &ControllerTesters{})
	app.AddRoute(Route{Pattern: ":_controller/:id"})
	app.AddRoute(Route{Pattern: ":_controller/:id/:_action"})
	
	for _, d := range matchRequestTests {
		var reader io.Reader
		app.SetBaseUrl(d.baseUrl)
		request, _ := http.NewRequest("GET", d.url, reader)
		result := app.matchRequest(request)
		if result.Success != d.success {
			t.Errorf("Expected %t, got %t", d.success, result.Success)
		}
		if result.ControllerName != d.controller {
			t.Errorf("Expected %s, got %s", d.controller, result.ControllerName)
		}
		if result.ActionName != d.action {
			t.Errorf("Expected %s, got %s", d.action, result.ActionName)
		}
	}
}