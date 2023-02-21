[go-md-book](https://github.com/hedongshu/go-md-book)

> 这是一款可以快捷的把 md 文件发布成 web 页面 的程序
> ![light](https://github.com/hedongshu/go-markdown-book/blob/main/light.jpg?raw=true) > ![dark](https://github.com/hedongshu/go-markdown-book/blob/main/dark.jpg?raw=true)

## 支持平台

> Windows、Linux、Mac OS

## 安装

1. 下载 [release](https://github.com/hedongshu/go-md-book/releases)

2. 解压

```
tar zxf markdown-book-V1.0.0-darwin-amd64.tar.gz
```

3. 创建 markdown-book 文件目录

```
cd markdown-book-linux-amd64
mkdir md && mkdir md/技术 md/生活
echo "### Hello World" > .md/技术/主页.md
echo "### Hello World" > .md/生活/主页.md
```

4. 运行

```
./markdown-book web
```

5. 访问 http://127.0.0.1:5006，查看效果

## 升级

1. 下载最新版 [release](https://github.com/hedongshu/go-md-book/releases)

2. 停止程序，解压替换 `markdown-book`

3. 重新启动程序

## 使用

### 命令

- markdown-book
  - -h 查看版本
  - web 运行博客服务
- markdown-book web
  - --config FILE 加载配置文件,
  - --origin value, -o value 指定数据源，可选：github,file，默认:file
  - --dir value, -d value 指定 markdown 文件夹，默认：./md/
  - --title value, -t value web 服务标题，默认："Blog Title"
  - --title2 value, -t2 value web 副标题，默认："Blog Title2"
  - --port value, -p value web 服务端口，默认：5006
  - --env value, -e value 运行环境, 可选：dev,test,prod，默认："prod"
  - --cache value, -c value 设置页面缓存时间，单位分钟，默认 3 分钟
  - --github.owner 设置 GitHub 数据源的用户名
  - --github.repo 设置 GitHub 数据源的仓库名
  - --analyzer-baidu value 设置百度分析统计器
  - --analyzer-google value 设置谷歌分析统计器
  - --gitalk.client-id value 设置 Gitalk ClientId, 默认为空
  - --gitalk.client-secret value 设置 Gitalk ClientSecret, 默认为空
  - --gitalk.repo value 设置 Gitalk Repo, 默认为空
  - --gitalk.owner value 设置 Gitalk Owner, 默认为空
  - --gitalk.admin 设置 Gitalk Admin, 默认为数组 [gitalk.owner]
  - --gitalk.labels 设置 Gitalk Admin, 默认为数组 ["gitalk"]
  - --ignore-file value 设置忽略文件, eg: demo.md
  - --ignore-path value 设置忽略文件夹, eg: demo
  - -h 查看版本

### 运行参数

> 支持从配置文件读取配置项，不过运行时指定参数优先于配置文件，配置内容参考 `config/config.yml.tmp`

### 配置文件

1. 新建配置文件 `config/config.yml`

2. 启动时加载配置文件

- 二进制文件

```
./markdown-book web --config ./config/config.yml
```

### 设置数据源

> 支持使用 github 或者本地的文件夹作为文章的数据源，设置 origin 为 `github` or file`

#### 使用 GitHub 为数据源

```yaml
origin: github
github:
  owner: "github 用户名"
  repo: "github 仓库名"
```

假如你有一个 github 的仓库，里面存了你的 markdown 文档，你只需要把用户名和仓库名填到配置文件里。
使用配置启动之后，会自动读取一次你的 github 仓库，把所有的文档加载到 go 里，后续如果你有 md 文件的改动，只需要访问一次

> http://127.0.0.1:5006/update
> 就可以完成数据的更新

- Docker

```
docker run -dit --rm --name=markdown-book \
-p 5006:5006 \
-v $(pwd)/md:/md -v $(pwd)/cache:/cache -v $(pwd)/config:/config \
hedongshu/markdown-book:latest --config ./config/config.yml
```

### 评论插件

> 评论插件使用的是 **Gitalk**，在使用前请阅读插件使用说明 [English](https://github.com/gitalk/gitalk/blob/master/readme.md) | [中文](https://github.com/gitalk/gitalk/blob/master/readme-cn.md)

#### 新增 `gitalk` 配置项，启动时加载配置文件即可

```yaml
gitalk:
  client-id: "你的 github oauth app client-id，必填。 如: ad549a9d085d7f5736d3"
  client-secret: "你的 github oauth app client-secret，必填。 如: 510d1a6bb875fd5031f0d613cd606b1d"
  repo: "你准备用于评论的项目名称，必填。 如: blog-issue"
  owner: "你的Github账号，必填。"
  admin:
    - "你的Github账号"
  labels:
    - "自定义issue标签，如: gitalk"
```

### 分析统计器

#### 百度

##### 1. 访问 https://tongji.baidu.com 创建站点，获取官方代码中的参数 `0952befd5b7da358ad12fae3437515b1`

```html
<script>
  var _hmt = _hmt || [];
  (function () {
    var hm = document.createElement("script");
    hm.src = "https://hm.baidu.com/hm.js?0952befd5b7da358ad12fae3437515b1";
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(hm, s);
  })();
</script>
```

##### 2. 配置

```shell
./markdown-book web --analyzer-baidu 0952befd5b7da358ad12fae3437515b1
```

#### 谷歌

##### 1. 访问 https://analytics.google.com 创建站点，获取官方代码中的参数 `G-MYSMYSMYS`

```html
<script
  async=""
  src="https://www.googletagmanager.com/gtag/js?id=G-MYSMYSMYS"
></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag() {
    dataLayer.push(arguments);
  }
  gtag("js", new Date());

  gtag("config", "G-MYSMYSMYS");
</script>
```

##### 2. 配置

```shell
./markdown-book web --analyzer-google G-MYSMYSMYS
```

### 标题栏图标

> 默认读取与程序运行同一级目录的 **favicon.ico** 文件

<!-- ### 导航排序

> 博客导航默认按照 `字典` 排序，可以通过 `@` 前面的数字来自定义顺序 -->

## 开发

1. 安装 `Golang` 开发环境

2. Fork [源码](https://github.com/gaowei-space/gocron)

3. 启动 web 服务

   运行之后访问地址 [http://localhost:5006](http://localhost:5006)，API 请求会转发给 `markdown-book` 程序

   ```
   make run
   ```

4. 编译

   在 **bin** 目录生成当前系统的压缩包，如：markdown-book-v1.1.0-darwin-amd64.tar

   ```
   make
   ```

5. 打包

   在 **package** 目录生成当前系统的压缩包，如：markdown-book-v1.1.0-darwin-amd64.tar

   ```
   make package
   ```

6. 生成 Windows、Linux、Mac 的压缩包

   在 **package** 生成压缩包，如：markdown-book-v1.1.0-darwin-amd64.tar markdown-book-v1.1.0-linux-amd64.tar.gz markdown-book-v1.1.0-windows-amd64.zip

   ```
   make package-all
   ```

### Docker

7. 下载

```
docker pull hedongshu/markdown-book:latest
```

2. 启动

   - 线上环境

   ```
   docker run -dit --rm --name=markdown-book \
   -p 5006:5006 \
   -v $(pwd)/md:/md -v $(pwd)/cache:/cache \
   hedongshu/markdown-book:latest
   ```

   - 开发环境

   ```
   docker run -dit --rm --name=markdown-book \
   -p 5006:5006 \
   -v $(pwd)/md:/md -v $(pwd)/cache:/cache \
   hedongshu/markdown-book:latest \
   -e dev
   ```

3. 访问 http://127.0.0.1:5006，查看效果

4. 其他用法

```
# 查看帮助
docker run -dit --rm --name=markdown-book \
    -p 5006:5006 \
    -v $(pwd)/md:/md -v $(pwd)/cache:/cache \
    hedongshu/markdown-book:latest -h


# 设置 title
docker run -dit --rm --name=markdown-book \
    -p 5006:5006 \
    -v $(pwd)/md:/md -v $(pwd)/cache:/cache \
    hedongshu/markdown-book:latest \
    -t "TechMan'Blog"


# 设置 谷歌统计
docker run -dit --rm --name=markdown-book \
    -p 5006:5006 \
    -v $(pwd)/md:/md -v $(pwd)/cache:/cache \
    hedongshu/markdown-book:latest \
    -t "TechMan'Blog" \
    --analyzer-google "De44AJSLDdda"
```

## 部署

> Nginx 反向代理配置文件参考

#### HTTP 协议

```
server {
    listen      80;
    listen [::]:80;
    server_name yourhost.com;

    location / {
         proxy_pass  http://127.0.0.1:5006;
         proxy_set_header   Host             $host;
         proxy_set_header   X-Real-IP        $remote_addr;
         proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
     }
}
```

#### HTTPS 协议（80 端口自动跳转至 443）

```
server {
    listen      80;
    listen [::]:80;
    server_name yourhost.com;

    location / {
        rewrite ^ https://$host$request_uri? permanent;
    }
}

server {
    listen          443 ssl;
    server_name     yourhost.com;
    access_log      /var/log/nginx/markdown-book.access.log main;


    #证书文件名称
    ssl_certificate /etc/nginx/certs/yourhost.com_bundle.crt;
    #私钥文件名称
    ssl_certificate_key /etc/ngpackinx/certs/yourhost.com.key;
    ssl_session_timeout 5m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
    ssl_prefer_server_ciphers on;

    location / {
         proxy_pass  http://127.0.0.1:5006;
         proxy_set_header   Host             $host;
         proxy_set_header   X-Real-IP        $remote_addr;
         proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
     }
 }
```
