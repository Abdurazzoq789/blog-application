package global

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/dgrijalva/jwt-go"
)

var NilUser User

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

func (u User) GetToken() string {
	byteSlc, _ := json.Marshal(u)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data" : string(byteSlc),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
