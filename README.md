# WEBSSH

基于vue、xterm、golang实现的web ssh客户端程序

## 特性
- 前后端分离，前端使用xterm、vue，后端使用golang写的服务
- 

## run demo

- 编译前端程序
```bash
$ cd front
$ npm -i
$ npm run build # 可以看到在front线生产一个dist目录，里头就是编译后的前端文件
```
- 编译golang程序

修改`handle.go`文件中目标主机和登录方式

```go
client, err := NewSSHClient(
    //使用私钥登录
    /* 		SSHClientConfigPulicKey(
        "host:22",
        "user",
        "/home/user/ssh/user.id_rsa",
    ), */
    //使用密码登录
    SSHClientConfigPassword(
        "host:22",
        "user",
        "pwd",
    ),
)
```

```bash
$ go build -o webssh main.go
$ ./webssh
```
- 用浏览器打开`http://localhost:8080/`

```bash
$ asciinema play rec/filename
```