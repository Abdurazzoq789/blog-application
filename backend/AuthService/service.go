package main

import (
	"context"
	"errors"
	"github.com/Abdurazzoq789/blog_application/global"
	proto "github.com/Abdurazzoq789/blog_application/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net"
	"regexp"
	"time"
)

type authServer struct{}

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	login, password := in.GetLogin(), in.GetPassword()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var user global.User

	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)

	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("Wrong Login Credentials provided ")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("Wrong Login Credentials provided ")
	}

	return &proto.AuthResponse{Token: user.GetToken()} ,nil
}


func (server authServer) Signup(_ context.Context, in *proto.SignupRequest) (*proto.AuthResponse, error)  {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(username) < 4 || len(username) > 20 || len(email) < 7 || len(email) > 35 || len(password) < 8 || len(password) > 128 || !match {
		return &proto.AuthResponse{}, errors.New("Validation failed")
	}

	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Println("Error returned from username used: ", err.Error())
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("Username is used")
	}

	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: email})
	if err != nil {
		log.Println("Error returned from email used: ", err.Error())
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		return &proto.AuthResponse{}, errors.New("Email is used")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := global.User{ID: primitive.NewObjectID(), Username: username, Email: email, Password: string(pw)}

	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = global.DB.Collection("user").InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error inserting newUser: ", err.Error())
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}

	return &proto.AuthResponse{Token: newUser.GetToken()}, nil
}

func (authServer) UsernameUsed(_ context.Context, in *proto.UsernameUsedRequest) (*proto.UsedResponse, error)  {
	username :=in.GetUsername()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User

	global.DB.Collection("user").FindOne(ctx, bson.M{"username" : username}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error)  {
	email :=in.GetEmail()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User

	global.DB.Collection("user").FindOne(ctx, bson.M{"email" : email}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) AuthUser(_ context.Context, in *proto.AuthUserRequest) (*proto.AuthUserResponse, error)  {
	token := in.GetToken()
	user := global.UserFromToken(token)
	return &proto.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})

	listener, err := net.Listen("tcp", ":5000")

	if err != nil {
		log.Fatal("Error creating listener")
	}
	server.Serve(listener)

}
