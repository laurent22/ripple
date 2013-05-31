package main

import (
	"./ripple"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ======================================
// Define some basic model for testing
// ======================================

type UserModel struct {
	Id int
	Name string	
}

type FriendshipModel struct {
	UserId1 int
	UserId2 int	
}

type UserCollection struct {
	users map[int]UserModel
	nextUserId int
}

func (this *UserCollection) Add(user UserModel) UserModel {
	this.nextUserId++
	user.Id = this.nextUserId
	this.users[this.nextUserId] = user
	return user
}

func (this *UserCollection) Get(id int) UserModel {
	return this.users[id]
}

func (this *UserCollection) Set(id int, user UserModel) UserModel {
	user.Id = id
	this.users[id] = user
	return user
}

func (this *UserCollection) GetAll() []UserModel {
	var output []UserModel
	for _, d := range this.users {
		output = append(output, d)
	}
	return output
}

var userCollection UserCollection
var friends []FriendshipModel

// ======================================
// Define controller
// ======================================

type UserController struct {}

func (this *UserController) Get(ctx *ripple.Context) {
	userId, _ := strconv.Atoi(ctx.Params["id"])
	if userId > 0 {
		ctx.Response.Body = userCollection.Get(userId)
	} else {
		ctx.Response.Body = userCollection.GetAll()
	}
}

func (this *UserController) Post(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	var user UserModel
	json.Unmarshal(body, &user)
	ctx.Response.Body = userCollection.Add(user)
}

func (this *UserController) Put(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	userId, _ := strconv.Atoi(ctx.Params["id"])
	var user UserModel
	json.Unmarshal(body, &user)
	ctx.Response.Body = userCollection.Set(userId, user)
}

func (this *UserController) GetFriends(ctx *ripple.Context) {
	userId, _ := strconv.Atoi(ctx.Params["id"])
	var output []UserModel
	for _, d := range friends {
		if d.UserId1 == userId {
			output = append(output, userCollection.Get(d.UserId2))
		} else if d.UserId2 == userId {
			output = append(output, userCollection.Get(d.UserId1))
		}
	} 
	ctx.Response.Body = output
}

func (this *UserController) PostFriends(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	userId, _ := strconv.Atoi(ctx.Params["id"])
	friendId, _ := strconv.Atoi(string(body))
	friends = append(friends, FriendshipModel{userId, friendId})
}

func main() {
	// Setup test models
	userCollection.users = make(map[int]UserModel)
	userCollection.Add(UserModel{ 0, "John" })
	userCollection.Add(UserModel{ 0, "Paul" })
	userCollection.Add(UserModel{ 0, "Ringo" })
	userCollection.Add(UserModel{ 0, "George" })
	
	friends = append(friends, FriendshipModel{1, 2})
	friends = append(friends, FriendshipModel{1, 3})
	friends = append(friends, FriendshipModel{2, 4})
	friends = append(friends, FriendshipModel{3, 4})
	
	// Build and run the REST application
	app := ripple.NewApplication()
	app.RegisterController("users", &UserController{})
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/:_action" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller/:id/" })
	app.AddRoute(ripple.Route{ Pattern: ":_controller" })
	http.ListenAndServe(":8080", app)
}