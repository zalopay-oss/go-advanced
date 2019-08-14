package main

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

type Nested struct {
    // validate định dạng email
    Email string `validate:"email"`
}
type T struct {
    // chỉ cho phép age = 10
    Age    int `validate:"eq=10"`
    Nested Nested
}

// validateEmail giúp xử lý các tag email
func validateEmail(input string) bool {
    if pass, _ := regexp.MatchString(
        `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, input,
    ); pass {
        return true
    }
    return false
}

// validate thực hiện công việc validate cho interface bất kỳ
// ở đây chỉ hiện thực cho kiểu T
func validate(v interface{}) (bool, string) {
    validateResult := true
    errmsg := "success"

    // xác định type và value của interface input
    vt := reflect.TypeOf(v)
    vv := reflect.ValueOf(v)

    // lần lượt duyệt trên mỗi field của struct
    for i := 0; i < vv.NumField(); i++ {
        // phân giải tag để áp dụng validate thích hợp
        fieldVal := vv.Field(i)
        tagContent := vt.Field(i).Tag.Get("validate")
        k := fieldVal.Kind()

        // điều kiện xét trên kiểu field của struct cần validate
        switch k {

        // trường hợp field là int
        case reflect.Int:
            // thực hiện validate cho tag eq=10
            val := fieldVal.Int()
            tagValStr := strings.Split(tagContent, "=")
            tagVal, _ := strconv.ParseInt(tagValStr[1], 10, 64)
            if val != tagVal {
                errmsg = "validate int failed, tag is: "+ strconv.FormatInt(
                    tagVal, 10,
                )
                validateResult = false
            }

        // trường hợp field là string
        case reflect.String:
            val := fieldVal.String()
            tagValStr := tagContent
            switch tagValStr {

            // nếu tag là email thì thực hiện validate tương ứng
            case "email":
                nestedResult := validateEmail(val)
                if nestedResult == false {
                    errmsg = "validate mail failed, field val is: "+ val
                    validateResult = false
                }
            }

        // nếu có struct lồng bên trong thì truyền
        // xuống đệ quy theo chiều sâu
        case reflect.Struct:
            valInter := fieldVal.Interface()
            nestedResult, msg := validate(valInter)
            if nestedResult == false {
                validateResult = false
                errmsg = msg
            }
        }
    }
    return validateResult, errmsg
}

func main() {
    // khởi tạo obj để test
    var a = T{Age: 10, Nested: Nested{Email: "abc@adfgom"}}

    validateResult, errmsg := validate(a)
    fmt.Println(validateResult, errmsg)
}

// kết quả:
// false validate mail failed, field val is: abc@adfgom