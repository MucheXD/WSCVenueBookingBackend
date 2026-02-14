# configs - sql

存储数据库操作文件，便于后期修改。

开发阶段需要更改表结构，建议直接重新建表并更新此处操作文件。

Comment 中标注 "DBOnly" 的为用于数据库内索引的字段（一般仅用于主键），不使用；

Comment 中标注 "RepoLayerOnly" 的为 Repository 层专用字段，不提供给业务层；