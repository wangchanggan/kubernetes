# Kubernetes源码分析

Source Code From
https://github.com/kubernetes/kubernetes/releases/tag/v1.21.0

参考Kubernetes源码分析（基于Kubernetes 1.14版本）（郑东旭/著）

## 目录

-   [Kubernetes源码分析](#kubernetes源码分析)
    -   [目录](#目录)
    -   [源码目录结构说明](#源码目录结构说明)
    -   [架构](#架构)
    -   [核心数据结构](#核心数据结构)
        -   [APIResourceList](#apiresourcelist)
        -   [Group](#group)
        -   [Version](#version)
        -   [Resource](#resource)
            -   [资源外部版本与内部版本](#资源外部版本与内部版本)
            -   [资源代码定义](#资源代码定义)
            -   [资源注册到资源注册表](#资源注册到资源注册表)
            -   [资源首选版本](#资源首选版本)
            -   [资源操作方法](#资源操作方法)
            -   [资源与命名空间](#资源与命名空间)
        -   [runtime.Object类型基石](#runtime.object类型基石)
        -   [Unstructured数据](#unstructured数据)
        -   [Scheme资源注册表](#scheme资源注册表)
            -   [Scheme资源注册表数据结构](#scheme资源注册表数据结构)
            -   [资源注册表注册方法](#资源注册表注册方法)
        -   [Codec编解码器](#codec编解码器)
            -   [Codec编解码实例化](#codec编解码实例化)
            -   [jsonSerializer与yamlSerializer序列化器](#jsonserializer与yamlserializer序列化器)
            -   [protobufSerializer序列号器](#protobufserializer序列号器)
        -   [Converter资源版本转换器](#converter资源版本转换器)
            -   [Converter转换器数据结构](#converter转换器数据结构)
            -   [Converter注册转换函数](#converter注册转换函数)
            -   [Converter资源版本转换原理](#converter资源版本转换原理)
    -   [kubectl命令行交互](#kubectl命令行交互)
        -   [创建资源对象的过程](#创建资源对象的过程)
            -   [实例化Factory接口](#实例化factory接口)
            -   [Builder构建资源对象](#builder构建资源对象)
            -   [Visitor多层匿名函数嵌套](#visitor多层匿名函数嵌套)
    -   [client-go编程式交互](#client-go编程式交互)
    -   [Etcd存储核心实现](#etcd存储核心实现)
        -   [RegistryStore存储服务通用操作](#registrystore存储服务通用操作)
        -   [Storage.Interface通用存储接口](#storage.interface通用存储接口)
        -   [CacherStorage缓存层](#cacherstorage缓存层)
            -   [CacherStorage缓存层设计](#cacherstorage缓存层设计)
            -   [watchCache 缓存滑动窗口](#watchcache-缓存滑动窗口)
        -   [UnderlyingStorage底层存储对象](#underlyingstorage底层存储对象)
        -   [Codec编解码数据](#codec编解码数据)
        -   [Strategy预处理](#strategy预处理)
            -   [创建资源对象时的预处理操作](#创建资源对象时的预处理操作)
            -   [更新资源对象时的预处理操作](#更新资源对象时的预处理操作)
            -   [删除资源对象时的预处理操作](#删除资源对象时的预处理操作)
    -   [kube-apiserver核心实现](#kube-apiserver核心实现)
        -   [热身概念](#热身概念)
            -   [go-restful核心原理](#go-restful核心原理)
            -   [OpenAPI/Swagger核心原理](#openapiswagger核心原理)
            -   [gRPC核心原理](#grpc核心原理)
        -   [kube-apiserver启动流程](#kube-apiserver启动流程)
            -   [资源注册](#资源注册)
            -   [Cobra命令行参数解析](#cobra命令行参数解析)
            -   [创建APIServer通用配置](#创建apiserver通用配置)
            -   [创建APIExtensionsServer](#创建apiextensionsserver)
            -   [创建KubeAPIServer](#创建kubeapiserver)
            -   [创建AggregatorServer](#创建aggregatorserver)
            -   [创建GenericAPIServer（以创建APIExtensionsSever为例）](#创建genericapiserver以创建apiextensionssever为例)
            -   [启动HTTP服务](#启动http服务)
            -   [启动HTTPS服务](#启动https服务)
        -   [认证](#认证)
            -   [BasicAuth认证](#basicauth认证)
            -   [ClientCA认证](#clientca认证)
            -   [TokenAuth认证](#tokenauth认证)
            -   [BootstrapToken认证](#bootstraptoken认证)
            -   [RequestHeader认证](#requestheader认证)
            -   [Webhook TokenAuth认证](#webhook-tokenauth认证)
            -   [Anonymous认证](#anonymous认证)
            -   [OIDC认证](#oidc认证)
            -   [ServiceAccountAuth认证](#serviceaccountauth认证)
        -   [授权](#授权)
            -   [AlwaysAllow授权](#alwaysallow授权)
            -   [AlwaysDeny授权](#alwaysdeny授权)
            -   [ABAC授权](#abac授权)
            -   [Webhook授权](#webhook授权)
            -   [RBAC授权](#rbac授权)
            -   [Node授权](#node授权)
        -   [准入控制器](#准入控制器)
            -   [AlwaysPullImages准入控制器](#alwayspullimages准入控制器)
            -   [PodNodeSelector准入控制器](#podnodeselector准入控制器)
        -   [进程信号处理机制](#进程信号处理机制)
            -   [常驻进程实现](#常驻进程实现)
            -   [进程的优雅关闭](#进程的优雅关闭)
            -   [向systemd报告进程状态](#向systemd报告进程状态)

## 源码目录结构说明
| 源码目录 | 说明 | 备注 |
| :----: | :---- | :---- |
| cmd/ | 存放可执行文件的入口代码,每个可执行文件都会对应一个 main函数 | |
| pkg/ | 存放核心库代码,可被项目内部或外部直接引用 | |
| vendor/ | 存放项目依赖的库代码,一般为第三方库代码 | |
| api/ | 存放 OpenAPI/Swagger 的spec文件,包括JSON、Protocol的定义等 | |
| build/ | 存放与构建相关的脚本 | |
| test/ | 存放测试工具及测试数据 | |
| docs/ | 存放设计或用户使用文档 | |
| hack/ | 存放与构建、测试等相关的脚本 | |
| third party/ | 存放第三方工具、代码或其他组件 | |
| plugin/ | 存放Kubernetes插件代码目录，例如认证、授权等相关插件 | |
| staging/ | 存放部分核心库的暂存目录 | 已将该目录下src中文件夹和文件迁移至vender/，否则代码无法编译运行 |
| translations/ | 存放il8n(国际化）语言包的相关文件,可以在不修改内部代码的情况下支持不同语言及地区 | |
| examples/ | 存放源码分析中的示例代码 | |



## 架构
见docs/架构.docx



## 核心数据结构
见docs/核心数据结构.doc


### APIResourceList
vendor/k8s.io/apimachinery/pkg/apis/meta/v1/type.go:1077

vendor/k8s.io/apimachinery/pkg/runtime/schema（资源数据结构）


### Group
vendor/k8s.io/apimachinery/pkg/apis/meta/v1/type.go:981


### Version
vendor/k8s.io/apimachinery/pkg/apis/meta/v1/type.go:953


### Resource
vendor/k8s.io/apimachinery/pkg/apis/meta/v1/type.go:1032

#### 资源外部版本与内部版本
资源外部版本：vendor/k8s.io/api/\<group>/\<version>/\<resource file>

资源内部版本：pkg/apis/\<group>

#### 资源代码定义
##### 资源内部版本
pkg\apis\apps

|--doc.go: GoDoc文件，定义了当前包的注释信息。在Kubernetes资源包中，它还担当了代码生成器的全局Tags描述文件。

|--register.go: 定义了资源组、资源版本及资源的注册信息。通过runtime.APIVersionIntermal (即_ intemal) 标识。

|--types.go:定义了在当前资源组、资源版本下所支持的资源类型。

|--v1、 v1beta1、 v1beta2: 定义了资源组下拥有的资源版本的资源(即外部版本)。

|--install: 把当前资源组下的所有资源注册到资源注册表中。

|--validation:定义了资源的验证方法。

|--zz_generated.deepcopy.go: 定义了资源的深复制操作，该文件由代码生成器自动生成。
##### 资源外部版本
pkg/apis/apps/{v1, v1beta1, v1beta2}

|--conversion.go: 定义了资源的转换函数(默认转换函数)，并将默认转换函数注册到资源注册表中。

|--zz_generated.conversion.go: 定义了资源的转换函数(自动生成的转换函数)，并将生成的转换函数注册到资源注册表中。该文件由代码生成器自动生成。

|--defaults.go: 定义了资源的默认值函数，并将默认值函数注册到资源注册

|--zz_generated.defaults.go:定义了资源的默认值函数(自动生成的默认值函数)，并将生成的默认值函数注册到资源注册表中。该文件由代码生成器自动生成。

|--register.go: 定义了资源组、资源版本及资源的注册信息。通过资源版本（Alpha、Beta、Stable）标识。

#### 资源注册到资源注册表
在每一个 Kubernetes资源组目录中，都拥有一个install/install.go代码文件，它负责将资源信息注册到资源注册表(Scheme)中。
以core核心资源组为例，代码示例:pkg/apis/core/install/install.go

#### 资源首选版本
首选版本（Preferred Version)，也称优选版本（Priority Version)，一个资源组下拥有多个资源版本，在一些场景下，如不指定资源版本，则使用该资源的首选版本。

以 apps资源组为例，注册资源时会注册多个资源版本：pkg/apis/apps/install/install.go

当通过资源注册表获取所有资源组下的首选版本时，将位于最前面的资源版本作为首选版本：vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:664

获取指定资源组的资源版本，按照优先顺序返回：vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:614

获取所有资源组的资源版本，按优先顺序返回：vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:640

#### 资源操作方法
通过metav1.Verbs数据结构进行描述：vendor/k8s.io/apimachinery/pkg/apis/meta/v1/types.go:1079

| 操作方法 (Verbs) | 操作方法接口(Interface) | 说明 |
| :----:| :----: | :---- |
| create | rest.Creater | 资源对象创建接口 |
| delete | rest.GracefulDeleter | 资源对象删除接口(单个资源对象) |
| deletecollection | rest.CllctionDeleter | 资源对象删除接口(多个资源对象) |
| update | rest.Updater | 资源对象更新接口 (完整资源对象的更新) |
| patch | rest.Patcher | 资源对象更新接口(局部资源对象的更新) |
| get | rest.Getter | 资源对象获取接口(单个资源对象) |
| list | rest.Lister | 资源对象获取接口(多个资源对象) |
| watch | rest.Watcher | 资源对象监控接口 |

接口定义：vendor/k8s.io/apiserver/pkg/registry/rest/rest.go

接口实现（以Pod资源对象为例）:

pkg/registry/core/pod/storage/storage.go

vendor/k8s.io/apiserver/pkg/registry/generic/registry/store.go

pkg/registry/core/pod/rest/log.go

#### 资源与命名空间
以Pod为例：vendor/k8s.io/apimachinery/pkg/apis/meta/v1/types.go:153

### runtime.Object类型基石
vendor/k8s.io/apimachinery/pkg/runtime/interfaces.go:299

examples/runtime_object.go

### Unstructured数据
vendor/k8s.io/apimachinery/pkg/runtime/interfaces.go:333

### Scheme资源注册表
#### Scheme资源注册表数据结构
vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:46

examples/scheme.go

注：UnversionedType类型的对象在通过scheme.AddUnversionedTypes方法注册时，会同时存在于4个map结构中vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:148

#### 资源注册表注册方法
以scheme.AddKnownTypes方法为例vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:166

### Codec编解码器
Codec编解码器通用接口定义：vendor/k8s.io/apimachinery/pkg/runtime/interfaces.go:94

#### Codec编解码实例化
vendor/k8s.io/apimachinery/pkg/runtime/serializer/codec_factory.go:52

#### jsonSerializer与yamlSerializer序列化器
序列化操作vendor/k8s.io/apimachinery/pkg/runtime/serializer/json/json.go:340

反序列化操作vendor/k8s.io/apimachinery/pkg/runtime/serializer/json/json.go:209

#### protobufSerializer序列号器
序列化操作vendor/k8s.io/apimachinery/pkg/runtime/serializer/protobuf/protobuf.go:173

反序列化操作vendor/k8s.io/apimachinery/pkg/runtime/serializer/protobuf/protobuf.go:99

vendor/k8s.io/apimachinery/pkg/runtime/serializer/protobuf/protobuf.go:381

### Converter资源版本转换器
#### Converter转换器数据结构
vendor/k8s.io/apimachinery/pkg/conversion/converter.go:41

#### Converter注册转换函数
支持5个注册转换函数vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:345

#### Converter资源版本转换原理
examples/converter.go

1.获取传入的资源对象的反射类型vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:495

2.从资源注册表中查找到传入的资源对象的GVK vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:510

3.从多个GVK中选出与目标资源对象相匹配的GVK vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:516

4.判断传入的资源对象是否属于Unversioned类型vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:538

5.执行转换操作

vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:558

vendor/k8s.io/apimachinery/pkg/conversion/converter.go:210

6.设置转换后资源对象的GVK vendor/k8s.io/apimachinery/pkg/runtime/scheme.go:601



## kubectl命令行交互
见docs/kubectl命令行交互.doc
### Cobra命令行参数解析
1.创建Command vendor/k8s.io/kubectl/pkg/cmd/cmd.go:472

2.以基础目命令（中级）get命令为例，添加命令行参数vendor/k8s.io/kubectl/pkg/cmd/get/get.go:155

3.执行命令cmd/kubectl/kubectl.go

### 创建资源对象的过程
#### 实例化Factory接口
cmd/kubectl-convert/kubectl-convert.go:41

vendor/k8s.io/kubectl/pkg/cmd/util/factory.go:44

#### Builder构建资源对象
pkg/kubectl/cmd/convert/convert.go:122

pkg/kubectl/cmd/convert/convert.go:142

#### Visitor多层匿名函数嵌套
vendor/k8s.io/cli-runtime/pkg/resource/interfaces.go:59

examples/visitor.go

Visitor中的VisitorList (存放Visitor 的集合)有两种

vendor/k8s.io/cli-runtime/pkg/resource/visitor.go:189

vendor/k8s.io/cli-runtime/pkg/resource/visitor.go:203



## client-go编程式交互
见https://github.com/wangchanggan/client-go



## Etcd存储核心实现
见docs/Etcd存储核心实现.doc
### RESTStorage存储服务通用接口
vendor/k8s.io/apiserver/pkg/registry/rest/rest.go:57

Kubernetes的每种资源实现的RESTStorage 接口一般定义在pkg/registy/<资源组>/<资源>/storage/storage.go中，它们通过NewStorage函数或NewREST函数实例化。
以Deployment资源为例pkg/registry/apps/deployment/storage/storage.go:51,128


### RegistryStore存储服务通用操作
RegistryStore结构vendor/k8s.io/apiserver/pkg/registry/generic/registry/store.go:94

以RegistryStore的Create方法（创建资源对象的方法）为例vendor/k8s.io/apiserver/pkg/registry/generic/registry/store.go:369


### Storage.Interface通用存储接口
Storage.Interface通用存储接口定义的资源操作方法vendor/k8s.io/apiserver/pkg/storage/interfaces.go:195

实现通用存储接口的资源存储对象

vendor/k8s.io/apiserver/pkg/storage/cacher/cacher.go:226

vendor/k8s.io/apiserver/pkg/storage/etcd3/store.go:67

实例化过程vendor/k8s.io/apiserver/pkg/server/options/etcd.go:249


### CacherStorage缓存层
#### CacherStorage缓存层设计
1.cacheWatcher

vendor/k8s.io/apiserver/pkg/storage/cacher/cacher.go:454,1394

2.watchCache

vendor/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go:281

3.Cacher

vendor/k8s.io/apiserver/pkg/storage/cacher/cacher.go:1261


#### watchCache 缓存滑动窗口
vendor/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go:135,350,578


### UnderlyingStorage底层存储对象
vendor/k8s.io/apiserver/pkg/storage/storagebackend/factory/etcd3.go:225

以Get操作为例vendor/k8s.io/apiserver/pkg/storage/etcd3/store.go:115,947


### Codec编解码数据
examples/codec.go


### Strategy预处理
vendor/k8s.io/apiserver/pkg/registry/generic/registry/store.go:72

#### 创建资源对象时的预处理操作
vendor/k8s.io/apiserver/pkg/registry/rest/create.go:40,80

#### 更新资源对象时的预处理操作
vendor/k8s.io/apiserver/pkg/registry/rest/update.go:40,96

#### 删除资源对象时的预处理操作
vendor/k8s.io/apiserver/pkg/registry/rest/delete.go:35,59,76



## kube-apiserver核心实现
见docs/kube-apiserver核心实现.doc

### 热身概念

#### go-restful核心原理
vendor/github.com/emicklei/go-restful/container.go:205

#### OpenAPI/Swagger核心原理
Kubernetes在注册go restful路由时，将资源信息与OpenAPI自定义扩展属性进行了关联vendor/k8s.io/apiserver/pkg/endpoints/installer.go:992

#### gRPC核心原理
1.引用类Tags

vendor/k8s.io/apimachinery/pkg/apis/meta/v1/time.go:32

vendor/k8s.io/apimachinery/pkg/apis/meta/v1/time_proto.go:26

vendor/k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto:1038

2.嵌入类Tags

vendor/k8s.io/apimachinery/pkg/api/resource/quantity.go:90

vendor/k8s.io/apimachinery/pkg/api/resource/generated.proto:86

3.go-to-protobuf的生成规则

vendor/k8s.io/code-generator/cmd/go-to-protobuf/protobuf/generator.go:99

vendor/k8s.io/code-generator/cmd/go-to-protobuf/protobuf/cmd.go:106


### kube-apiserver启动流程
#### 资源注册
以KubeAPIServer (API核心服务)为例cmd/kube-apiserver/app/server.go:73

1.初始化Scheme资源注册表pkg/api/legacyscheme/scheme.go

2.注册Kubernetes所支持的资源pkg/controlplane/import_known_versions.go

#### Cobra命令行参数解析
cmd/kube-apiserver/app/server.go:107

#### 创建APIServer通用配置
1.genericConfig实例化

cmd/kube-apiserver/app/server.go:461

pkg/controlplane/instance.go:665

2.OpenAPI/Swagger配置 cmd/kube-apiserver/app/server.go:487

3.StorageFactory存储(Etcd)配置 cmd/kube-apiserver/app/server.go:500

4.Authentication认证配置 pkg/kubeapiserver/authenticator/config.go:207

5.Authorization授权配置 pkg/kubeapiserver/authorizer/config.go:77

6.Admission准入控制器配置

vendor/k8s.io/apiserver/pkg/admission/plugins.go:39

vendor/k8s.io/apiserver/pkg/server/plugins.go

pkg/kubeapiserver/options/plugins.go:110

以AlwaysPullImages准入控制器为例，注册方法plugin/pkg/admission/alwayspullimages/admission.go:45

#### 创建APIExtensionsServer
1.创建 GenericAPIServer vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go:130

2.实例化 CustomResourceDefinitions vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go:136

3.实例化 APIGroupInfo

vendor/k8s.io/apiserver/pkg/server/genericapiserver.go:64

vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go:149

4.InstallAPIGroup注册APIGroup

vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go:185

vendor/k8s.io/apiserver/pkg/endpoints/groupversion.go:109

#### 创建KubeAPIServer
1.创建GenericAPIServer pkg/controlplane/instance.go:355

2.实例化Master pkg/controlplane/instance.go:396

3.InstallLegacyAPI注册/api资源 pkg/controlplane/instance.go:403,546

4.InstallAPIs注册/apis资源 pkg/controlplane/instance.go:592

#### 创建AggregatorServer
1.创建GenericAPIServer vendor/k8s.io/kube-aggregator/pkg/apiserver/apiserver.go:171

2.实例化APIAggregator vendor/k8s.io/kube-aggregator/pkg/apiserver/apiserver.go:186

3.实例化APIGroupInfo

vendor/k8s.io/kube-aggregator/pkg/apiserver/apiserver.go:209

vendor/k8s.io/kube-aggregator/pkg/registry/apiservice/rest/storage_apiservice.go:34

4.installAPIGroup注册APIGroup vendor/k8s.io/kube-aggregator/pkg/apiserver/apiserver.go:215

#### 创建GenericAPIServer（以创建APIExtensionsSever为例）
vendor/k8s.io/apiserver/pkg/server/handler.go:73

#### 启动HTTP服务
vendor/k8s.io/apiserver/pkg/server/deprecated_insecure_serving.go:43

vendor/k8s.io/apiserver/pkg/server/secure_serving.go:207

#### 启动HTTPS服务
vendor/k8s.io/apiserver/pkg/server/secure_serving.go:147


### 认证
vendor/k8s.io/apiserver/pkg/endpoints/filters/authentication.go:46

vendor/k8s.io/apiserver/pkg/authentication/request/union/union.go:54

#### BasicAuth认证
通过go语言标准库实现 $GOROOT/src/net/http/request.go:922

#### ClientCA认证
vendor/k8s.io/apiserver/pkg/authentication/request/x509/x509.go:201

#### TokenAuth认证
vendor/k8s.io/apiserver/plugin/pkg/authenticator/token/tokentest/tokentest.go:39

#### BootstrapToken认证
plugin/pkg/auth/authenticator/token/bootstrap/bootstrap.go:95

#### RequestHeader认证
vendor/k8s.io/apiserver/pkg/authentication/request/headerrequest/requestheader.go:160

#### Webhook TokenAuth认证
vendor/k8s.io/apiserver/pkg/authentication/token/cache/cached_token_authenticator.go:127,138

vendor/k8s.io/apiserver/plugin/pkg/authenticator/token/webhook/webhook.go:86

#### Anonymous认证
vendor/k8s.io/apiserver/pkg/authentication/request/anonymous/anonymous.go:35

#### OIDC认证
vendor/k8s.io/apiserver/plugin/pkg/authenticator/token/oidc/oidc.go:537

#### ServiceAccountAuth认证
pkg/serviceaccount/jwt.go:261

### 授权
vendor/k8s.io/apiserver/pkg/authorization/authorizer/interfaces.go

#### AlwaysAllow授权
vendor/k8s.io/apiserver/pkg/authorization/authorizerfactory/builtin.go:32,,37

#### AlwaysDeny授权
vendor/k8s.io/apiserver/pkg/authorization/authorizerfactory/builtin.go:62,68

#### ABAC授权
pkg/auth/authorizer/abac/abac.go:228,242

#### Webhook授权
vendor/k8s.io/apiserver/plugin/pkg/authorizer/webhook/webhook.go:161,233

#### RBAC授权
1.RBAC授权详解

examples\rbac.go

plugin/pkg/auth/authorizer/rbac/rbac.go:75

2.内置集群角色
plugin/pkg/auth/authorizer/rbac/bootstrappolicy/policy.go:34，580

#### Node授权
plugin/pkg/auth/authorizer/node/node_authorizer.go:94

### 准入控制器
vendor/k8s.io/apiserver/pkg/admission/interfaces.go:129

vendor/k8s.io/apiserver/pkg/admission/chain.go:247,36,53

#### AlwaysPullImages准入控制器
plugin/pkg/admission/alwayspullimages/admission.go:64,85

#### PodNodeSelector准入控制器
plugin/pkg/admission/podnodeselector/admission.go:63,104,143


### 进程信号处理机制
#### 常驻进程实现
vendor/k8s.io/apiserver/pkg/server/signal.go:37

#### 进程的优雅关闭
vendor/k8s.io/apiserver/pkg/server/genericapiserver.go:334

#### 向systemd报告进程状态
vendor/k8s.io/apiserver/pkg/server/genericapiserver.go:422

## kube-scheduler核心实现
见docs/kube-scheduler核心实现.doc
### kube-scheduler组件的启动流程
#### Cobra命令行参数解析
cmd/kube-scheduler/app/server.go:64

#### 内置调度算法的注册
pkg/scheduler/factory.go:193

### 实例化Scheduler对象
1.实例化所有的Informer cmd/kube-scheduler/app/server.go:318
2.实例化调度算法函数 pkg/scheduler/apis/config/types.go:135
3.为所有Informer对象添加对资源事件的监控 pkg/scheduler/eventhandlers.go:359

#### 运行EventBroadcaster事件管理器
cmd/kube-scheduler/app/server.go:159

#### 运行HTTP或HTTPS服务
kube-scheduler组件也拥有自己的HTTP服务，但功能仅限于监控及监控检查等，其运行原理与kube-apiserver组件的类似。

/healthz：用于健康检查。

/metries：用于监控指标，一般用于Prometheus指标采集。

/debug/pprof：用于pprof性能分析。

#### 运行Informer同步资源
cmd/kube-scheduler/app/server.go:208

#### 领导者选举实例化
cmd/kube-scheduler/app/server.go:211

#### 运行sched.Run调度器
pkg/scheduler/scheduler.go:314


### 亲和性调度
vendor/k8s.io/api/core/v1/types.go:2688

#### NodeAffinity
vendor/k8s.io/api/core/v1/types.go:2825

#### PodAffinity
vendor/k8s.io/api/core/v1/types.go:2707

#### PodAntiAffinity
vendor/k8s.io/api/core/v1/types.go:2747


### 内置调度算法
#### 预选调度算法
pkg/scheduler/testing/fake_extender.go:36

#### 优选调度算法
pkg/scheduler/framework/extender.go:41


### 调度器核心实现
#### 调度过程
1.预选调度前的性能优化 pkg/scheduler/core/generic_scheduler.go:190

2.预选调度过程

pkg/scheduler/core/generic_scheduler.go:285

pkg/scheduler/framework/runtime/framework.go:651

3.优选调度过程

pkg/scheduler/core/generic_scheduler.go:425

以 leastRequestedPriority 优选调度算法为例 pkg/scheduler/framework/plugins/noderesources/least_allocated.go:97

4.选择一个最佳节点 pkg/scheduler/core/generic_scheduler.go:156

#### Preempt抢占机制
1.判断当前Pod资源对象是否有资格抢占其他Pod资源对象所在的节点 

pkg/scheduler/framework/plugins/defaultpreemption/default_preemption.go:133

2.从预选调度失败的节点中尝试找到能够调度成功的节点列表(潜在的节点列表)

pkg/scheduler/framework/plugins/defaultpreemption/default_preemption.go:141,276

3.从潜在的节点列表中尝试找到能够抢占成功的节点列表(驱逐的节点列表)

pkg/scheduler/framework/plugins/defaultpreemption/default_preemption.go:333

4.从驱逐的节点列表中选择一个节点用于最终被抢占的节点(被抢占节点)

pkg/scheduler/framework/plugins/defaultpreemption/default_preemption.go:475

5.获取被抢占节点上的所有NominatedPods列表

pkg/scheduler/framework/plugins/defaultpreemption/default_preemption.go:744

#### bind绑定机制
pkg/scheduler/scheduler.go:400


### 领导者选举机制
#### 资源锁
vendor/k8s.io/client-go/tools/leaderelection/resourcelock/interface.go:47,86

#### 领导者选举过程
vendor/k8s.io/client-go/tools/leaderelection/leaderelection.go:196

1.资源锁获取过程

vendor/k8s.io/client-go/tools/leaderelection/leaderelection.go:241,327

2.领导者节点定时更新租约过程

vendor/k8s.io/client-go/tools/leaderelection/leaderelection.go:270





