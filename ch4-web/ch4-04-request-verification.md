# 4.4 Kiểm tra yêu cầu validator

Một số lập trình viên thích chế giễu cấu trúc của PHP bằng hình sau:

<div align="center">
	<img src="../images/ch5-04-validate.jpg">
	<br/>
	<span align="center">
		<i>Quá trình validator</i>
	</span>
</div>
<br/>

Thực tế đây là một trường hợp không liên quan gì tới ngôn ngữ. Có nhiều trường hợp mà các trường cần phải xác nhận (validate). Form hoặc JSON submit chỉ là một ví dụ điển hình. Chúng ta sử dụng Go để viết một ví dụ validate giống với ở trên, sau đó sẽ xem xét để cải thiện nó theo từng bước.

## 4.4.1 Tái cấu trúc hàm request validation

Giả sử dữ liệu được liên kết tới một struct cụ thể thông qua binding bằng một thư viện opensource.

```go
type RegisterReq struct {
    Username       string   `json:"username"`
    PasswordNew    string   `json:"password_new"`
    PasswordRepeat string   `json:"password_repeat"`
    Email          string   `json:"email"`
}

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

Làm thế nào để tối ưu đoạn code trên?

Rất đơn giản, có một giải pháp đã được đưa ra trong [Refactoring: Guard Clauses](https://refactoring.com/catalog/replaceNestedConditionalWithGuardClauses.html)

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

Thế là đoạn code trở nên "clean" hơn và nhìn bớt kì cục. Mặc dù phương thức tái cấu trúc được sử dụng để làm cho code của quy trình validate trông thanh lịch hơn, chúng ta vẫn phải viết một tập các hàm tương tự như `validate()` cho mỗi yêu cầu `http`. Có cách nào tốt hơn để giúp chúng ta cải thiện hơn không? Câu trả lời là sử dụng validator.

## 4.4.2 Cải tiến với validator

Từ quan điểm thiết kế, chúng ta chắc chắn sẽ phải khai báo một cấu trúc cho mỗi request. Các trường hợp validate được đề cập trong phần trước đều có thể được thực hiện thông qua validator. Đoạn code sau lấy lại struct trong phần trước làm ví dụ. Để cho gọn chúng ta sẽ bỏ qua thẻ json.

Ở đây ta sử dụng một thư viện validator mới: <https://github.com/go-playground/validator>

```go
import "gopkg.in/go-playground/validator.v9"

type RegisterReq struct {
    // gt = 0 cho biết độ dài chuỗi phải > 0，gt = greater than
    Username       string   `validate:"gt=0"`
    // như trên
    PasswordNew    string   `validate:"gt=0"`
    // eqfield kiểm tra các trường bằng nhau
    PasswordRepeat string   `validate:"eqfield=PasswordNew"`
    // kiểm tra định dạng email thích hợp
    Email          string   `validate:"email"`
}

validate := validator.New()

func validate(req RegisterReq) error {
    err := validate.Struct(req)
    if err != nil {
        doSomething()
        return err
    }
    ...
}
```

Điều này loại bỏ sự cần thiết phải viết các hàm `validate()` trùng lặp trước khi mỗi request đi vào logic nghiệp vụ. Trong ví dụ này, chỉ có một vài tính năng của validator này được liệt kê.

Ta thử thực thi chương trình này với tham số đầu vào được set:

```go
//...

var req = RegisterReq {
    Username       : "Xargin",
    PasswordNew    : "ohno",
    PasswordRepeat : "ohn",
    Email          : "alex@abc.com",
}

err := validate(req)
fmt.Println(err)

// Key: 'RegisterReq.PasswordRepeat' Error:Field validation for
// 'PasswordRepeat' failed on the 'eqfield' tag
```

Khi trả về error message cho người dùng thì không nên viết trực tiếp bằng tiếng Anh. Thông tin về error có thể được tổ chức theo từng tag và người đọc theo đó tự tìm hiểu.

## 4.4.3 Các nguyên tắc

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
	<img src="../images/ch5-04-validate-struct-tree.png">
	<br/>
	<span align="center">
		<i>Cây validator</i>
	</span>
</div>
<br/>

Việc validate các trường có thể đi qua cây cấu trúc này (bằng cách duyệt chiều sâu hoặc theo chiều rộng). Thử viết một ví dụ duyệt cây theo chiều sâu:

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
    Email string `validate:"email"`
}
type T struct {
    Age    int `validate:"eq=10"`
    Nested Nested
}

func validateEmail(input string) bool {
    if pass, _ := regexp.MatchString(
        `^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, input,
    ); pass {
        return true
    }
    return false
}

func validate(v interface{}) (bool, string) {
    validateResult := true
    errmsg := "success"
    vt := reflect.TypeOf(v)
    vv := reflect.ValueOf(v)
    for i := 0; i < vv.NumField(); i++ {
        fieldVal := vv.Field(i)
        tagContent := vt.Field(i).Tag.Get("validate")
        k := fieldVal.Kind()

        switch k {
        case reflect.Int:
            val := fieldVal.Int()
            tagValStr := strings.Split(tagContent, "=")
            tagVal, _ := strconv.ParseInt(tagValStr[1], 10, 64)
            if val != tagVal {
                errmsg = "validate int failed, tag is: "+ strconv.FormatInt(
                    tagVal, 10,
                )
                validateResult = false
            }
        case reflect.String:
            val := fieldVal.String()
            tagValStr := tagContent
            switch tagValStr {
            case "email":
                nestedResult := validateEmail(val)
                if nestedResult == false {
                    errmsg = "validate mail failed, field val is: "+ val
                    validateResult = false
                }
            }
        case reflect.Struct:
            // nếu có struct lồng bên trong thì truyền
            // xuống đệ quy theo chiều sâu
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
    var a = T{Age: 10, Nested: Nested{Email: "abc@abc.com"}}

    validateResult, errmsg := validate(a)
    fmt.Println(validateResult, errmsg)
}
```

Ví dụ trên kiểm tra định dạng email theo tag age, bạn có thể thực hiện thay đổi đơn giản trong chương trình để xem kết quả validate cụ thể. Để tinh giản việc xử lý lỗi và các tiến trình phức tạp, ví dụ `reflect.Int8/16/32/64`, `reflect.Ptr`, nếu bạn viết thư viện xác minh cho môi trường doanh nghiệp, hãy đảm bảo cải thiện chức năng và khả năng chịu lỗi.

Component validation opensource được giới thiệu trong phần trước phức tạp hơn về mặt chức năng so với ví dụ ở đây. Nhưng nguyên tắc chung rất đơn giản là duyệt cây của một struct với reflection. Việc ta phải sử dụng một lượng lớn các reflection khi verify struct nhưng vì reflection trong Go không hiệu quả lắm nên đôi khi sẽ ảnh hưởng đến hiệu suất của chương trình. Ngữ cảnh đòi hỏi nhiều verify struct thường xuất hiện trong các web service. Đây không hẳn là vấn đề thắt cổ chai hiệu năng của chương trình. Hiệu quả thực tế là đưa ra phán đoán chính xác hơn từ thư viện `pprof`.

Điều gì xảy ra nếu quá trình verify dựa trên reflection thực sự trở thành thắt cổ chai hiệu năng trong service của bạn? Có một ý tưởng để tránh dùng reflection: sử dụng "Trình phân tích cú pháp tích hợp của Go (Parser)" để quét mã nguồn và sau đó tạo mã xác thực dựa trên định nghĩa của struct. Chúng ta có thể đưa tất cả các struct cần được verify vào trong một package riêng. Vấn đề này được để lại cho người đọc khám phá.
