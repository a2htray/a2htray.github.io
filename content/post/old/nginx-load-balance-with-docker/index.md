---
title: "基于 Docker 的 Nginx 负载均衡配置"
date: 2023-02-16T21:22:17+08:00
draft: false
reward: false
tags:
 - Nginx
 - Docker
 - Docker Compose
 - Load Balance
categories:
  - 后端技术
  - Nginx
---

本文主要介绍如何基于 Docker 生成容器化应用以及配置 Nginx 来实现 Web 服务的负载均衡，包含以下内容：

1. Go 简易 Web 程序（hello world）
2. 通过 Dockerfile 生成镜像
3. Nginx 配置
4. Docker Compose 多容器启动

## Go Web 程序

本次配置使用 Go Web 程序作为后端服务，其作用是接收请求并返回“hello-world”，代码如下：

```go
// main.go
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		logFile.WriteString("request arrived\n")
		w.Write([]byte("hello world\n"))
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}

```

`/var/log/app.log` 是程序的日志文件，用于标识服务器端接收到了请求，作为后续分析流量的依据。执行以下命令查看运行结果：

```bash
$ go run main.go
# 另启一个终端并执行
$ curl localhost:8000/hello-world
hello-world
$ cat /var/log/app.log
request arrived
```

从结果输出中可以看出，程序正常运行并成功记录请求日志。

## Web 服务镜像

Web 服务镜像就是将上一节的 Go Web 服务封装到镜像中，基于该镜像可以生成容器并提供服务。在与 `main.go` 文件同级的目录下创建 `Dockerfile` 文件，内容如下：

```dockerfile
FROM golang:1.19.3

WORKDIR /usr/src/app

COPY main.go ./

RUN go build -o /usr/local/bin/app ./main.go

EXPOSE 8000

CMD ["app"]
```

各行作用为：

1. FROM xxx：以 golfing:1.19.3 作为基础镜像，后续修改在基础镜像上修改
2. WORKDIR xxx：设置工作目录，为 COPY、RUN 指令指定当前目录
3. COPY xxx：复制 main.go 到当前目录
4. RUN xxx：生成二进制可执行文件
5. EXPOSE xxx：对外暴露 8000 端口
6. CMD xxx：容器运行时要执行的命令

下面通过 `Dockerfile` 生成一个镜像并运行容器：

```bash
$ docker build -t goweb:0.1 .
$ docker run -d -p 8000:8000 goweb:0.1
$ curl localhost:8000/hello-world
hello-world
```

## Nginx 配置

此次配置中只涉及到两个 Nginx 配置指令，分别是 proxy_pass 和 upstream。

### proxy_pass

proxy_pass，即代理转发指令，后面跟要转发到的地址，需要在 location 块中声明，其格式为 `proxy_pass: URL`，如：

```nginx
server {
  # 其它配置 ...
  location / {
    proxy_pass http://localhost:9001/;
  }
}
```

上述配置会将所有请求转发到 http://localhost:9001/。

### upstream

如果希望将请求转发到一组后端服务，也就是要实现负载均衡，则需要配合 upstream 指令，如：

```nginx
upstream backend {
  server localhost:9001;
  server localhost:9002;
  server localhost:9003;
}

server {
  location / {
    proxy_pass http://backend/;
  }
}
```

如上，upstream 指令定义一组后端服务，backend 作为组的名称供后面配置中引用，Nginx 默认的负载均衡机制为轮询，也就是请求会按序被引流到各个后端服务。

## Docker Compose

Docker Compose 是 Docker 提供的容器编排工具，通过声明 yml 配置文件可以同时运行多个容器，并且各容器间可建立一定的联系。以下是本次试验中的节点拓扑图：

```bash
request --> Nginx --> rp-app1
                  --> rp-app2
                  --> rp-app3
```

请求到达 Nginx 服务器后会依照轮询机制转发到 3 个后端服务。因为在 docker-compose.yml 文件中使用到了相对路径，因此需要提前说明下当前的目录结构，如下：

```bash
- go1.19.3
 Dockerfile
 main.go
- nginx
 default.conf
docker-compose.yml
```

相应的 docker-compose.yml 内容如下：

```yaml
version: "3.9"
services:
  app1:
    build: ./go1.19.3
    container_name: rp-app1
  app2:
    build: ./go1.19.3
    container_name: rp-app2
  app3:
    build: ./go1.19.3
    container_name: rp-app3
  nginx:
    image: "nginx"
    container_name: rp-nginx
    ports:
      - "8000:80"
    volumes:
      - ./nginx/:/etc/nginx/conf.d/
    links:
      - app1
      - app2
      - app3
```

运行以下命令启动容器：

```bash
$ docker-compose up -d
$ docker ps -a
CONTAINER ID   IMAGE                     COMMAND                  CREATED         STATUS                     PORTS                  NAMES
7df4d168608f   nginx                     "/docker-entrypoint.…"   6 seconds ago   Up 5 seconds               0.0.0.0:8000->80/tcp   rp-nginx
1bac1946f4d6   reverseproxy_app2         "app"                    7 seconds ago   Up 5 seconds               8000/tcp               rp-app2
cbaee01c7d80   reverseproxy_app1         "app"                    7 seconds ago   Up 6 seconds               8000/tcp               rp-app1
0533a2cf39d1   reverseproxy_app3         "app"                    7 seconds ago   Up 6 seconds               8000/tcp               rp-app3
```

## 流量结果

容器化应用正常启动后，用 `ab` 命令行进行并发请求，然后查看日志分析各后端服务接收到的请求数。

```bash
$ ab -n 1000 -c 100 localhost:8000/hello-world
$ docker exec -it rp-app1 wc -l /var/log/app.log
335 /var/log/app.log
$ docker exec -it rp-app2 wc -l /var/log/app.log
334 /var/log/app.log
$ docker exec -it rp-app2 wc -l /var/log/app.log
331 /var/log/app.log
```

可以看到，rp-app1、rp-app2 和 rp-app3 分别接收到 335、334 和 331 个请求，共 1000 个请求。

## 总结

本文完成了以下工作：

1. 基于 Go 实现简易的 Web 程序作为后端服务（`net/http`）
2. 通过 Dockerfile 将后端服务打包成镜像（`Dockerfile 指令`）
3. 介绍了 Nginx 中用于负载均衡的两个指令（`proxy_pass`、`upstream`）
4. 借助 Docker Compose 运行多个容器（`docker-compose up -d`）
5. 完成了流量分析（`ab`、`docker exec`）

全部代码及配置见：https://github.com/a2htray/code-notebook/tree/main/Nginx/ReverseProxy。