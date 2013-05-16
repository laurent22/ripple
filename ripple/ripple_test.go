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

type ControllerTesters struct {}
func (this *ControllerTesters) Get() {}
func (this *ControllerTesters) Post() {}
func (this *ControllerTesters) Patch() {}
func (this *ControllerTesters) GetFriends() {}

func TestMatchRequest(t *testing.T) {
	type MatchRequestTest struct {
		method string
		url string
		controller string
		action string
	}
	app := NewApplication()
	app.RegisterController("testers", &ControllerTesters{})
	app.AddRoute(Route{ Pattern: ":_controller/:id/:_action" })
	var reader io.Reader
	request, _ := http.NewRequest("GET", "http://localhost:8080/testers/123/friends", reader)
	result := app.matchRequest(request)
	
	
	t.Log("test: " + result.ControllerName + " " + result.ActionName)
}