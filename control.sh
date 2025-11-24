#!/bin/bash

# Go Trading Bot 제어 스크립트

# 설정
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_NAME="go-trading-bot"
BINARY_PATH="${SCRIPT_DIR}/bin/${APP_NAME}"
PID_FILE="${SCRIPT_DIR}/tmp/${APP_NAME}.pid"
LOG_FILE="${SCRIPT_DIR}/logs/${APP_NAME}.log"
TMP_DIR="${SCRIPT_DIR}/tmp"
LOG_DIR="${SCRIPT_DIR}/logs"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 디렉토리 생성
mkdir -p "${TMP_DIR}"
mkdir -p "${LOG_DIR}"

# 함수: 프로세스 상태 확인
is_running() {
    if [ -f "${PID_FILE}" ]; then
        PID=$(cat "${PID_FILE}")
        if ps -p "${PID}" > /dev/null 2>&1; then
            return 0
        else
            # PID 파일은 있지만 프로세스가 없는 경우
            rm -f "${PID_FILE}"
            return 1
        fi
    fi
    return 1
}

# 함수: 빌드
build() {
    echo -e "${BLUE}Building ${APP_NAME}...${NC}"
    cd "${SCRIPT_DIR}" || exit 1
    
    # bin 디렉토리 생성
    mkdir -p bin
    
    # 빌드
    go build -o "${BINARY_PATH}" ./cmd/server/main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Build successful${NC}"
        return 0
    else
        echo -e "${RED}✗ Build failed${NC}"
        return 1
    fi
}

# 함수: 서비스 시작 (백그라운드)
start() {
    if is_running; then
        echo -e "${YELLOW}⚠ ${APP_NAME} is already running (PID: $(cat ${PID_FILE}))${NC}"
        return 1
    fi
    
    # 바이너리가 없으면 빌드
    if [ ! -f "${BINARY_PATH}" ]; then
        echo -e "${YELLOW}Binary not found. Building...${NC}"
        build || return 1
    fi
    
    echo -e "${BLUE}Starting ${APP_NAME}...${NC}"
    
    # 백그라운드로 실행하고 PID 저장
    nohup "${BINARY_PATH}" >> "${LOG_FILE}" 2>&1 &
    echo $! > "${PID_FILE}"
    
    # 시작 확인 (1초 대기)
    sleep 1
    
    if is_running; then
        echo -e "${GREEN}✓ ${APP_NAME} started successfully (PID: $(cat ${PID_FILE}))${NC}"
        echo -e "${BLUE}  Log file: ${LOG_FILE}${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to start ${APP_NAME}${NC}"
        echo -e "${YELLOW}  Check log file: ${LOG_FILE}${NC}"
        return 1
    fi
}

# 함수: 서비스 시작 (포그라운드)
foreground() {
    if is_running; then
        echo -e "${YELLOW}⚠ ${APP_NAME} is already running in background (PID: $(cat ${PID_FILE}))${NC}"
        echo -e "${YELLOW}  Stop it first using: $0 stop${NC}"
        return 1
    fi
    
    # 바이너리가 없으면 빌드
    if [ ! -f "${BINARY_PATH}" ]; then
        echo -e "${YELLOW}Binary not found. Building...${NC}"
        build || return 1
    fi
    
    echo -e "${BLUE}Starting ${APP_NAME} in foreground mode...${NC}"
    echo -e "${BLUE}Press Ctrl+C to stop${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # 포그라운드로 실행 (Ctrl+C로 종료)
    "${BINARY_PATH}"
    
    echo -e "\n${BLUE}${APP_NAME} stopped${NC}"
}

