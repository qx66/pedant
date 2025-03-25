# pedant (书呆子)

pedant是一个对接大模型语言的，提供Rest Api接口的模块

主要初衷是对接测试市面上的大模型语言，给内部一个简单的Http Rest Api访问接口

如果其他模块需要直接快速的原生的使用大模型语言，可以直接倒入使用。(聚合仓库的作用)

## 初始化数据库

本模块使用MySQL数据库进行存储数据

[init.sql](docs%2Finit.sql)

下载 init.sql，将创建表语句复制执行即可

## Api参考

[pedant.postman_collection.json](docs%2Fpedant.postman_collection.json)

## ChatGpt

需要设置全局代理

## BaiduCloud

[文心一言](https://yiyan.baidu.com/)

[文心插件](https://yiyan.baidu.com/pluginSubmission)

[飞桨应用](https://aistudio.baidu.com/index/creations/application)

## 火山引擎

[火山引擎](https://www.volcengine.com/docs/82379/1319853)

## Google Gemini AI

[GoogleGemini](https://makersuite.google.com/app/prompts/new_freeform)

```markdown
message: User location is not supported for the API use

解决办法: 设置全局代理
```


