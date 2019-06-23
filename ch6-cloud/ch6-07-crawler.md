# 6.7 Trình thu thập thông tin phân tán

Sự bùng nổ thông tin trong kỷ nguyên Internet là một vấn đề khiến nhiều người cảm thấy đau đầu. Vô số tin tức, thông tin và video đang xâm chiếm thời gian của chúng ta. Mặt khác, khi chúng ta thực sự cần dữ liệu, chúng ta cảm thấy rằng dữ liệu không dễ dàng gì có được. Ví dụ: chúng tôi muốn biết về những gì mọi người đang thảo luận và những gì họ quan tâm. Nhưng chúng ta không có thời gian để đọc từng cuốn tiểu thuyết yêu thích, lúc này chúng ta muốn sử dụng công nghệ để đưa những thông tin cần vào cơ sở dữ liệu của chúng ta. Quá trình này có thể tốn một vài tháng hoặc một năm. Cũng có thể chúng tôi muốn lưu những thông tin hữu ích và thoáng qua trên Internet, chẳng hạn như các cuộc thảo luận chất lượng cao của những người giỏi được tập hợp trong một diễn đàn rất nhỏ. Vào một lúc nào đó trong tương lai, chúng ta tìm lại được những thông tin đó và rút ra được những kết luận mà chúng ta chưa từng có.

Ngoài nhu cầu giải trí, có rất nhiều tài liệu mở quý giá trên Internet. Trong những năm gần đây, học sâu đã và đang rất `hot`, nhưng học máy thường không quan tâm đến việc thiết lập đúng hay không, các thông số có được điều chỉnh chính xác không, mà là ở Giai đoạn khởi tạo ban đầu: không có dữ liệu.

Một công việc nhập môn của thu thập thông tin, khả năng viết một chượng trình thu thập thông tin đơn giản hoặc phức tạp rất quan trọng.

## 6.7.1 Trình thu thập thông tin độc lập dựa trên collly

"Lập trình ngôn ngữ Go" đưa ra một ví dụ về trình thu thập thông tin đơn giản. Sau nhiều năm lập trình, việc dùng Go sẽ cực kì thuận tiện để viết một trình thu thập thông tin cho trang web, chẳng hạn như việc thu thập thông tin trang web (www.abcdefg.com là trang web ảo):


