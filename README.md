## 2023-Dcard-Intern

實作一個共用的 Key-Value 列表系統，用於 Dcard 的個人化推薦文章列。
| 工具 | 版本 |
| ------ | ------ |
| Go | 1.20 |
| PostgreSQL | 15.2 |
功能如下:
1.  `Golang` 開發並設計 Restful API 。
2. 執行 `go test` 即可執行 Integration Test
3. 避免 storage 儲存太多，需要定期清除不用的內容。




## Storage 選擇 - PostgreSQL

本身之前較常使用的是 Mysql，選擇PostgreSQL的原因有: 
- 数据类型: MySQL 和 PostgreSQL 支持的数据类型有所不同，PostgreSQL 支持更多的数据类型
- 性能：MySQL 更适合处理大量的简单查询和事务，而 PostgreSQL 更适合处理复杂的查询和高并发的环境。


## 定期清除不用的內容

使用的是PostgreSQL 的 Time to Live 功能
- TTL 功能可以讓 PostgreSQL 自動刪除超過一定時間限制的資料。
- 需注意的是PostgreSQL 9.6以上才可以使用



## UnitTest 測試

```sh
go test
go test -run TestGetHeadHandler
go test -run TestGetPageHandler
go test -run TestSetHandler
```

![Test Result](https://i.postimg.cc/yNCVrxjL/Capture1.png)