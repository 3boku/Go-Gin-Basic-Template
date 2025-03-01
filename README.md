# Go-Gin-Basic-Template
Go Gin을 사용하는 3 tier 아키텍처를 사용한 웹 어플리케이션 템플릿 입니다.

#  .env template
```dotenv
PORT=:8080

POSTGRES_HOST=
POSTGRES_USER=
POSTGRES_PASS=
POSTGRES_DB=
POSTGRES_PORT=
```

# Deploy
배포는 쿠버네티스 쓸려고 하는데 이건 각 프로젝트에서 직접 구현하는게 나을 거 같아용 </br>
하지만 쿠버네티스를 안쓰는 사람들도 있으니 docker-compose 파일은 추가합니다
```yaml
services:
    app:
        build:
            context: .
            dockerfile: Dockerfile
        env_file:
            - ./secret/.env
        ports:
            - "8080:8080"
```
- docker-compose.yaml