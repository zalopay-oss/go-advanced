package main

import (
	"gopkg.in/go-playground/validator.v9"
	"fmt"
	// "time"
)

type RegisterReq struct {
    // gt = 0 cho biết độ dài chuỗi phải > 0，gt: greater than
    Username       string   `json:"username" validate:"gt=0"`
    // như trên
    PasswordNew    string   `json:"password_new" validate:"gt=0"`
    // eqfield kiểm tra các trường bằng nhau
    PasswordRepeat string   `json:"password_repeat" validate:"eqfield=PasswordNew"`
    // kiểm tra định dạng email thích hợp
    Email          string   `json:"email" validate:"email"`
}


func validatefunc(req RegisterReq) error {
	err := validate.Struct(req)
    if err != nil {
        return err
    }
    return nil
}
var validate *validator.Validate

func main() {
	validate = validator.New()

	a := RegisterReq{
		Username: "Alex",
		PasswordNew: "z",
		PasswordRepeat: "z",
		Email: "zz@abccom",
	}

	err := validatefunc(a)
	fmt.Println(err)

}