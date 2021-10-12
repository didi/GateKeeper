# Integration test use document
-1. The goconvey command must be installed to ensure that the command can be executed
```
cd $GOPATH
go get github.com/smartystreets/goconvey
goconvey
```

-2. conf/dev must be rationally configured

To modify mysql, please edit conf/dev/mysql_map.toml
To modify redis, please edit conf/dev/redis_map.toml

-3. Start the test server

```sh test_httpserver.sh``` will get the following output.
```
[INFO][2021-10-11T17:37:00.055+0800]||start test httpserver::8881
[INFO][2021-10-11T17:37:00.055+0800]||start test httpserver::8882
```

-4. Use goconvey to test

Execute the ```sh goconvey.sh``` command, it will automatically open http://127.0.0.1:8080/ under normal circumstances
Observe the interface to view the test output.

-5. Use go test to test

Execute ```go test -v``` to complete the use case test

# 集成测试使用文档

- 1、必须安装goconvey命令，确保命令可以执行
```
cd $GOPATH
go get github.com/smartystreets/goconvey
goconvey
```

- 2、必须合理化配置 conf/dev

修改mysql，请编辑 conf/dev/mysql_map.toml
修改redis，请编辑 conf/dev/redis_map.toml

- 3、启动测试服务器

```sh test_httpserver.sh``` 会得到如下输出。
```
[INFO][2021-10-11T17:37:00.055+0800]||start test httpserver::8881
[INFO][2021-10-11T17:37:00.055+0800]||start test httpserver::8882
```

- 4、使用 goconvey 测试

执行 ```sh goconvey.sh``` 命令 ，正常情况下会自动打开 http://127.0.0.1:8080/   
观察界面查看测试输出即可。

- 5、使用 go test 测试

执行 ```go test -v``` 完成用例测试