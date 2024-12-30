# docker基本知识

![docker基本知识](./docker-1.png)
```sh
# arg 通过命令改变运行时参数
docker build -t test --build-arg a=12
```

# 禁用 Docker BuildKit
# 禁用使用 Docker CLI 
 
export DOCKER_BUILDKIT=0
export COMPOSE_DOCKER_CLI_BUILD=0
