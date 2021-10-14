# Integration test use document
-1. The goconvey command must be installed to ensure that the command can be executed
```
cd $GOPATH
go get github.com/smartystreets/goconvey
goconvey
```

-2. Create a test database and import the data structure
```
mysql -h 127.0.0.1 -u root -p -e "CREATE DATABASE gatekeeper_test DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h 127.0.0.1 -u root -p gatekeeper_test <gatekeeper_test.sql --default-character-set=utf8
```

-3. The configuration of conf/dev must be rationalized

To modify mysql, please edit conf/dev/mysql_map.toml
To modify redis, please edit conf/dev/redis_map.toml

-4. Use go test to test

Execute ```go test -v``` to complete the use case test

Or use interface testing tools
Execute the ```sh goconvey.sh``` command, it will automatically open http://127.0.0.1:8080/ under normal circumstances
Observe the interface to view the test output.


# 集成测试使用文档

- 1、必须安装goconvey命令，确保命令可以执行
```
cd $GOPATH
go get github.com/smartystreets/goconvey
goconvey
```

- 2、创建测试数据库并导入数据结构
```
mysql -h 127.0.0.1 -u root -p -e "CREATE DATABASE gatekeeper_test DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;"
mysql -h 127.0.0.1 -u root -p gatekeeper_test < gatekeeper_test.sql --default-character-set=utf8
```

- 3、必须合理化配置 conf/dev

修改mysql，请编辑 conf/dev/mysql_map.toml
修改redis，请编辑 conf/dev/redis_map.toml

- 4、使用 go test 测试

执行 ```go test -v``` 完成用例测试

或者使用界面化测试工具
执行 ```sh goconvey.sh``` 命令 ，正常情况下会自动打开 http://127.0.0.1:8080/   
观察界面查看测试输出即可。