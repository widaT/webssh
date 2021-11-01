# WEBSSH

基于vue、xterm、golang实现的web ssh客户端程序

## 特性
- 前后端分离，前端使用xterm、vue，后端使用golang写的服务
- 支持录像审计，支持录像回看

## run demo

- 编译前端程序
```bash
$ cd front
$ npm -i
$ npm run build # 可以看到在front线生产一个dist目录，里头就是编译后的前端文件
```
- 编译golang程序

修改`main.go`文件中目标主机和登录方式

```go
confing := &webssh.WebSSHConfig{
		Record:     true,
		RecPath:    "./rec/cast/",
		RemoteAddr: "localhost:22",
		User:       "wida",
		Password:   "wida",
		AuthModel:  webssh.PASSWORD,
	}
```

```bash
$ go build -o webssh main.go
$ ./webssh
```
- 用浏览器打开`http://localhost:8080/`

## 查看录像

修改 `rec/index.html`的`src`

```html
   <asciinema-player src="./1.cast"></asciinema-player>
```

- 用浏览器打开`http://localhost:8080/rec/`