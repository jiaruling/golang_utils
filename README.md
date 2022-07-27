# Alpine Linux 系统下开发注意事项

## 安装编译环境

``` shell
$ apk add gcc
$ apk add musl-dev
```

## 修改时区

``` shell
# url: https://blog.csdn.net/isea533/article/details/87261764

# 安装时区设置
$ apk add tzdata
# 复制上海时区
$ cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 指定为上海时区
$ echo "Asia/Shanghai" > /etc/timezone
# 验证时区
$ date -R
Fri, 22 Jul 2022 17:07:46 +0800
# 删除其它时区配置, 节省空间 【可选】
$ apk del tzdata
```

# golang 开发常用工具包功能介绍

## 链路日志

## http client

## Tips

## 数据操作

### MySQL

### Redis

