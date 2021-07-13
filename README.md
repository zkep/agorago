# agorago
根据声网restfulapi接口的go版本实现

```
    ├── LICENSE
    ├── README.md
    ├── access_token.go     
    ├── cloud_recording.go
    ├── cloud_recording_test.go
    ├── cloud_recording_types.go
    ├── generate       // 根据声网swagger的yaml文档生成对应的golang代码
    │   ├── README.md
    │   ├── g_api.go         
    │   ├── g_api_docs.go    
    │   ├── g_api_test.go   // 测试
    │   ├── server_restfulapi_cn.yaml    // 测试yaml文件
    │   └── swagger.go     // swagger文件解析
    ├── go.mod
    ├── go.sum
    ├── request.go   
    └── token_builder.go
```