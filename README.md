# CDN Proxy

## 项目功能

CDN Proxy 是一个使用 Go 开发的，用于代理公共静态资源 CDN 的工具。

我的基于 Hexo 的站点一直在使用公共 CDN 来加载一些库，比如：

- animejs
- medium-zoom
- moment

我先后使用过 [baomitu](https://cdn.baomitu.com)、[Staticfile](https://www.staticfile.net)，但是都出现了服务不能稳定使用的问题。最终还是觉得 [jsdelivr](https://www.jsdelivr.com) 的服务比较稳定，但是直接访问 jsdelivr 的速度太慢了，所以有了这个工具来帮助代理提速。

具有以下特性：

1. 支持本地缓存资源文件
2. 支持 Socks5 代理加速

## 部署方式

可使用 Docker Compose 来部署：

```yaml
services:
  cdn-proxy:
    image: ixuewen/cdn-proxy
    container_name: cdn-proxy
    restart: always
    volumes:
      # 提供配置文件到容器中，可替换为你自己的本地路径
      - ./cdn-proxy/config.yaml:/app/config.yaml
      # 持久化缓存目录，可选
      - ./cdn-proxy/cache:/app/cache
```

配置文件说明如下：

```yaml
# 端口号，可根据实际情况修改
port: 80
# 允许跨域的 origin，多个 origin 支持使用英文逗号分隔
cors_origin: https://a.com,https://b.com
# 缓存文件的保存目录
cache_dir: /app/cache
# 访问源 CDN 时的代理，可选
socks_proxy_url: socks5://127.0.0.1:1080
# 支持的 CDN 列表，可以配置多个，会根据请求时的 CDN_NAME 自动选择
cdn:
  - name: jsdelivr
    url: https://cdn.jsdelivr.net
```

## 使用方式

假设部署在 cdn.abc.com 域名下，则接受的请求格式为：

```
https://cdn.abc.com/{CDN_NAME}/{PATH_TO_RESOURCE}
```

例如：

| 原始链接                                                     | 代理链接                                                         |
| ------------------------------------------------------------ | ---------------------------------------------------------------- |
| https://cdn.jsdelivr.net/npm/animejs@3.2.1/lib/anime.min.js  | https://cdn.abc.com/jsdelivr/npm/animejs@3.2.1/lib/anime.min.js  |
| https://cdn.jsdelivr.net/npm/moment@2.30.1/min/moment.min.js | https://cdn.abc.com/jsdelivr/npm/moment@2.30.1/min/moment.min.js |
