# Integration test use document

- 1. The goconvey command must be installed to ensure that the command can be executed
   cd $GOPATH
   go get github.com/smartystreets/goconvey
   goconvey

- 2. The configuration of conf/dev must be rationalized
   To modify mysql, please edit conf/dev/mysql_map.toml
   To modify redis, please edit conf/dev/redis_map.toml

- 3. Start goconvey
   Execute sh goconvey.sh command, it will automatically open http://127.0.0.1:8080/ under normal circumstances

# 集成测试使用文档

- 1、必须安装goconvey命令，确保命令可以执行
cd $GOPATH
go get github.com/smartystreets/goconvey
goconvey

- 2、必须合理化配置 conf/dev
修改mysql，请编辑 conf/dev/mysql_map.toml
修改redis，请编辑 conf/dev/redis_map.toml

- 3、启动 goconvey
执行 sh goconvey.sh 命令 ，正常情况下会自动打开 http://127.0.0.1:8080/