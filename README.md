# Kubernetes源码分析

Source Code From https://github.com/kubernetes/kubernetes/releases/tag/v1.21.0

参考Kubernetes源码分析（基于Kubernetes 1.14版本）（郑东旭/著）

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
