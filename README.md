## Storage 選擇 - PostgreSQL

本身之前較常使用的是 Mysql，選擇PostgreSQL的原因有: 
- 数据类型: MySQL 和 PostgreSQL 支持的数据类型有所不同，PostgreSQL 支持更多的数据类型
- 性能：MySQL 更适合处理大量的简单查询和事务，而 PostgreSQL 更适合处理复杂的查询和高并发的环境。

## UnitTest 測試
```sh
go test -run TestGetHeadHandler
go test -run TestGetPageHandler
go test -run TestSetHandler
```