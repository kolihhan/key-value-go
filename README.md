## Storage 選擇 - PostgreSQL

本身之前較常使用的是 Mysql，選擇PostgreSQL的原因有:

- 數據類型: MySQL 和 PostgreSQL 支持的數據類型有所不同，PostgreSQL 支持更多的數據類型
- 性能：MySQL 更適合處理大量的簡單查詢和事務，而 PostgreSQL 更適合處理複雜的查詢和高並發的環境。

## UnitTest 測試
```sh
go test -run TestGetHeadHandler
go test -run TestGetPageHandler
go test -run TestSetHandler
```
