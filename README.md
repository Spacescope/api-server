# Data-infra-api-backend
Data-infra-api backend server, provide data api.

### Regenerate swagger doc
swagger doc defined in router api comment.
if edited these comments, need to regenerate swagger doc.
```shell script
swag init -g cmd/api-server/main.go
```

### Swagger doc
swagger doc please refer to
`http://127.0.0.1:7003/api-server/swagger/index.html`

### How to make
```
make # make to see help
```
### Run
1. Config -- data-infra-backend/config
    ```
    DB of Postgresql
    DB of Data Observable
    KV of redis
    ```
2. Run
```
docker run -v /home/ec2-user/api-server/config:/etc/api-server/conf -p 7003:7003 -d 129862287110.dkr.ecr.us-east-2.amazonaws.com/data-infra/data-api-server:commitId
```

### Refer
1. [Data api system design doc](https://docs.google.com/document/d/1QoA4ZNfGSCqvZPHf5l2D11NZchuTTy9oh3-qkVUC3BM/edit#heading=h.qsczq5hva7v4)
2. [Data api delopyment doc](https://docs.google.com/document/d/1oOXyUavXw4uGl-DBJEi5HuISR4tmVkUAgmdt4ixNUi8/edit#heading=h.j5s77vdguv8h)
