# Build stage
FROM golang:1.24.10-alpine AS builder

# 빌드에 필요한 패키지 설치
RUN apk add --no-cache git gcc musl-dev

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사 및 의존성 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드 (CGO_ENABLED=0로 정적 링크)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o trading-bot ./cmd/server/main.go

# Runtime stage
FROM alpine:latest

# 타임존 및 CA 인증서 설치
RUN apk --no-cache add ca-certificates tzdata

# 타임존 설정 (한국 시간)
ENV TZ=Asia/Seoul

# 비-root 사용자 생성
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 작업 디렉토리 설정
WORKDIR /app

# 로그 디렉토리 생성 및 권한 설정
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

# 빌더에서 바이너리 복사
COPY --from=builder --chown=appuser:appuser /app/trading-bot .

# application.json 복사
COPY --chown=appuser:appuser application.json .

# 사용자 전환
USER appuser

# 포트 노출 (문서화 목적)
# 런타임에 .env의 PORT 환경 변수가 사용됨
EXPOSE 3000 5000 8080

# 애플리케이션 실행
CMD ["./trading-bot"]

