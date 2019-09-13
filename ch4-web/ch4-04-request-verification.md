# 4.4 Kiểm tra tính hợp lệ của request

Một nguyên tắc quan trọng trong lập trình web là không được hoàn toàn tin những gì mà user gửi lên, luôn phải có các cơ chế xác thực, kiểm tra tính hợp lệ của các request từ client để tránh nguy cơ bảo mật, phá rối hệ thống. Từ câu chuyện xác thực request đã nảy sinh các vấn đề xung quanh khác mà chúng ta sẽ giải quyết tiếp sau đây.

Có lẽ bạn đã bắt gặp đâu đó tấm hình mà mọi người dùng để  chế giễu cấu trúc của PHP:

<div align="center">
	<img src="../images/ch5-04-validate.jpg">
	<br/>
	<span align="center">
		<i>'Hadouken' if-else</i>
	</span>
    <br/>
</div>

Thực tế đây là một trường hợp không liên quan gì tới ngôn ngữ mà chỉ là cách tổ chức code rườm rà khi gặp trường hợp mà nhiều field cần phải validate.

Trong phần này chúng ta sẽ dùng Go để viết một ví dụ validate và xem xét cải tiến nó theo 2 bước. Cuối cùng là phân tích cơ chế để hiểu rõ hơn cách một validator hoạt động.

## 4.4.1 Cải tiến 1: Tái cấu trúc hàm validation

Giả sử dữ liệu được liên kết tới một struct cụ thể thông qua binding bằng một thư viện Open source.

```go
type RegisterReq struct {
    // tag giúp json package encode giá trị của Username
    // thành giá trị tương ứng với key username trong json obj
    Username       string   `json:"username"`
    PasswordNew    string   `json:"password_new"`
    PasswordRepeat string   `json:"password_repeat"`
    Email          string   `json:"email"`
}

// register nhận vào obj kiểu RegisterReq và thực hiện validate
// các trường trong đó.
func register(req RegisterReq) error{
    if len(req.Username) > 0 {
        if len(req.PasswordNew) > 0 && len(req.PasswordRepeat) > 0 {
            if req.PasswordNew == req.PasswordRepeat {
                if emailFormatValid(req.Email) {
                    createUser()
                    return nil
                } else {
                    return errors.New("invalid email")
                }
            } else {
                return errors.New("password and reinput must be the same")
            }
        } else {
            return errors.New("password and password reinput must be longer than 0")
        }
    } else {
        return errors.New("length of username cannot be 0")
    }
}
```

Giờ code của chúng ta có vẻ khá giống một *"Hadouken"* nhắc ở phần đầu rồi, vậy làm thế nào để tối ưu đoạn code trên?