```go
package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

var visited = map[string]bool{}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.abcdefg.com"),
		colly.MaxDepth(1),
	)

	// We think the matching page is the details page of the site
	detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
	// Matching the following pattern is the list page of the site
	listRegex, _ := regexp.Compile(`/t/\d+#\w+`)

	// All a tags, set the callback function
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Visited details page or list page, skipped
		if visited[link] && (detailRegex.Match([]byte(link)) || listRegex.Match([]byte(link))) {
			return
		}

		// is neither a list page nor a detail page
		// Then it's not what we care about, skip it
		if !detailRegex.Match([]byte(link)) && !listRegex.Match([]byte(link)) {
			println("not match", link)
			return
		}

		// Because most websites have anti-reptile strategies
		// So there should be sleep logic in the crawler logic to avoid being blocked
		time.Sleep(time.Second)
		println("match", link)

		visited[link] = true

		time.Sleep(time.Millisecond * 2)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	err := c.Visit("https://www.abcdefg.com/go/go")
	if err != nil {fmt.Println(err)}
}
```

## 6.7.2 Trình thu thập thông tin phân tán

Hãy tưởng tượng rằng hệ thống phân tích thông tin của bạn đang chạy rất nhanh. Tốc độ thu thập thông tin đã trở thành nút cổ chai. Mặc dù bạn có thể sử dụng tất cả các tính năng đồng thời tuyệt vời  Go để sử dụng hết hiệu suất CPU và băng thông mạng, và bạn vẫn muốn tăng tốc độ thu thập thông tin của trình thu thập thông tin. Trong nhiều ngữ cảnh, tốc độ mang nhiều ý nghĩa:

1. Đối với thương mại điện tử, cụ thể là cuộc chiến giá cả, tôi sẽ hy vọng mình có được giá mới nhất của đối thủ khi chúng thay đổi, và hệ thống sẽ tự động điều chỉnh giá của sản phẩm sao cho phù hợp.
2. Đối với các dịch vụ cung cấp thông tin `Feed`, tính kịp thời của thông tin rất quan trọng. Nếu tin tức thu thập là tin tức của ngày hôm qua, nó sẽ không có ý nghĩa gì với người dùng.

Vì vậy, chúng ta cần hệ thống thu thập thông tin phân tán. Về bản chất, các trình thu thập thông tin phân tán là tập hợp các hệ thống phân phối và thực thi tác vụ. Trông hệ thống phân phối tác vụ phổ biến, sẽ có sự sai lệch tốc độ giữa upstream và downstream nên sẽ luôn có hàng đợi tin nhắn.

![dist-crawler](../images/ch6-dist-crawler.png)

*Hình 6-14 Luồng công việc*

Công việc chính của upstream là thu thập thông tin tất cả các "trang" đích từ một điểm bắt đầu được cấu hình sẵn. Nội dung html của trang danh sách sẽ chứa các liên kết đến các trang chi tiết. Số lượng trang chi tiết thường gấp 10 đến 100 lần so với trang danh sách, vì vậy chúng tôi xem các trang chi tiết này như "tác vụ" và phân phối chúng thông qua hàng đợi tin nhắn.

Để thu thập thông tin trang, điều quan trọng là không có sự lặp lại thường xuyên trong quá trình thực thi, vì nó sẽ tạo kết quả bất định (ví dụ trên sẽ chỉ thu thập nội dung trang chứ không phải phần bình luận).

Trong phần này, chúng ta sẽ hiện thực trình thu thập thông tin đơn giản dựa trên hàng đợi tin nhắn. Cụ thể ở đây là sử dụng các nats để phân phối tác vụ. Trong thực tế, tuỳ vào yêu cầu về độ tin cậy của thông điệp và cơ sở hạ tầng của công ty nên sẽ ảnh hưởng tới việc chọn công nghệ của từng doanh nghiệp.

### 6.7.2.1 Giới thiệu về nats

Nats là một hàng đợi tin nhắn phân tán hiệu suất cao được hiện thực bằng Go cho các tình huống yêu cầu tính đồng thời cao, `throughput` cao. Những phiên bản nats ban đầu mang thiên hướng về tốc độ và không hỗ trợ tính `persistence`. Kể từ 16 năm trước, nats đã hỗ trợ tính `persistence` dựa trên log thông qua nats-streaming, cũng như nhắn tin đáng tin cậy. Vì đây là những ví dụ đơn giản, chúng tôi chỉ sử dụng các nats trong phần này.

Máy chủ của nats là gnatsd. Phương thức giao tiếp giữa máy khách và gnatsd là giao thức văn bản dựa trên tcp:

Gửi tin nhắn đi có chứa chủ đề cho một tác vụ:

![nats-Protocol-pub](../images/ch6-09-nats-protocol-pub.png)

*Hình 6-15 Pub trong giao thức nats*

Theo dõi các tác vụ có chủ đề bằng hàng đợi các worker:

![nats-protocol-sub](../images/ch6-09-nats-protocol-sub.png)

*Hình 6-16 Sub trong giao thức nats*

Tham số hàng đợi là tùy chọn. Nếu bạn muốn tải cân bằng tác vụ ở phía người tiêu dùng phân tán, thay vì mọi người nhận cùng một thông báo, bạn nên gán cùng tên hàng đợi cho người tiêu dùng.

#### Sản xuất tin nhắn cơ bản

Thông điệp sản xuất có thể được chỉ định bằng cách chỉ định chủ đề:

`` `đi
Nc, err: = nats.Connect (nats.DefaultURL)
Nếu err! = Nil {return}

