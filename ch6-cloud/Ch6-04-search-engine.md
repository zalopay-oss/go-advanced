# 6.4 Công cụ tìm kiếm phân tán

Trong chương Web, chúng tôi đã nói rằng MySQL rất yếu ớt. Bản thân hệ cơ sở dữ liệu này đảm bảo theo tính "real-time" và nhất quán mạnh mẽ, do đó chức năng chính của nó là để đáp ứng yêu cầu về tính nhất quán. Ví dụ: thiết kế ghi nhật ký trước (write ahead log), chỉ mục (index) và tổ chức dữ liệu dựa trên cây B+ và transaction dựa trên MVCC.

Cơ sở dữ liệu quan hệ thường được sử dụng để hiện thực các hệ thống OLTP, tạo sao gọi là OLTP, trích dẫn theo wikipedia:

> OLTP has also been used to refer to processing in which the system responds immediately to user requests. An automated teller machine (ATM) for a bank is an example of a commercial transaction processing application. Online transaction processing applications have high throughput and are insert- or update-intensive in database management. These applications are used concurrently by hundreds of users. The key goals of OLTP applications are availability, speed, concurrency and recoverability. Reduced paper trails and the faster, more accurate forecast for revenues and expenses are both examples of how OLTP makes things simpler for businesses. However, like many modern online information technology solutions, some systems require offline maintenance, which further affects the cost-benefit analysis of an online transaction processing system.

In the business scenario of the Internet, there are also scenarios where the real-time requirements are not high (a delay of many seconds can be accepted), but the query complexity is high. For example, in an e-commerce WMS system, or in most CRM or customer service systems with rich business scenarios, it may be necessary to provide random combination query functions for dozens of fields. The data dimensions of such a system are inherently numerous, such as a description of a piece of goods in an e-commerce WMS, which may have the following fields:
Trong kịch bản kinh doanh trên Internet, có những kịch bản mà yêu cầu thời gian thực không cao (có thể chấp nhận độ trễ vài giây), nhưng độ phức tạp khi truy vấn cao. Ví dụ: trong hệ thống WMS của thương mại điện tử hoặc trong hầu hết các hệ thống CRM hoặc dịch vụ khách hàng có kịch bản kinh doanh đa dạng, nó có thể cần phải cung cấp các hàm truy vấn kết hợp ngẫu nhiên cho rất nhiều trường. Kích thước dữ liệu của một hệ thống như vậy vốn đã rất nhiều, chẳng hạn như mô tả về một phần hàng hóa trong WMS thương mại điện tử, có thể có các trường sau:

> Warehouse id, warehousing time, location id, storage shelf id, warehousing operator id, outbound operator id, inventory quantity, expiration time, SKU type, product brand, product category, number of internals

Ngoài các thông tin trên, nếu hàng hóa được lưu thông trong kho. Chúng có thể có id quá trình liên quan, trạng thái luồng hiện tại, v.v.

Hãy tưởng tượng nếu chúng ta đang điều hành một công ty thương mại điện tử lớn với hàng chục triệu đơn hàng mỗi ngày, sẽ rất khó để truy vấn và xây dựng index thích hợp trong cơ sở dữ liệu này.

In CRM or customer service systems, there is often a need to search by keyword, and large Internet companies receive tens of thousands of user complaints every day. Considering the source of the incident, the user's complaint must be at least 2 to 3 years. It is also tens of millions or even hundreds of millions of data. Performing a like query based on the keyword may hang the entire MySQL directly.
Trong CRM hoặc hệ thống dịch vụ khách hàng, ta thường có nhu cầu tìm kiếm theo từ khóa, và các công ty Internet lớn nhận được hàng chục ngàn khiếu nại của người dùng mỗi ngày. Xem xét nguồn gốc của vụ việc, khiếu nại của người dùng phải có ít nhất 2 đến 3 năm. Có đến hàng chục triệu hoặc thậm chí hàng trăm triệu dữ liệu. Thực hiện một truy vấn "like" dựa trên từ khóa có thể trực tiếp làm treo toàn bộ MySQL.