Có một giải pháp đã được đưa ra trong [Refactoring.com - Guard Clauses](https://refactoring.com/catalog/replaceNestedConditionalWithGuardClauses.html), thử áp dụng cho trường hợp của chúng ta:

```go
func register(req RegisterReq) error{
    if len(req.Username) == 0 {
        return errors.New("length of username cannot be 0")
    }

    if len(req.PasswordNew) == 0 || len(req.PasswordRepeat) == 0 {
        return errors.New("password and password reinput must be longer than 0")
    }

    if req.PasswordNew != req.PasswordRepeat {
        return errors.New("password and reinput must be the same")
    }

    if emailFormatValid(req.Email) {
        return errors.New("invalid email")
    }

    createUser()
    return nil
}
```

Nhờ bỏ đi cách viết if-else lồng nhau mà code trở nên "clean" hơn. Tuy vậy chúng ta vẫn phải viết khá nhiều hàm validate cho mỗi field trong một kiểu request.

Có một cách giúp chúng ta giảm khá nhiều code là sử dụng validator.

## 4.4.2 Cải tiến 2: Sử dụng validator

<div align="center">
	<img src="../images/validator.png" width="100">
	<br/>
</div>

Thư viện [validator](https://github.com/go-playground/validator) hỗ trợ việc validate bằng cách sử dụng các tag lúc định nghĩa struct. Một ví dụ nhỏ:

```go
import (
    "gopkg.in/go-playground/validator.v9"
    "fmt"
)

// RegisterReq là struct cần được validate
type RegisterReq struct {
    // gt = 0 cho biết độ dài chuỗi phải > 0，gt: greater than
    Username       string   `json:"username" validate:"gt=0"`
    PasswordNew    string   `json:"password_new" validate:"gt=0"`

    // eqfield kiểm tra các trường bằng nhau
    PasswordRepeat string   `json:"password_repeat" validate:"eqfield=PasswordNew"`

    // kiểm tra định dạng email thích hợp
    Email          string   `json:"email" validate:"email"`
}

// dùng 1 instance của Validate, cache lại struct info
var validate *validator.Validate

// validatefunc để wrap hàm validate.Struct
func validatefunc(req RegisterReq) error {
    err := validate.Struct(req)
    if err != nil {
        return err
    }
    return nil
}

func main() {
    validate = validator.New()

    // khởi tạo obj để test validator
    a := RegisterReq{
        Username        : "Alex",
        PasswordNew     : "",
        PasswordRepeat  : "z",
        Email           : "z@z.z",
    }

    err := validatefunc(a)
    fmt.Println(err)
}

// kết quả:
// Key: 'RegisterReq.PasswordNew' Error:Field validation for 'PasswordNew' failed on the 'gt' tag
// Key: 'RegisterReq.PasswordRepeat' Error:Field validation for 'PasswordRepeat' failed on the 'eqfield' tag
```

Một lưu ý nhỏ là  error message trả về cho người dùng thì không nên viết trực tiếp bằng tiếng Anh mà thông tin về error nên được tổ chức theo từng tag để người dùng theo đó tra cứu.

## 4.4.3 Cơ chế của validator

Từ quan điểm cấu trúc, mỗi struct có thể được xem như một cây. Giả sử chúng ta có một struct được định nghĩa như sau:

```go
type Nested struct {
    Email string `validate:"email"`
}
type T struct {
    Age    int `validate:"eq=10"`
    Nested Nested
}
```

Sẽ được vẽ thành một cây như bên dưới:

<div align="center">
	<img src="../images/ch5-04-validate-struct-tree.png" width="300">
	<br/>
	<span align="center">
		<i>Cây validator</i>
	</span>
</div>

Việc validate các trường có thể thực hiện khi đi qua cấu trúc cây này (bằng cách duyệt theo chiều sâu hoặc theo chiều rộng). Tiếp theo chúng ta sẽ minh hoạ cơ chế validate trên một cấu trúc như thế, mục đích để hiểu rõ hơn cách mà validator thực hiện.

Đầu tiên xác định 2 struct như hình trên:

```go
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
```

Định nghĩa hàm validate:

```go
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

```

Sau đây là cách sử dụng trong hàm main:

```go
func main() {
    // khởi tạo obj để test
    var a = T{Age: 10, Nested: Nested{Email: "abc@adfgom"}}

    validateResult, errmsg := validate(a)
    fmt.Println(validateResult, errmsg)
}

// kết quả:
// false validate mail failed, field val is: abc@adfgom
```

Thư viện validator được giới thiệu trong phần trước phức tạp hơn về mặt chức năng so với ví dụ ở đây. Nhưng nguyên tắc chung cũng là duyệt cây của một struct với reflection.

## 4.4.4 Xác thực request bằng JWT

Phần trên đã trình bày quá trình validate các thông tin về email và password khi đăng ký một tài khoản. Sau đó, nếu họ đăng nhập vào tài khoản bằng email và password thì trạng thái phiên làm việc của họ sẽ được giữ cho các yêu cầu kế tiếp. Có một số giải pháp để lưu trữ phiên làm việc bằng [session/cookie](https://astaxie.gitbooks.io/build-web-application-with-golang/en/06.1.html), một giải pháp khác là dùng cơ chế cấp token [JWT](https://jwt.io/) sau khi đăng nhập, và dùng token này để xác thực các yêu cầu về sau.

Không chỉ lưu trữ phiên làm việc, token JWT cũng hay đi kèm trong các lệnh gọi API để xác thực phía client khi gọi đến web service. Sau đây là một đoạn chương trình middleware xác thực yêu cầu bằng JWT:

***[auth.go](https://github.com/thoainguyen/go-hackercamp/blob/master/go-contacts/app/auth.go):***

```go
var JwtAuthentication = func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // danh sách các API không cần xác thực bằng token
        notAuth := []string{"/api/user/new", "/api/user/login"}
        requestPath := r.URL.Path
        for _, value := range notAuth {
            if value == requestPath {
                next.ServeHTTP(w, r)
                return
            }
        }

        response := make(map[string] interface{})
        tokenHeader := r.Header.Get("Authorization") 
        // thiếu jwt token, trả về lỗi
        if tokenHeader == "" {
            response = u.Message(false, "Missing auth token")
            w.WriteHeader(http.StatusForbidden)
            w.Header().Add("Content-Type", "application/json")
            u.Respond(w, response)
            return
        }
        // thông thường chuỗi token có định dạng: Bearer {token-body}, nên cần tách phần token ra
        splitted := strings.Split(tokenHeader, " ")
        if len(splitted) != 2 {
            response = u.Message(false, "Invalid/Malformed auth token")
            w.WriteHeader(http.StatusForbidden)
            w.Header().Add("Content-Type", "application/json")
            u.Respond(w, response)
            return
        }
        // chuỗi jwt token trong phần header của request
        tokenPart := splitted[1]
        tk := &models.Token{}

        token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("token_password")), nil
        })

        if err != nil {
            response = u.Message(false, "Malformed authentication token")
            w.WriteHeader(http.StatusForbidden)
            w.Header().Add("Content-Type", "application/json")
            u.Respond(w, response)
            return
        }

        if !token.Valid {
            response = u.Message(false, "Token is not valid.")
            w.WriteHeader(http.StatusForbidden)
            w.Header().Add("Content-Type", "application/json")
            u.Respond(w, response)
            return
        }
        ctx := context.WithValue(r.Context(), "user", tk.UserId)
        r = r.WithContext(ctx)
        // tiếp tục thực hiện request
        next.ServeHTTP(w, r)
    });
}

// github: https://github.com/thoainguyen/go-hackercamp/tree/master/go-contacts
```

<div style="display: flex; justify-content: space-around;">
<span> <a href="ch4-03-middleware.md">&lt Phần 4.3</a>
</span>
<span><a href="../SUMMARY.md"> Mục lục</a>  </span>
<span> <a href="ch4-05-database.md">Phần 4.5 &gt</a> </span>
</div>