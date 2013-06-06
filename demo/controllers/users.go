package rippledemo

import (
	"../../ripple"
	"../models"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// ==========================================
// Define some basic controller for testing
// ==========================================

type UserController struct {
	userCollection rippledemo.UserCollection
	friends        []rippledemo.FriendshipModel
}

func NewUserController() *UserController {
	output := new(UserController)

	// Build some data for testing. In a real application, this would probably
	// come from a database. Also rather than being built here, it would
	// be accessed directly from each individual controller action.

	output.userCollection.Users = make(map[int]rippledemo.UserModel)
	output.userCollection.Add(rippledemo.UserModel{0, "John"})
	output.userCollection.Add(rippledemo.UserModel{0, "Paul"})
	output.userCollection.Add(rippledemo.UserModel{0, "Ringo"})
	output.userCollection.Add(rippledemo.UserModel{0, "George"})

	output.friends = append(output.friends, rippledemo.FriendshipModel{1, 2})
	output.friends = append(output.friends, rippledemo.FriendshipModel{1, 3})
	output.friends = append(output.friends, rippledemo.FriendshipModel{2, 4})
	output.friends = append(output.friends, rippledemo.FriendshipModel{3, 4})

	return output
}

func (this *UserController) Get(ctx *ripple.Context) {
	userId, _ := strconv.Atoi(ctx.Params["id"])
	if userId > 0 {
		ctx.Response.Body = this.userCollection.Get(userId)
	} else {
		ctx.Response.Body = this.userCollection.GetAll()
	}
}

func (this *UserController) Post(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	var user rippledemo.UserModel
	json.Unmarshal(body, &user)
	ctx.Response.Body = this.userCollection.Add(user)
}

func (this *UserController) Put(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	userId, _ := strconv.Atoi(ctx.Params["id"])
	var user rippledemo.UserModel
	json.Unmarshal(body, &user)
	ctx.Response.Body = this.userCollection.Set(userId, user)
}

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

func (this *UserController) PostFriends(ctx *ripple.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	userId, _ := strconv.Atoi(ctx.Params["id"])
	friendId, _ := strconv.Atoi(string(body))
	this.friends = append(this.friends, rippledemo.FriendshipModel{userId, friendId})
}