/ / Chỉ định chủ đề là nhiệm vụ, nội dung của tin nhắn là miễn phí
Err = nc.Publish ("tác vụ", [] byte ("nội dung nhiệm vụ của bạn"))

nc.Flush ()
`` `

#### Tiêu thụ tin nhắn cơ bản

Việc sử dụng trực tiếp API đăng ký của Nats không đạt được mục đích phân phối tác vụ, vì bản thân phụ pub được phát sóng. Tất cả người tiêu dùng sẽ nhận được chính xác cùng một thông điệp.

Ngoài đăng ký bình thường, nats cũng cung cấp chức năng đăng ký hàng đợi. Bằng cách cung cấp tên nhóm xếp hàng (tương tự như nhóm người tiêu dùng trong Kafka), các tác vụ có thể được phân phối cho người tiêu dùng một cách cân bằng.

`` `đi
Nc, err: = nats.Connect (nats.DefaultURL)
Nếu err! = Nil {return}

// đăng ký hàng đợi tương đương với cân bằng chi nhánh để phân phối nhiệm vụ giữa người tiêu dùng
// Tiền đề là tất cả người tiêu dùng sử dụng công nhân hàng đợi này
// Hàng đợi trong nats tương tự về mặt khái niệm với nhóm người tiêu dùng ở Kafka
Sub, err: = nc.QueueSubscribeSync ("tác vụ", "công nhân")
Nếu err! = Nil {return}

Thông điệp Var * nats.Msg
Dành cho {
	Msg, err = sub.NextMsg (thời gian. * 10000)
	Nếu err! = Nil {break}
	// tiêu thụ đúng thông điệp
	// có thể xử lý các tác vụ với đối tượng nats.Msg
}
`` `

## 6.7.3 Kết hợp việc tạo tin nhắn với nats và colly

Chúng tôi tùy chỉnh một trình thu thập tương ứng cho mỗi trang web và đặt các quy tắc tương ứng, chẳng hạn như abcdefg, Hijklmn (hư cấu), sau đó sử dụng một phương thức xuất xưởng đơn giản để ánh xạ trình thu thập tới máy chủ của nó. Mỗi trang web leo lên trang danh sách. Bạn cần phân tích tất cả các liên kết trong chương trình hiện tại và gửi URL của trang đích đến hàng đợi tin nhắn.

`` `đi
Gói chính

Nhập khẩu (
	"fmt"
	"mạng / url"

	"github.com/gocolly/colly"
)

Var domain2Collector = map [chuỗi] * colly.Collector {}
Var nc * nats. Kết nối
Biến tối đa = 10
Var natsURL = "nats: // localhost: 4222"

Nhà máy Func (chuỗi urlStr) * colly.Collector {
	u, _: = url.Pude (urlStr)
	Trả về domain2Collector [u.Host]
}

Func initABCDECollector () * colly.Collector {
	c: = colly.NewCollector (
		colly.AllowedDomains ("www.abcdefg.com"),
		colly.MaxDepth (maxDepth),
	)

	c.OnResponse (func (resp * colly.Response) {
		// Thực hiện một số công việc chăm sóc sau khi leo núi
		// Ví dụ: xác nhận rằng trang đã được thu thập thông tin được lưu trữ trong MySQL.
	})

	c.OnHTML ("a [href]", func (e * colly.HTMLEuity) {
		// Chiến lược chống bò sát cơ bản
		Liên kết: = e.Attr ("href")
		time.S ngủ (time.Second * 2)

		// trang danh sách phù hợp thông thường, sau đó truy cập
		Nếu listRegex.Match ([] byte (liên kết)) {
			c.Visit (e.Request.AbsoluteURL (liên kết))
		}
		// Trang đích phù hợp thông thường, gửi tin nhắn xếp hàng
		Nếu chi tiếtRegex.Match ([] byte (liên kết)) {
			Err = nc.Publish ("tác vụ", [] byte (liên kết))
			nc.Flush ()
		}
	})
	Trả lại c
}