Lúc này, chúng tôi cần một công cụ tìm kiếm để cứu lấy trò chơi trên.

## Công cụ tìm kiếm

Elasticsearch is the leader of the open source distributed search engine, which relies on the Lucene implementation and has made many optimizations in deployment and operation and maintenance. Building a distributed search engine today is much easier than the Sphinx era. Simply configure the client IP and port.
Elaticsearch(ES) là đi đầu trong các công cụ tìm kiếm phân tán nguồn mở, dựa trên việc hiện thực Lucene và có nhiều tối ưu hóa trong triển khai, vận hành và bảo trì. Xây dựng một công cụ tìm kiếm phân tán ngày nay dễ dàng hơn nhiều so với thời đại Sphinx. Đơn giản chỉ cần cấu hình IP và cổng cho client.

### Inverted list

Although es is customized for the search scenario, as mentioned earlier, es is often used as a database in practical applications because of the nature of the inverted list. You can understand the inverted index with a simpler perspective:
Mặc dù ES được thiết kế cho việc tìm kiếm, nhưng ES thường được sử dụng làm cơ sở dữ liệu trong các ứng dụng vì bản chất của inverted list. Có thể dễ dàng hiểu được inverted index như sau:

![posting-list](../images/ch6-posting_list.png)

*Figure 6-10 Inverted list*

When querying data in Elasticsearch, the essence is to find a sequence of multiple ordered sequences. Non-numeric type fields involve word segmentation problems. In most internal usage scenarios, we can use the default bi-gram word segmentation directly. What is a bi-gram participle:
Khi truy vấn dữ liệu trong Elaticsearch, điều quan trọng là tìm một chuỗi nhiều chuỗi theo thứ tự. Các trường loại không số liên quan đến các vấn đề phân đoạn từ. Trong hầu hết các kịch bản sử dụng nội bộ, chúng ta có thể sử dụng phân đoạn từ bi-gram mặc định trực tiếp. Một phân từ bi-gram là gì:

Putting all `Ti` and `T(i+1)` into one word (called term in Elasticsearch), and then rearranging its inverted list, so our inverted list is probably like this:
Đặt tất cả `Ti` và`T(i+1)` vào một từ (được gọi là term trong Elaticsearch), và sau đó sắp xếp lại danh sách đảo ngược của nó, vì vậy danh sách đảo ngược của chúng ta sẽ như thế này:

![terms](../images/ch6-terms.png)

*Figure 6-11 "今天天气很好" word segmentation results*

When the user searches for "the weather is good", it is actually seeking: the weather, the gas is very good, the intersection of the three groups of inverted lists, but the equality judgment logic here is somewhat special, with pseudo code:
Khi người dùng tìm kiếm "天气很好", thực tế họ đang tìm kiếm: 天气, 气很, 很好 giao điểm của ba nhóm danh sách đảo ngược, nhưng logic phán đoán bình đẳng ở đây có phần đặc biệt, với mã giả:
  
```go
func equal() {
  if postEntry.docID of '天气' == postEntry.docID of '气很' &&
    postEntry.offset + 1 of '天气' == postEntry.offset of '气很' {
      return true
  }

  if postEntry.docID of '气很' == postEntry.docID of '很好' &&
    postEntry.offset + 1 of '气很' == postEntry.offset of '很好' {
    return true
  }

  if postEntry.docID of '天气' == postEntry.docID of '很好' &&
    postEntry.offset + 2 of '天气' == postEntry.offset of '很好' {
    return true
  }

  return false
}
```

The time complexity of multiple ordered lists to find intersections is: `O(N * M)`, where N is the smallest set of elements in a given list, and M is the number of given lists.
Độ phức tạp thời gian của nhiều danh sách được sắp xếp để tìm giao điểm là: `O(N * M)`, trong đó N là tập hợp phần tử nhỏ nhất trong danh sách đã cho và M là số lượng danh sách đã cho.

