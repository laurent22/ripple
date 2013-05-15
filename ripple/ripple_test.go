package ripple

import (
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

func TestRegisterInvalidController1(t *testing.T) {
	type InvalidController struct {
	}
	defer func() {
		recover()
	}()
	checkControllerType(&InvalidController{})
	t.Error("Function did not panic.")
}

func TestRegisterInvalidController2(t *testing.T) {
	type InvalidController struct {
		Dep string
	}
	defer func() {
		recover()
	}()
	checkControllerType(&InvalidController{})
	t.Error("Function did not panic.")
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