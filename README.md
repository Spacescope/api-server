# FVM-Explorer-API-Server
### Regenerate swagger doc
swagger doc defined in router api comment.
if edited these comments, need to regenerate swagger doc.
```shell script
swag init -g cmd/api-server/main.go
```

### Swagger doc
swagger doc please refer to
`http://127.0.0.1:7006/api-server/swagger/index.html`

### How to make
```
make # make to see help
```
### Run
1. Config -- data-infra-backend/config
    ```
    task DB
    api server DB
    ```
2. Run
```
docker run -v /home/ec2-user/api-server/config:/etc/api-server/conf -p 7006:7006 -d 129862287110.dkr.ecr.us-east-2.amazonaws.com/extraction/api-server:commitId
```

### Refer
1. https://drive.google.com/drive/u/0/folders/1ptiBCy4lsO78KJqQR3oYv2TXrk3BrH8p
