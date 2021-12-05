package main

import (
	"context"
	"github.com/Abdurazzoq789/blog_application/global"
	"github.com/Abdurazzoq789/blog_application/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "test@gmail.com", Username: "Carl", Password: string(pw)})
	server := authServer{}
	_, err := server.Login(context.Background(), &proto.LoginRequest{Login: "test@gmail.com", Password: "example"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "something", Password: "something"})
	if err == nil {
		t.Error("2: error was nil")
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "Carl", Password: "example"})
	if err != nil {
		t.Error("3: An error was returned: ", err.Error())
	}
}

func TestAuthServer_UsernameUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Carl"})
	server := authServer{}
	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "carlo"})
	if err != nil {
		t.Error("An error was returned: ", err.Error())
	}
	if res.GetUsed() {
		t.Error("1: Wrong result")
	}
	res, err = 	server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Carl"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}

	if !res.GetUsed() {
		t.Error("2: Wrong result")
	}
}

func TestAuthServer_EmailUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "carl@test.com"})
	server := authServer{}
	res, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "carlo@test.com"})
	if err != nil {
		t.Error("An error was returned: ", err.Error())
	}
	if res.GetUsed() {
		t.Error("1: Wrong result")
	}
	res, err = 	server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "carl@test.com"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}

	if !res.GetUsed() {
		t.Error("2: Wrong result")
	}
}

func TestAuthServer_Signup(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "carl", Email: "carl@gmail.com"})
	server := authServer{}
	_, err := server.Signup(context.Background(), &proto.SignupRequest{Username: "carl", Email: "example@gmail.com", Password: "examplestring"})
	if err.Error() != "Username is used" {
		t.Error("1: No or the wrong Error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "carl@gmail.com", Password: "examplestring"})
	if err.Error() != "Email is used" {
		t.Error("2: No or the wrong Error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "example@gmail.com", Password: "examplestring"})
	if err != nil{
		t.Error("3: an error was returned")
	}

	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "example", Email: "example@gmail.com", Password: "exam"})
	if err.Error() != "Validation failed" {
		t.Error("3: No or the wrong Error was returned")
	}
}

