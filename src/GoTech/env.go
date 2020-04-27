package GoTech

// Go先是从$GOROOT中查找包myfunc，如果没找到就从$GOPATH中查找，结果都没有找到，我们可以使用go env输出Go的环境变量设置

/*

go build 编译包，如果是main包则在当前目录生成可执行文件，其他包不会生成.a文件；
go install 编译包，同时复制结果到$GOPATH/bin，$GOPATH/pkg等对应目录下；
	1. 对单个文件使用go install，就会出现这个错误,应该用 go install hello 编译
	就会编译出 bin/hello 文件,之后用 ./bin/hello 运行
	2. 只有当环境变量GOPATH中只包含一个工作区的目录路径时，go install命令才会把命令源码安装到当前工作区的bin目录下
go run gofiles... 编译列出的文件，并生成可执行文件然后执行。注意只能用于main包，否则会出现go run: cannot run non-main package的错误。
	go run是不需要设置$GOPATH的，但go build和go install必须设置。go run常用来测试一些功能，这些代码一般不包含在最终的项目中。

go get
在使用go的时候如果依赖导入github上的，比如下面样式
import "github.com/go-sql-driver/mysql"
我们需要先执行get操作
go get github.com/go-sql-driver/mysql
它会下载到你的gopath目录下

go mod
使用go mod ，利用Go 的 module 特性，你再也不需要关心GOPATH了（当然GOPATH变量还是要存在的，但只需要指定一个目录，而且以后就不用我们关心了），
你可以任性的在你的硬盘任何位置新建一个Golang项目了。

因为这两个包不在同一个项目路径下，你想要导入本地包，并且这些包也没有发布到远程的github或其他代码仓库地址。
这个时候我们就需要在go.mod文件中使用replace指令。



*/
