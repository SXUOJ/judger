# go-sandbox

<p align="center">
<a href="https://sxu.sx.cn"><img alt="Doc" src="https://img.shields.io/badge/doc-%E6%96%87%E6%A1%A3%E5%9C%B0%E5%9D%80-green?style=for-the-badge"></a>
<img alt="GitHub last commit (branch)" src="https://img.shields.io/github/last-commit/isther/judger/main?style=for-the-badge">
<img alt="GitHub Workflow Status (branch)" src="https://img.shields.io/github/workflow/status/isther/judger/Docker-Image%20CI/main?style=for-the-badge">
<img src="https://img.shields.io/github/languages/top/isther/judger?style=for-the-badge">
</p>

**国内镜像地址： hkccr.ccs.tencentyun.com/sxu-oj/judger**

## Docker 部署方式

```bash
git clone https://github.com/isther/judger.git
```

### 使用已构建的镜像

```bash
# docker-compose.yml
version: "3"
services:
  judger:
    container_name: judger
    image: hkccr.ccs.tencentyun.com/sxu-oj/judger
    restart: always
    tmpfs:
      - /tmp
    volumes:
      - $PWD/log:/sxu-judger/log
    ports:
      - 9000:9000
```

### 自行构建

```bash
# docker-compose.yml
version: "3"
services:
  judger:
    container_name: judger
    build:
      context: .
      dockerfile: Dockerfile
    restart: always
    tmpfs:
      - /tmp
    volumes:
      - $PWD/log:/sxu-judger/log
    ports:
      - 9000:9000
```
