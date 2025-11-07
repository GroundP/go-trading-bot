# Go Trading Bot

Go 1.24.10 기반 암호화폐 거래 봇 프로젝트

## 환경 구성

### Go 버전
- **Go 1.24.10** (linux/amd64)

### 환경 변수
프로젝트를 사용하려면 다음 환경 변수가 설정되어 있어야 합니다 (~/.bashrc에 추가됨):

```bash
export GOROOT=~/go1.24/go
export PATH=~/go1.24/go/bin:$PATH
export GOPATH=~/go
export GOBIN=$GOPATH/bin
```

### 환경 변수 적용
새 터미널을 열거나 다음 명령어를 실행하세요:

```bash
source ~/.bashrc
```

## 프로젝트 구조

```
go-trading-bot/
├── cmd/
│   └── bot/              # 메인 애플리케이션 진입점
├── internal/
│   ├── config/           # 설정 관리
│   ├── exchange/         # 거래소 API 연동
│   └── strategy/         # 거래 전략
├── pkg/
│   └── utils/            # 공통 유틸리티
├── go.mod
├── main.go
└── README.md
```

## 빌드 및 실행

### 프로그램 실행
```bash
go run main.go
```

### 빌드
```bash
go build -o trading-bot main.go
```

### 실행 결과 예시
```
=== Go Trading Bot ===
Go Version: go1.24.10
OS/Arch: linux/amd64
Environment configured successfully!
```

## 개발 가이드

### 의존성 추가
```bash
go get <package-name>
```

### 의존성 정리
```bash
go mod tidy
```

### 테스트
```bash
go test ./...
```

## 기술 스택
- **언어**: Go 1.24.10
- **아키텍처**: 비동기/병렬 프로그래밍
- **분야**: 암호화폐 거래, 외환거래(FX)

## 주요 기능 (예정)
- [ ] 다중 거래소 API 연동
- [ ] 실시간 시장 데이터 수집
- [ ] 자동 거래 전략 실행
- [ ] 포트폴리오 관리
- [ ] 백테스팅 기능
- [ ] 리스크 관리

## 참고사항
- WSL2 환경에서 개발됨
- 비동기 프로그래밍을 통한 고성능 처리
- Go의 goroutine과 channel을 활용한 병렬 처리

