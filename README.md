#webhook ci项目

主要用于小项目托管在[码云](http://gitee.com)和[Coding](http://coding.net)等平台的项目进行自动部署
可多项目部署，如有多个项目，增加多个 `.json` 文件配置即可

# 支持

- [x] 支持Windows/Linux平台下部署
- [x] 支持Gogs
- [x] 支持Coding
- [x] 支持Gitee

# 命令列表

```shell
PROG -p 7442     # 启动监听端口为7442的服务器
hook -h          # 帮助信息
```

# 路由列表

* /
自动解析，将根据Header头中的参数自动解析

* /gitee
解析[码云](http://gitee.com)的webhook通知，支持 `ContentType: application/json` 

* /coding
解析[Coding](http://coding.net)的webhook通知，支持V2版本的通知、`ContentType: application/json`

* /gogs
解析Gogs系列的webhook通知，支持`ContentType: application/json`
