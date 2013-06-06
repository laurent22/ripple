package rippledemo

// ======================================
// Define some basic model for testing
// ======================================

type UserModel struct {
	Id   int
	Name string
}

type FriendshipModel struct {
	UserId1 int
	UserId2 int
}

type UserCollection struct {
	Users      map[int]UserModel
	nextUserId int
}

func (this *UserCollection) Add(user UserModel) UserModel {
	this.nextUserId++
	user.Id = this.nextUserId
	this.Users[this.nextUserId] = user
	return user
}

func (this *UserCollection) Get(id int) UserModel {
	return this.Users[id]
}

func (this *UserCollection) Set(id int, user UserModel) UserModel {
	user.Id = id
	this.Users[id] = user
	return user
}

func (this *UserCollection) GetAll() []UserModel {
	var output []UserModel
	for _, d := range this.Users {
		output = append(output, d)
	}
	return output
}
