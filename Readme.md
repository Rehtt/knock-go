# Knock-go

一个简单使用UDP广播发现局域网上的Hostname的程序

## install

go >= 1.21

```shell
go install github.com/Rehtt/knock-go@latest
```

## use

### server
```shell
knock-go -s
```
```shell
2023/09/08 15:55:46 INFO running
2023/09/08 15:55:48 INFO knock addr=192.168.56.1:56142
```

### client
```shell
knock-go
```
```shell
IP               Hostname
192.168.14.1     PC
```