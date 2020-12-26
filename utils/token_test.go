package utils

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	token, err := GenerateToken("sdfs")
	if err!=nil{
		fmt.Print(err.Error())
	}else {
		fmt.Println("token is : ",token)
	}
	uid,err := ValidToken(token)
	fmt.Println("uid: ",uid)
	if err!=nil{
		fmt.Println(err.Error())
	}else {
		fmt.Println("success")
	}

}



func TestVerifyToken(t *testing.T){
	uid, err := ValidToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzQzMDc4NjgsImlhdCI6MTU3NDMwNDI2OCwidWlkIjoiNWRkNWYyN2VlY2E0NjIyYzQ1ZDFjYzkzIn0.L5ADhV2cSxAELqN6IkQPiXbBSv6dHNj9JAmYAOqmrU4")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(uid)
}
