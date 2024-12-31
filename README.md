# 编译命令

```sh
# 编译
docker build -t yarn-api-img .

# 运行
docker run -it -p 8080:8080  -d --name yarn-api yarn-api-img /bin/bash
```


