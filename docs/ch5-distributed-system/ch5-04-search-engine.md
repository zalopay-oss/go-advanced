# 5.4 Công cụ tìm kiếm phân tán

<div align="center">
	<img src="../images/ch6-dis-search.png" width="500">
	<br/>
	<span align="center">
		<i>Distributed search engine</i>
	</span>
</div>

Trong chương Web, chúng ta đã biết MySQL có rất nhiều rủi ro. Nó là hệ thống cơ sở dữ liệu đảm bảo tính thời gian thực và tính nhất quán cao, do đó các chức năng của MySQL được thiết kế để đáp ứng tính nhất quán này. Ví dụ: thiết kế ghi nhật ký trước (write ahead log), index, tổ chức dữ liệu dựa trên [B+ tree](https://en.wikipedia.org/wiki/B%2B_tree) và transaction dựa trên [MVCC](https://en.wikipedia.org/wiki/Multiversion_concurrency_control).

Cơ sở dữ liệu quan hệ ([Relational database](https://en.wikipedia.org/wiki/Relational_database)) thường được sử dụng để triển khai các hệ thống OLTP. OLTP là gì? Theo như wikipedia định nghĩa:

Trong các kịch bản business, sẽ có những kịch bản mà yêu cầu tính thời gian thực không cao (có thể chấp nhận độ trễ vài giây), nhưng lại có độ phức tạp khi truy vấn cao. Ví dụ: trong hệ thống thương mại điện tử [WMS](https://vi.wikipedia.org/wiki/Hệ_thống_quản_lý_kho), các hệ thống [CRM](https://vi.wikipedia.org/wiki/Quản_lý_quan_hệ_khách_hàng) hoặc các dịch vụ khách hàng có kịch bản kinh doanh phức tạp, lúc này ta cần phải tạo ra các hàm truy vấn kết hợp rất nhiều trường khác nhau. Kích thước dữ liệu của hệ thống như vậy cũng rất lớn, chẳng hạn như mô tả hàng hóa trong hệ thống thương mại điện tử WMS, có thể có các trường sau:

> warehouse id, warehousing time, location id, storage shelf id, warehousing operator Id, outbound operator id, stock quantity, expiration time, SKU type, product brand, product category, number of internals

Ngoài các thông tin trên, nếu hàng hóa nằm trong kho. Có thể có id quá trình đang thực hiện, trạng thái hiện tại, ....

Hãy tưởng tượng nếu chúng ta đang điều hành một công ty thương mại điện tử lớn với hàng chục triệu đơn hàng mỗi ngày, sẽ rất khó để truy vấn và xây dựng các index thích hợp trong cơ sở dữ liệu này.

Trong CRM hoặc hệ thống dịch vụ khách hàng, thường có nhu cầu tìm kiếm theo từ khóa, bên cạnh đó thì các công ty lớn sẽ nhận được hàng chục ngàn khiếu nại của người dùng mỗi ngày. Từ đó, ta thấy rằng các khiếu nại của người dùng phải có ít nhất 2 đến 3 năm để có thể xử lý xong. Số lượng dữ liệu có thể là hàng chục triệu hoặc thậm chí hàng trăm triệu. Thực hiện một truy vấn dựa trên từ khóa có thể trực tiếp làm treo toàn bộ MySQL.

Lúc này, chúng ta cần một công cụ tìm kiếm có thể xử lý tốt trường hợp này.

## 5.4.1 Công cụ tìm kiếm

[Elaticsearch](https://www.elastic.co/) là công cụ dẫn đầu trong các công cụ tìm kiếm phân tán Open source. Elaticsearch dựa trên việc triển khai [apache Lucene](https://lucene.apache.org) và kết hợp cùng nhiều tối ưu hóa trong quá trình triển khai, vận hành và bảo trì. Xây dựng một công cụ tìm kiếm phân tán ngày nay dễ dàng hơn nhiều so với thời đại trước. Đơn giản chỉ cần cấu hình IP máy khách và cổng.

### 5.4.2 Inverted index

Mặc dù Elaticsearch được tạo ra cho mục đích tìm kiếm, nhưng Elaticsearch thường được sử dụng làm cơ sở dữ liệu trong các ứng dụng thực tế, vì tính chất của [inverted index](https://en.wikipedia.org/wiki/Inverted_index). Hiểu đơn giản thì dữ liệu là 1 quyển sách, để tìm kiếm nhanh thì người ta sinh ra 1 cái là mục lục đánh dấu nội dung, thì cái mục lục bản chất giống như việc đánh index vậy.

Việc đánh index có vẻ đơn giản nhưng bên dưới Elaticsearch làm khá nhiều việc. Mysql thường sẽ đánh theo các trường trong bảng như (name, email …). Tuy nhiên Elaticsearch sẽ đánh index theo đơn vị là term, cụ thể như sau:

```sh
Title (A1) = "Advanced Go Book"
Question (A2) = "What is Go"
Answer (A3) = "Go is simple"
```

Inverted index sẽ như sau:

```sh
"Advanced" => {A1}
"Go" => {A1,A2,A3}
"Book" => {A1}
"What" => {A2}
"is" => {A2,A3}
"simple" => {A3}
```

Ta thấy được chuỗi ban đầu là tổ hợp của nhiều Term. Và việc tìm kiếm sẽ dựa trên tổ hợp các term này. Nhưng làm sao Elaticsearch tách được chuỗi thành các Term? Câu trả lời là Elaticsearch sử dụng 2 kỹ thuật:

1. **N-Gram Morphological Analysis:** là kỹ thuật chia các chuỗi to thành các chuỗi con theo trọng số với độ dài N, N = (1..3), ví dụ N = 2 (cấu hình mặc định của Elaticsearch), khi tách chuỗi "ADVANCED GO BOOK" ta sẽ được các term như sau:

```sh
"ADVANCED GO BOOK" => {"AD","DV","VA","AN","NC","CE","ED","D "," G","GO","O "," B","BO","OO","OK"}
```

2. **[Morphological Analysis](<https://en.wikipedia.org/wiki/Morphology_(linguistics)>)** là kỹ thuật xử lý ngôn ngữ tự nhiên (National Language Procesing). Đơn giản là kỹ thuât tách các chuỗi thành từ có nghĩa dựa theo ngôn ngữ, ví dụ "ADVANCED GO BOOK" sẽ được phân tích như sau:

```sh
"ADVANCED GO BOOK" => {"ADVANCED", "GO","BOOK}
```

### 5.4.3 Truy vấn DSL ([Domain-specific Language](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html))

Khi chúng ta sử dụng Elaticsearch làm cơ sở dữ liệu, chúng ta cần sử dụng truy vấn Bool của nó. Ví dụ:

```json
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "field_1": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_2": {
              "query": "2",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_3": {
              "query": "3",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_4": {
              "query": "4",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```

Trông khá phực tạp, nhưng ta có thể hiểu đơn giản như sau:

```go
if field_1 == 1 && field_2 == 2 && field_3 == 3 && field_4 == 4 {
    return true
}
```

Logic OR trong truy vấn bool:

```json
{
  "query": {
    "bool": {
      "should": [
        {
          "match": {
            "field_1": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "field_2": {
              "query": "3",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```

Code go cho logic trên:

```go
if field_1 == 1 || field_2 == 2 {
	return true
}
```

Các biểu thức theo sau `if` trong các đoạn code Go ở trên để diễn đạt [Boolean Expression](https://en.wikipedia.org/wiki/Boolean_expression):

```go
4 > 1
5 == 2
3 < i && x > 10
```

`Bool Query` là dùng json để diễn tả Boolean Expression, tại sao sao lại sử dụng nó? Vì json có thể biểu thị cấu trúc cây, code của chúng ta sẽ trở thành [AST](https://en.wikipedia.org/wiki/Abstract_syntax_tree) (Abstract Syntax Tree) sau khi được trình biên dịch phân tích cú pháp và trở thành cây cú pháp trừu tượng AST. Boolean Expression ở đây được tạo bởi trình biên dịch Parse và kết quả là một cấu trúc cây, và đây chỉ là một bước nhỏ trong toàn bộ quá trình thực hiện của trình biên dịch.

### 5.4.4 Sử dụng client SDK của Elaticsearch

Khởi tạo client trong Elaticsearch:

```go
import (
  // sử dụng elastic version 3
	elastic "gopkg.in/olivere/elastic.v3"
)

var esClient *elastic.Client

func initElasticsearchClient(host string, port string) {
	var err error
	esClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", host, port)),
		elastic.SetMaxRetries(3),
	)

	if err != nil {
		// log error
	}
}
```

Chèn một document vào Elaticsearch:

```go
func insertDocument(db string, table string, obj map[string]interface{}) {

  id := obj["id"]

  var indexName, typeName string
  // database/table trong cơ sở dữ liệu được ánh xạ (mapping) vào index và type của Elaticsearch
  // lưu ý _type trong Elaticsearch chỉ là một field document
  // sử dụng nhiều index sẽ dẫn đến việc giảm hiệu suất
  // ở phiên bản mới thì type không con sử dụng
  // để làm cho dữ liệu của các bảng khác nhau có index khác nhau, trong ví dụ này chúng tôi sử dụng table+name làm index

  indexName = fmt.Sprintf("%v_%v", db, table)
  typeName = table

  // thực hiện việc chèn
  res, err := esClient.Index().Index(indexName).Type(typeName).Id(id).BodyJson(obj).Do()
  if err != nil {
    // xử lý error
  } else {
    // xủ lý success
  }
}
```

Để lấy dữ liệu chúng ta làm như sau:

```go
func query(indexName string, typeName string) (*elastic.SearchResult, error) {
	// thêm điều kiện truy vấn bằng bool must và bool should 
	q := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("id", 1),
	elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("male", "m")))

	q = q.Should(
		elastic.NewMatchPhraseQuery("name", "alex"),
		elastic.NewMatchPhraseQuery("name", "xargin"),
	)

	searchService := esClient.Search(indexName).Type(typeName)
	res, err := searchService.Query(q).Do()
	if err != nil {
		// log error
		return nil, err
	}

	return res, nil
}
```

Để xóa dữ liệu chúng ta làm như sau:

```go
func deleteDocument(
	indexName string, typeName string, obj map[string]interface{},
) {
	id := obj["id"]

	res, err := esClient.Delete().Index(indexName).Type(typeName).Id(id).Do()
	if err != nil {
		// xử lý error
	} else {
		// xử lý success
	}
}
```

Do bản chất của Lucene, dữ liệu trong công cụ tìm kiếm này là bất biến. Vì vậy nếu bạn muốn cập nhật tài liệu, thực chất việc chèn sẽ diễn ra.

Khi sử dụng Elaticsearch làm cơ sở dữ liệu, bạn cần lưu ý rằng Elaticsearch có hoạt động hợp nhất index, vì vậy phải mất một thời gian để dữ liệu được chèn vào Elaticsearch mới có thể truy vấn được (cấu hình refresh_interval của Elaticsearch). Vì vậy, không sử dụng Elaticsearch như một cơ sở dữ liệu quan hệ [strong consistency](https://en.wikipedia.org/wiki/Strong_consistency).

### 5.4.5 Chuyển đổi SQL sang DSL

Ví dụ: chúng ta cần một biểu thức bool `user_id = 1 and (product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1)`, SQL của nó như sau:

```sql
select * from xxx where user_id = 1 and (
	product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1
)
```

DSL trong Elaticsearch có dạng:

```json
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "user_id": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "match": {
            "product_id": {
              "query": "1",
              "type": "phrase"
            }
          }
        },
        {
          "bool": {
            "should": [
              {
                "match": {
                  "star_num": {
                    "query": "4",
                    "type": "phrase"
                  }
                }
              },
              {
                "match": {
                  "star_num": {
                    "query": "5",
                    "type": "phrase"
                  }
                }
              }
            ]
          }
        },
        {
          "match": {
            "banned": {
              "query": "1",
              "type": "phrase"
            }
          }
        }
      ]
    }
  },
  "from": 0,
  "size": 1
}
```

Dù bạn hiểu rõ DSL của Elaticsearch nhưng vẫn sẽ tốn công để viết được nó. Chúng ta đã biết client SDK của Elaticsearch nhưng nó cũng không đủ linh hoạt.

Phần WHERE trong SQL là Boolean Expression. Như chúng ta đã biết, bool expresion này tương tự như cấu trúc DSL của Elaticsearch sau khi được phân tích cú pháp. Vậy có thể chuyển đổi qua lại giữa DSL và SQL không?

Câu trả lời là chắc chắn được. Bây giờ, chúng ta thử so sánh cấu trúc của SQL với cấu trúc của Parse và cấu trúc của DSL của Elaticsearch :

<div align="center">
	<img src="../images/ch6-ast-dsl.png">
	<br/>
	<span align="center">
		<i>Sự tương ứng giữa AST và DSL</i>
	</span>
</div>

Ta thấy cấu trúc chúng khá giống nhau nên chúng ta có thể chuyển đổi logic của chúng cho nhau. Trước tiên, chúng ta duyệt cây AST theo chiều rộng, sau đó chuyển đổi biểu thức nhị phân thành chuỗi json và tổng hợp nó lại.

Do quá trình khá phức tạp nên ví dụ không được đưa vào bài viết này. Vui lòng tham khảo [ở đây](github.com/cch123/elasticsql) để biết thêm chi tiết.

### 5.4.6 Đồng bộ hóa dữ liệu không đồng nhất

Trong các ứng dụng thực tế, chúng ta hiếm khi ghi dữ liệu trực tiếp vào công cụ tìm kiếm. Một cách phổ biến hơn là đồng bộ hóa dữ liệu từ MySQL hoặc loại databsae khác vào công cụ tìm kiếm. Người dùng của công cụ tìm kiếm chỉ có thể truy vấn dữ liệu mà không thể sửa đổi và xóa nó.

Có hai chương trình đồng bộ hóa phổ biến:

#### 5.4.1 Đồng bộ dữ liệu theo timestamp

<div align="center">
	<img src="../images/ch6-sync.png">
	<br/>
	<span align="center">
		<i>Dựa trên timestamp để đồng bộ hóa dữ liệu</i>
	</span>
</div>


Phương thức đồng bộ hóa này cần phải phù hợp với nhu cầu của business. Ví dụ, đối với đơn hàng trong hệ thống WMS, chúng ta không cần tính `realtime` cao và việc xử lý chậm có thể chấp nhận được. Vì vậy, chúng tôi có thể xử lý đơn hàng cứ mỗi 10 phút, logic cụ thể giống câu SQL sau:

```sql
select * from wms_orders where update_time >= date_sub(now(), interval 10 minute);
```

Khi xem xét các giá trị biên, chúng ta nên lấy dữ liệu với khoảng thời gian trùng với một phần khoảng thời gian trước đó:

```sql
select * from wms_orders where update_time >= date_sub(
	now(), interval 11 minute
);
```

Sau khi tăng khoảng thời gian lên 11 phút thì mọi thứ đã ổn hơn. Nhưng rõ ràng, phương pháp này có khá nhiều thiếu sót và có điều kiện về tính chất thời gian. Ví dụ: phải có trường `update_time` và cập nhật nó mỗi khi tạo hoặc cập nhật, và giá trị thời gian này phải chính xác. Nếu không việc đồng bộ hoá có thể mất dữ liệu.

### 5.4.2 Đồng bộ hóa dữ liệu với binlog

<div align="center">
	<img src="../images/ch6-binlog-sync.png">
	<br/>
	<span align="center">
		<i>Đồng bộ hóa dữ liệu dựa trên binlog</i>
	</span>
</div>

[Canal](https://github.com/alibaba/canal) là Open source của **Alibaba** và nó được dùng để phân tích cú pháp `binlog` và đồng bộ hóa bởi nhiều công ty. Canal sẽ hoạt động như một thư viện phụ thuộc MySQL, nó sẽ phân tích cú pháp `bincode` của từng dòng và gửi nó đến hàng đợi tin nhắn theo định dạng dễ hiểu hơn (chẳng hạn như json).

[Downstream Kafka consumer](https://kafka.apache.org/20/documentation.html) chịu trách nhiệm ghi khóa chính tự tăng của bảng dữ liệu upstream dưới dạng ID của Elaticsearch. Mỗi khi nhận được binlog, dữ liệu có ID tương ứng sẽ được cập nhật mới nhất. Binlog của một row trong MySQL sẽ cung cấp tất cả các trường của record cho downstream. Vì vậy, trên thực tế, khi đồng bộ hóa dữ liệu, bạn không cần xem xét dữ liệu được chèn hay cập nhật, miễn là bạn có chèn ID.

Mô hình này cũng yêu cầu business tuân thủ một điều kiện của bảng dữ liệu, bảng phải có ID khóa chính là `duy nhất` để đảm bảo rằng dữ liệu chúng ta nhập vào Elaticsearch sẽ không bị trùng lặp. Khi điều kiện này không được tuân theo, nó sẽ dẫn đến sự trùng lặp dữ liệu khi đồng bộ hóa.

[Tiếp theo](ch5-05-load-balance.md)