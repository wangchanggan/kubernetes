package main

import (
	"github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
)

func main() {
	// 在main函数中未实例化Container 操作，默认会使用restful.defaultContainer和DefaultServeMux，可以通过rsful.NewContainer函数来创建一个Container。

	// 实例化restful.WebService，为WebService添加一个Router。在Router中定义请求方法(GET方法)、请求路径(hello)、 Handler 函数(hello 函数)。
	// 其中请求方法接收Request和Response，并与用户数据进行交互。
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	// 通过Add方法将Router 添加到WebService中。
	restful.Add(ws)
	// 最后通过Go语言标准库http.ListenAndServer监控地址和端口。
	log.Fatal(http.ListenAndServe(":8080", nil))

	// 运行以上程序，监控8080端口。当我们发送GET请求到http://localhost:8080/hello时，会得到响应“world”。
	// 该程序提供了HTTP短连接(非持久连接)服务请求。客
	//户端和服务端在进行一次HTTP请求/响应之后，会关闭连接，下一次的HTTP请求/ 响应操作需要重新建立连接。
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world\n")
}