Func initHIJKLCollector () * colly.Collector {
	c: = colly.NewCollector (
		colly.AllowedDomains ("www.hijklmn.com"),
		colly.MaxDepth (maxDepth),
	)

	c.OnHTML ("a [href]", func (e * colly.HTMLEuity) {
	})

	Trả lại c
}

Func init () {
	domain2Collector ["www.abcdefg.com"] = initV2exCollector ()
	domain2Collector ["www.hijklmn.com"] = initV2fxCollector ()

	Lỗi lỗi
	Nc, err = nats. Kết nối (natsURL)
	Nếu err! = Nil {os.Exit (1)}
}

Func chính () {
	Các Url: = [] chuỗi {"https://www.abcdefg.com", "https://www.hijklmn.com"}
	Đối với _, url: = phạm vi url {
		Sơ thẩm: = nhà máy (url)
		dụ.Visit (url)
	}
}

`` `

## 6.7.4 Kết hợp tiêu thụ tin nhắn với collly

Phía người tiêu dùng đơn giản hơn một chút, chúng tôi chỉ cần đăng ký chủ đề tương ứng và truy cập trực tiếp vào trang chi tiết của trang web (trang sàn).

`` `đi
Gói chính

Nhập khẩu (
	"fmt"
	"mạng / url"

	"github.com/gocolly/colly"
)

Var domain2Collector = map [chuỗi] * colly.Collector {}
Var nc * nats. Kết nối
Biến tối đa = 10
Var natsURL = "nats: // localhost: 4222"

Nhà máy Func (chuỗi urlStr) * colly.Collector {
	u, _: = url.Pude (urlStr)
	Trả về domain2Collector [u.Host]
}

Func initV2exCollector () * colly.Collector {
	c: = colly.NewCollector (
		colly.AllowedDomains ("www.abcdefg.com"),
		colly.MaxDepth (maxDepth),
	)
	Trả lại c
}

Func initV2fxCollector () * colly.Collector {
	c: = colly.NewCollector (
		colly.AllowedDomains ("www.hijklmn.com"),
		colly.MaxDepth (maxDepth),
	)
	Trả lại c
}

Func init () {
	domain2Collector ["www.abcdefg.com"] = initABCDECollector ()
	domain2Collector ["www.hijklmn.com"] = initHIJKLCollector ()

	Lỗi lỗi
	Nc, err = nats. Kết nối (natsURL)
	Nếu err! = Nil {os.Exit (1)}
}

Func startConsumer () {
	Nc, err: = nats.Connect (nats.DefaultURL)
	Nếu err! = Nil {return}

	Sub, err: = nc.QueueSubscribeSync ("tác vụ", "công nhân")
	Nếu err! = Nil {return}

	Thông điệp Var * nats.Msg
	Dành cho {
		Msg, err = sub.NextMsg (thời gian. * 10000)
		Nếu err! = Nil {break}

		urlStr: = chuỗi (dir.Data)
		Ins: = nhà máy (urlStr)
		// Bởi vì phần hạ lưu nhất phải là trang đích của trang web tương ứng.
		// Vì vậy, don lồng phải đưa ra những đánh giá bổ sung, chỉ cần leo lên nội dung trực tiếp.
		in.Visit (urlStr)
		// ngăn chặn bị chặn
		time.S ngủ (time.Second)
	}
}

Func chính () {
	startConsumer ()
}
`` `

Ở cấp độ mã, các nhà sản xuất và người tiêu dùng ở đây về cơ bản là giống nhau. Nếu chúng tôi muốn linh hoạt hỗ trợ tăng và giảm thu thập thông tin của các trang web khác nhau trong tương lai, chúng tôi nên suy nghĩ về cách định cấu hình các chiến lược và tham số của trình thu thập thông tin này càng nhiều càng tốt.

Việc sử dụng một số hệ thống cấu hình đã được đề cập trong phần Cấu hình phân tán của chương này và độc giả có thể tự mình dùng thử, vì vậy tôi sẽ không đi sâu vào chi tiết ở đây.