One of the decisive factors in the whole algorithm is the length of the shortest inverted list, followed by the sum of words, and the general number of words is not very large (Imagine, would you type hundreds of words in the search engine to search?) , so the decisive role is generally the length of the shortest one of all the inverted list.

Therefore, in the case where the total number of documents is large, the search speed is also very fast when the shortest one of the inverted lists of search terms is not long. If you use a relational database, you need to scan slowly by index (if any).

### Query DSL

Es defines a set of query DSL, when we use es as a database, we need to use its bool query. for example:

ES định nghĩa một tập hợp DSL truy vấn, khi chúng ta sử dụng ES làm cơ sở dữ liệu, chúng ta cần sử dụng truy vấn bool của nó. ví dụ:

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

Trông có vẻ rườm rà, nhưng ý nghĩa của biểu thức rất đơn giản:

```go
If field_1 == 1 && field_2 == 2 && field_3 == 3 && field_4 == 4 {
    Return true
}
```

Use bool should query to represent the logic of or:
Sử dụng bool "should" để truy vấn có logic "or":

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

Giống với biểu thức dưới đây:

```go
If field_1 == 1 || field_2 == 2 {
  Return true
}
```

The expressions followed by `if` in these Go codes have proper nouns in the programming language to express `Boolean Expression`:
Các biểu thức được theo sau `if` trong code Go ở trện có các danh từ riêng trong ngôn ngữ lập trình để diễn tả `Boolean Expression`:

```go
4 > 1
5 == 2
3 < i && x > 10
```

The `Bool Query` scheme of es actually uses json to express the Boolean Expression in this programming language. Why can this be done? Because json itself can express the tree structure, our program code will become AST after being parse by the compiler, and the AST abstract syntax tree, as its name implies, is a tree structure. In theory, json can fully express the result of a piece of program code being parse. The Boolean Expression here is also generated by the compiler Parse and generates a similar tree structure, and is only a small subset of the entire compiler implementation.
Lược đồ `Bool Query` của ES sử dụng json để diễn tả Boolean Expression trong ngôn ngữ lập trình này. Tại sao lại hiện thực bằng cách này? Vì chính json có thể biểu thị cấu trúc cây, code chương trình của chúng ta sẽ trở thành AST (abstract syntax tree) sau khi được trình biên dịch phân tích cú pháp. Về lý thuyết, json hoàn toàn có thể diễn tả kết quả của một đoạn mã chương trình được phân tích cú pháp. Boolean Expression ở đây cũng được tạo bởi trình biên dịch Parse và tạo ra một cấu trúc tương tự cây, và là một phần nhỏ của cách hiện thực trình biên dịch.

### Phát triển SDK cho người dùng

Khởi tạo:

```go
// When using the elastic version
// Note that it corresponds to the elasticsearch you use.
Import (
  Elastic "gopkg.in/olivere/elastic.v3"
)

Var esClient *elastic.Client

Func initElasticsearchClient(host string, port string) {
  Var err error
  esClient, err = elastic.NewClient(
    elastic.SetURL(fmt.Sprintf("http://%s:%s", host, port)),
    elastic.SetMaxRetries(3),
  )

  If err != nil {
    // log error
  }
}
```

Chèn:

```go
Func insertDocument(db string, table string, obj map[string]interface{}) {

  Id := obj["id"]

  Var indexName, typeName string
  // The database/table concept in the database can be simply mapped to the index and type of es
  // but note, because the _type in es is essentially just a field of document
  // So too much single index content can cause performance issues
  // type is deprecated in the new version
  / / In order to make the data of different tables fall into different indexes, here we use table + name as the name of index
  indexName = fmt.Sprintf("%v_%v", db, table)
  typeName = table

  // normal situation
  Res, err := esClient.Index().Index(indexName).Type(typeName).Id(id).BodyJson(obj).Do()
  If err != nil {
    // handle error
  } else {
    // insert success
  }
}
```

Đạt được:

