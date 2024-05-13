# GoFly快速开发后台管理系统介绍


## 二、优势简介

1. 基于优秀成熟框架集成，保证系统可靠性。集成的主要有 Gin、Arco Design 、Mysql 等主流框架技术《我们不生产框架，我们是优秀框架的搬运工》。
2. 系统已集成开发常用基础功能，开箱即用，快速开始您业务开发，快人一步，比同行节省成本，降本增效首选。
3. 框架提供其他开发者开发的插件，可快速安装或卸载，让开个资源共享，同意功能无需重复造车，一键安装即可使用。 框架搭建了一键 CRUD 生成前后端代码，建数据库一键生成，节省您的复制粘贴时间，进一步为您节省时间。
4. 框架自带 API 接口文档管理，接口带有请求 token 等配置，添加接口只需配置路径和数据库或者备注，其部分信息如数据字段，系统自动根据数据库字段补齐，开发配套接口文档尽可能的为您节省一点时间。不需要其他接口文档工具复制粘贴，登录注册等时间。还有一个重点！接口文档可以一键生成接口 CRUD 的代码和通用的操作数据的 CRUD 接口，根据您的业务选择自己写接口代码、一键生成接口代码、不用写和生成代码调用通用接口。让写接口工作节省更多时间。
5. 前后端分离解耦业务，让前段人员与后端人协调开发，提高项目交付，并且可以开发出功能复杂度高的项目。
6. 前端用 Vue3+TypeScript 的 UI 框架 [Arco Design](https://arco.design/vue/component/button)，好用的 UI 框架前端可以设计出优秀且交互不错的界面，完善的大厂 UI 支持，前端开发效率也很高！ 以上只是框架比较明显优势点，还有很多优势等你自己体验，我们从各个开发环节，努力为您节省每一分时间。
7. 集成操作简单的 ORM 框架，操作数据非常简单，就像使用php的Laravel一样，您可以去文档看看 [框架的ROM数据库操作文档](https://doc.goflys.cn/docview?id=25&fid=289)
   例如下面语句就可以查找一条数据：

```
  db.Table("users").Fields("uid,name,age").First()
```

8. 框架以“大道至简，唯快不破”为思想，在每个细节处理都坚持让“开发”变得简单，即使你是新手也可以跟着开发文档快手上手并能开发出企业级产品。
9. 框架有软著可靠可控、不留后门、不加密、不设限制、不收费，开发者可以放心使用。

## 三、目录结构

```
├── app                     # 应用目录
│   ├── admin               # 后台管理应用模块
│   ├── business            # 业务端应用模块
│   ├── common              # 公共应用模块
│   ├── home                # 可以编写平台对应网站
│   ├── wxapp               # 微信小程序模块
│   ├── wxoffi              # 微信公众号模块
│   └── controller.go       # 应用控制器
├── bootstrap               # 工具方法
├── global                  # 全局变量
├── model                   # 数据模型
├── resource                # 静态资源和config配置文件
├── route                   # 路由
├── runtime                 # 运行日志文件
├── tmp                     # 开发是使用fresh热编译 产生临时文件
├── utils                   # 工具包
├── go.mod                  # 依赖包管理工具
├── go.sum     
├── main.go                 # main函数    
└── README.md               # 项目介绍
```

开发时仅需在app目录下添加你新的需求，app外部文件建议不要改动，除了config配置需要改，其他不要修改，
框架已经为您封装好，你只需在app应用目录书写你的业务，路由、访问权限、跨域、限流、Token验证、ORM等
框架已集成好，开发只需添加新的方法或者新增一个文件即可。

## 四、快速安装

1. 首先在GOPATH路径下的src目录下现在放代码的文件夹下载代码解压到项目文件夹中（或者直接git clone 代码到src目录下）。
2. 再运行服务 go run main.go 或者 编译 fresh (go install github.com/pilu/fresh@latest 安装fresh热编译工具)，启动成功如下：
   ![运行启动命令](https://api.goflys.cn/common/uploadfile/get_image?url=resource/uploads/20230912/00ab0aa6dbbaea7135421d9d58fc7d53.png)
   在浏览器打开安装界面进行安装：
   ![安装界面](https://api.goflys.cn/common/uploadfile/get_image?url=resource/uploads/20240219/30b533c7a6d3bf711498089dd0f1337f.png)

注意：前端代码安装设置是安装时同时把前端vue代码安装到开发前端代码目录下，为了防止热编译效率框架不建议把前端代码放到go目录下。

## 五、在线预览

 [1.GoFly全栈开发社区了解更多](https://goflys.cn/home)

 [2.Go快速后台系统开发框架完整代码包下载](https://goflys.cn/prdetail?id=6)

 [3.Go快速后台系统开发文档](https://doc.goflys.cn/docview?id=25)

 [4.A端在线预览](https://sg.goflys.cn/webadmin)

 [5.B端在线预览](https://sg.goflys.cn/webbusiness)



## 七、安装及部署打包说明

### 1. 后端代码

#### 安装fresh 热更新-边开发边编译

go install github.com/pilu/fresh@latest

#### 初始化mod

go mod tidy

#### 热编译运行

bee run 或 fresh

#### 打包

go build main.go

#### 打包（此时会打包成Linux上可运行的二进制文件，不带后缀名的文件）

```
SET GOOS=linux
SET GOARCH=amd64
go build
```

#### widows

```
// 配置环境变量
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

go build main.go

// 编译命令
```

#### 编译成Linux环境可执行文件

```

// 配置参数
SET CGO_ENABLED=0 
SET GOOS=linux 
SET GOARCH=amd64 

go build main.go

// 编译命令
```

#### 服务器部署

部署是把打包生成的二进制文件(Linux:gofly，windows:gofly.exe)和资源文件resource复制过去即可。

### 2. 前端端代码

#### 初始化依赖

```
 npm run install 或者 yarn install
```

如果第一次使用Arco Design Pro install初始化可以报错，如果保存请运行下面命令（安装项目模版的工具）：

```
npm i -g arco-cli
```

#### 运行

```
npm run serve 或者  yarn serve
```

#### 打包

```
npm run build 或者 yarn build
```

</div>
