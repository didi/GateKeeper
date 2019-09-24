tester存放的是集成测试的内容

目录

# suites

存放测试套件，每个测试套件

## suites/xxx

这里存放测试套件，测试套件文件夹需要包含before.go和after.go两个文件

before.go存放有

* SetUp() 函数，这个函数在Suite运行之前会运行
* Before() 函数，这个函数在所有Case运行之前运行

after.go存放有

* TearDown() 函数，这个函数在Suite运行之后会运行
* After() 函数，这个函数在Suite运行之后运行

# conf

测试环境的配置

# report

存放报告的地址

代码覆盖率需要额外跑脚本

在tester目录下运行：

`sh bootstrap.sh` 会自动打开浏览器，并执行自动化测试
`sh coverage.sh` 会在report下生成coverage.out和coverage.html，并自动打开浏览器