```go
Func query(indexName string, typeName string) (*elastic.SearchResult, error) {
  // Add bool query conditions via bool must and bool should
  q := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("id", 1),
  elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("male", "m")))

  q = q.Should(
    elastic.NewMatchPhraseQuery("name", "alex"),
    elastic.NewMatchPhraseQuery("name", "xargin"),
  )

  searchService := esClient.Search(indexName).Type(typeName)
  Res, err := searchService.Query(q).Do()
  If err != nil {
    // log error
    Return nil, err
  }

  Return res, nil
}
```

Xoá:

```go
Func deleteDocument(
  indexName string, typeName string, obj map[string]interface{},
) {
  Id := obj["id"]

  Res, err := esClient.Delete().Index(indexName).Type(typeName).Id(id).Do()
  If err != nil {
    // handle error
  } else {
    // delete success
  }
}
```

Because of the nature of Lucene, the data in the search engine is essentially immutable, so if you want to update the document, it is actually completely covered by id, so it is the same as the insertion.
Do bản chất của Lucene, dữ liệu trong công cụ tìm kiếm là bất biến. Vì vậy, nếu bạn muốn cập nhật tài liệu, thì nó cũng là một lệnh chèn.

When using es as a database, you need to be aware that because es has an operation of index merging, it takes a while for the data to be inserted into es to be queried (determined by ref's refresh_interval). So don't use es as a strong and consistent relational database.
Khi sử dụng ES làm cơ sở dữ liệu, bạn cần lưu ý rằng ES có hoạt động hợp nhất index, cần một khoản thời gian để dữ liệu được chèn vào ES và có thể truy vấn được (được xác định bởi refresh_interval). Vì vậy, không sử dụng es như một cơ sở dữ liệu quan hệ "strong consistent".

### Chuyển đổi sql sang DSL

For example, we have a bool expression, `user_id = 1 and (product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1)`, written in SQL as follows:
Ví dụ: chúng ta có một biểu thức bool, `user_id = 1 and (product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1)`, được viết bằng SQL như sau:

```sql
Select * from xxx where user_id = 1 and (
  Product_id = 1 and (star_num = 4 or star_num = 5) and banned = 1
)
```

DSL của ES được viết như sau:

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

Although es DSL is well understood, it is very hard to write. The SDK-based approach was provided earlier, but it is not flexible enough.
Mặc dù ES DSL khá dễ hiểu, nhưng nó rất khó để viết. Cách sử dụng dựa vào SDK đã được cung cấp ở trên, nhưng nó không đủ linh hoạt.

The where part of SQL is the boolean expression. As we mentioned before, this bool expression is similar to the es DSL structure after being parsed. Can we directly help us convert SQL to DSL through this "almost" guess?
Phần "where" của SQL là biểu thức boolean. Như chúng ta đã biết, biểu thức bool này tương tự như cấu trúc ES DSL sau khi được "parse". Với phỏng đoán này, liệu chúng ta có thể trực tiếp chuyển đổi SQL sang DSL được không?

Of course, we can compare the structure of SQL with the structure of Parse and the structure of es DSL:
Câu trả lời là tất nhiên được, so sánh cấu trúc của SQL sau khi Parse và cấu trúc của ES DSL:

![ast](../images/ch6-ast-dsl.png)

*Figure 6-12 Correspondence between AST and DSL*

Since the structure is completely consistent, we can logically convert each other. We traverse the AST tree with breadth first, then convert the binary expression into a json string, and then assemble it. Due to space limitations, no examples are given in this article. Readers can check:
Vì cấu trúc hoàn toàn giống nhau, nên chúng ta có thể chuyển đổi chúng. Chúng ta duyệt cây AST theo chiều ngang trước, sau đó chuyển đổi biểu thức nhị phân thành chuỗi json, và tổng hợp nó lại. Về chi tiết khá phức tạp nên không có ví dụ cụ thể trong bài viết này. Các bạn có thể tham khảo chi tiết ở đây:

> github.com/cch123/elasticsql

To learn the specific implementation.
Để biết cách hiện thực cụ thể.

## Đồng bộ dữ liệu không đồng nhất

In practical applications, we rarely write data directly to search engines. A more common way is to synchronize data from MySQL or other relational data into a search engine. The user of the search engine can only query the data and cannot modify and delete it.

There are two common synchronization schemes:

Trong thực tế, chúng ta hiếm khi ghi dữ liệu trực tiếp vào công cụ tìm kiếm. Một cách phổ biến hơn là đồng bộ hóa dữ liệu từ MySQL hoặc dữ liệu quan hệ khác vào công cụ tìm kiếm. Người dùng của công cụ tìm kiếm chỉ có thể truy vấn dữ liệu và không thể sửa đổi hoặc xóa nó.

Có hai sơ đồ đồng bộ hóa phổ biến:

### Đồng bộ hóa dữ liệu tăng dần bằng timestamp

![sync to es](../images/ch6-sync.png)

*Figure 6-13 Time-based data synchronization*

This synchronization method is strongly bound to the business, such as the outbound order in the WMS system, we do not need very real time, a little delay is acceptable, then we can get the last ten from the MySQL outbound order table every minute. All the outbound orders created in minutes are taken out and stored in es in batches. The specific logic is actually a SQL:
Phương thức đồng bộ hóa này bị ràng buộc chặt chẽ với nhu cầu của doanh nghiệp, chẳng hạn như đơn hàng bên ngoài trong hệ thống WMS, chúng ta không cần xử lý "real time", chậm một chút có thể chấp nhận được, lúc này chúng ta có thể lấy mười giá trị mới nhất mỗi phút từ bảng đặt hàng của MySQL. Tất cả các đơn đặt hàng được tạo trong vài phút sẽ được lấy ra và lưu trữ trong ES theo "batch". Logic cụ thể như câu SQL dưới đây:

```sql
Select * from wms_orders where update_time >= date_sub(now(), interval 10 minute);
```

Of course, considering the boundary situation, we can make the data of this time period overlap with the previous one:
Tất nhiên, khi xem xét các tình huống xảy ra ở giữa 2 khoảng thời gian, chúng ta nên lấy dữ liệu với khoảng thời gian có một ít trùng lặp với khoảng thời gian trước đó:

```sql
Select * from wms_orders where update_time >= date_sub(
  Now(), interval 11 minute
);
```

Update the data coverage changed to es in the last 11 minutes. The shortcomings of this approach are obvious, and we must require business data to strictly adhere to certain specifications. For example, there must be an update_time field, and each time it is created and updated, the field must have the correct time value. Otherwise our synchronization logic will lose data.

Cập nhật phạm vi bảo hiểm dữ liệu đã thay đổi thành ES trong 11 phút cuối. Những thiếu sót của phương pháp này là rõ ràng và chúng tôi phải yêu cầu dữ liệu kinh doanh tuân thủ nghiêm ngặt các thông số kỹ thuật nhất định. Ví dụ: phải có trường update_time và mỗi lần được tạo và cập nhật, trường phải có giá trị thời gian chính xác. Nếu không, việc đồng bộ hóa sẽ mất dữ liệu.

### Synchronizing data with binlog

![binlog-sync](../images/ch6-binlog-sync.png)

*Figure 6-13 Data synchronization based on binlog*

The industry uses more of Ali's open source Canal for binlog parsing and synchronization. Canal will pretend to be a MySQL slave library, then parse the bincode of the line format and send it to the message queue in a more easily parsed format (such as json).

The downstream Kafka consumer is responsible for writing the self-incrementing primary key of the upstream data table as the id of the es document, so that each time the binlog is received, the corresponding id data is updated to the latest. MySQL's Row format binlog will provide all the fields of each record to the downstream, so in fact, when synchronizing data to heterogeneous data targets, you don't need to consider whether the data is inserted or updated, as long as you always cover by id.

This model also requires the business to adhere to a data table specification, that is, the table must have a unique primary key id to ensure that the data we enter into es will not be duplicated. Once the specification is not followed, it will result in data duplication when synchronizing. Of course, you can also customize the consumer logic for each required table, which is not the scope of the general system discussion.