# 함수: 서비스 중지
stop() {
    if ! is_running; then
        echo -e "${YELLOW}⚠ ${APP_NAME} is not running${NC}"
        return 1
    fi
    
    PID=$(cat "${PID_FILE}")
    echo -e "${BLUE}Stopping ${APP_NAME} (PID: ${PID})...${NC}"
    
    # SIGTERM 전송
    kill -TERM "${PID}" 2>/dev/null
    
    # 최대 10초 대기
    for i in {1..10}; do
        if ! ps -p "${PID}" > /dev/null 2>&1; then
            rm -f "${PID_FILE}"
            echo -e "${GREEN}✓ ${APP_NAME} stopped successfully${NC}"
            return 0
        fi
        sleep 1
    done
    
    # 강제 종료
    echo -e "${YELLOW}Process did not stop gracefully. Force killing...${NC}"
    kill -9 "${PID}" 2>/dev/null
    rm -f "${PID_FILE}"
    
    if ps -p "${PID}" > /dev/null 2>&1; then
        echo -e "${RED}✗ Failed to stop ${APP_NAME}${NC}"
        return 1
    else
        echo -e "${GREEN}✓ ${APP_NAME} force stopped${NC}"
        return 0
    fi
}

# 함수: 서비스 재시작
restart() {
    echo -e "${BLUE}Restarting ${APP_NAME}...${NC}"
    stop
    sleep 2
    start
}

# 함수: 상태 확인
status() {
    if is_running; then
        PID=$(cat "${PID_FILE}")
        UPTIME=$(ps -p "${PID}" -o etime= | tr -d ' ')
        MEM=$(ps -p "${PID}" -o rss= | awk '{printf "%.2f MB", $1/1024}')
        CPU=$(ps -p "${PID}" -o %cpu= | tr -d ' ')
        
        echo -e "${GREEN}● ${APP_NAME} is running${NC}"
        echo -e "  PID: ${PID}"
        echo -e "  Uptime: ${UPTIME}"
        echo -e "  Memory: ${MEM}"
        echo -e "  CPU: ${CPU}%"
        echo -e "  Log: ${LOG_FILE}"
    else
        echo -e "${RED}● ${APP_NAME} is not running${NC}"
        return 1
    fi
}

# 함수: 로그 보기
logs() {
    if [ ! -f "${LOG_FILE}" ]; then
        echo -e "${YELLOW}⚠ Log file not found: ${LOG_FILE}${NC}"
        return 1
    fi
    
    LINES=${1:-50}
    echo -e "${BLUE}Showing last ${LINES} lines of log:${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    tail -n "${LINES}" "${LOG_FILE}"
}

# 함수: 로그 실시간 보기
follow() {
    if [ ! -f "${LOG_FILE}" ]; then
        echo -e "${YELLOW}⚠ Log file not found: ${LOG_FILE}${NC}"
        return 1
    fi
    
    echo -e "${BLUE}Following log file (Ctrl+C to stop):${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    tail -f "${LOG_FILE}"
}

# 함수: 로그 삭제
clean() {
    echo -e "${BLUE}Cleaning up...${NC}"
    
    if is_running; then
        echo -e "${YELLOW}⚠ Cannot clean while ${APP_NAME} is running${NC}"
        return 1
    fi
    
    rm -f "${LOG_FILE}"
    rm -f "${PID_FILE}"
    echo -e "${GREEN}✓ Cleaned up logs and pid files${NC}"
}

# 함수: 사용법 출력
usage() {
    echo "Usage: $0 {build|start|foreground|stop|restart|status|logs|follow|clean}"
    echo ""
    echo "Commands:"
    echo "  build      - Build the application"
    echo "  start      - Start the trading bot in background"
    echo "  foreground - Start the trading bot in foreground (Ctrl+C to stop)"
    echo "  stop       - Stop the trading bot"
    echo "  restart    - Restart the trading bot"
    echo "  status     - Show the trading bot status"
    echo "  logs       - Show last 50 lines of log (use 'logs N' for N lines)"
    echo "  follow     - Follow log file in real-time"
    echo "  clean      - Clean up logs and pid files"
    echo ""
    echo "Examples:"
    echo "  $0 start          # Start in background"
    echo "  $0 foreground     # Start in foreground for debugging"
    echo "  $0 logs 100       # Show last 100 lines"
    echo "  $0 restart        # Restart the service"
}

# 메인 로직
case "${1}" in
    build)
        build
        ;;
    start)
        start
        ;;
    foreground|fg|run)
        foreground
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs "${2}"
        ;;
    follow)
        follow
        ;;
    clean)
        clean
        ;;
    *)
        usage
        exit 1
        ;;
esac

exit $?

