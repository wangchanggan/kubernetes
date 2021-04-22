# Kubernetes源码分析

Source Code From https://github.com/kubernetes/kubernetes/releases/tag/v1.21.0

参考Kubernetes源码分析（基于Kubernetes 1.14版本）（郑东旭/著）

## 源码目录结构说明
| 源码目录 | 说明 |
| :----: | :---- |
| cmd/ | 存放可执行文件的入口代码,每个可执行文件都会对应一个 main函数 |
| pkg/ | 存放核心库代码,可被项目内部或外部直接引用 |
| vendor/ | 存放项目依赖的库代码,一般为第三方库代码 |
| api/ | 存放 OpenAPI/Swagger 的spec文件,包括JSON、Protocol的定义等 |
| build/ | 存放与构建相关的脚本 |
| test/ | 存放测试工具及测试数据 |
| docs/ | 存放设计或用户使用文档 |
| hack/ | 存放与构建、测试等相关的脚本 |
| third party/ | 存放第三方工具、代码或其他组件 |
| plugin/ | 存放Kubernetes插件代码目录，例如认证、授权等相关插件 |
| staging/ | 存放部分核心库的暂存目录 |
| translations | 存放il8n(国际化）语言包的相关文件,可以在不修改内部代码的情况下支持不同语言及地区 |
