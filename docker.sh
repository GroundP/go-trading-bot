#!/bin/bash

# Go Trading Bot Docker 관리 스크립트

set -e

PROJECT_NAME="go-trading-bot"
COMPOSE_FILE="docker-compose.yml"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# .env 파일 확인
check_env_file() {
    if [ ! -f .env ]; then
        log_warning ".env 파일이 존재하지 않습니다."
        log_info ".env.example을 복사하여 .env 파일을 만듭니다..."
        cp .env.example .env
        log_warning ".env 파일을 열어서 ACCESS_KEY와 SECRET_KEY를 설정해주세요!"
        exit 1
    fi
    
    # API 키 확인
    if grep -q "your_upbit_access_key_here" .env || grep -q "your_upbit_secret_key_here" .env; then
        log_error "Upbit API 키가 설정되지 않았습니다!"
        log_info ".env 파일을 열어서 ACCESS_KEY와 SECRET_KEY를 설정해주세요."
        exit 1
    fi
}

# 명령어 처리
case "$1" in
    build)
        log_info "Docker 이미지를 빌드합니다..."
        docker-compose -f "$COMPOSE_FILE" build --no-cache
        log_success "빌드 완료!"
        ;;
    
    start)
        check_env_file
        log_info "Trading Bot을 시작합니다..."
        docker-compose -f "$COMPOSE_FILE" up -d
        log_success "Trading Bot이 백그라운드에서 실행 중입니다."
        log_info "로그 확인: ./docker.sh logs"
        ;;
    
    stop)
        log_info "Trading Bot을 중지합니다..."
        docker-compose -f "$COMPOSE_FILE" stop
        log_success "Trading Bot이 중지되었습니다."
        ;;
    
    restart)
        log_info "Trading Bot을 재시작합니다..."
        docker-compose -f "$COMPOSE_FILE" restart
        log_success "Trading Bot이 재시작되었습니다."
        ;;
    
    logs)
        log_info "로그를 확인합니다... (Ctrl+C로 종료)"
        docker-compose -f "$COMPOSE_FILE" logs -f --tail=100
        ;;
    
    status)
        log_info "컨테이너 상태:"
        docker-compose -f "$COMPOSE_FILE" ps
        ;;
    
    shell)
        log_info "컨테이너 내부로 접속합니다..."
        docker-compose -f "$COMPOSE_FILE" exec trading-bot sh
        ;;
    
    down)
        log_warning "모든 컨테이너를 중지하고 제거합니다..."
        docker-compose -f "$COMPOSE_FILE" down
        log_success "컨테이너가 제거되었습니다."
        ;;
    
    clean)
        log_warning "모든 컨테이너와 볼륨을 제거합니다..."
        read -p "정말로 모든 데이터를 삭제하시겠습니까? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose -f "$COMPOSE_FILE" down -v
            log_success "컨테이너와 볼륨이 제거되었습니다."
        else
            log_info "취소되었습니다."
        fi
        ;;
    
    rebuild)
        check_env_file
        log_info "이미지를 다시 빌드하고 재시작합니다..."
        docker-compose -f "$COMPOSE_FILE" down
        docker-compose -f "$COMPOSE_FILE" build --no-cache
        docker-compose -f "$COMPOSE_FILE" up -d
        log_success "재빌드 및 재시작 완료!"
        ;;
    
    stats)
        log_info "컨테이너 리소스 사용량:"
        docker stats $PROJECT_NAME --no-stream
        ;;
    
    api-test)
        log_info "API 헬스체크 테스트..."
        PORT=$(grep "^PORT=" .env | cut -d '=' -f2)
        PORT=${PORT:-5000}
        
        response=$(curl -s -w "\n%{http_code}" http://localhost:$PORT/api/v1/health)
        http_code=$(echo "$response" | tail -n 1)
        body=$(echo "$response" | head -n -1)
        
        if [ "$http_code" = "200" ]; then
            log_success "API 서버가 정상 작동 중입니다!"
            echo "Response: $body"
        else
            log_error "API 서버 응답 실패 (HTTP $http_code)"
            echo "Response: $body"
        fi
        ;;
    
    *)
        echo "Go Trading Bot - Docker 관리 스크립트"
        echo ""
        echo "사용법: $0 {command}"
        echo ""
        echo "명령어:"
        echo "  build      - Docker 이미지 빌드"
        echo "  start      - Trading Bot 시작 (백그라운드)"
        echo "  stop       - Trading Bot 중지"
        echo "  restart    - Trading Bot 재시작"
        echo "  logs       - 실시간 로그 확인"
        echo "  status     - 컨테이너 상태 확인"
        echo "  shell      - 컨테이너 내부 접속"
        echo "  down       - 컨테이너 중지 및 제거"
        echo "  clean      - 컨테이너 및 볼륨 제거"
        echo "  rebuild    - 이미지 재빌드 및 재시작"
        echo "  stats      - 리소스 사용량 확인"
        echo "  api-test   - API 서버 헬스체크"
        echo ""
        echo "예시:"
        echo "  $0 build          # 이미지 빌드"
        echo "  $0 start          # 봇 시작"
        echo "  $0 logs           # 로그 확인"
        echo "  $0 stop           # 봇 중지"
        exit 1
        ;;
esac

