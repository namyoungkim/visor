# visor

Claude Code용 효율성 대시보드. 캐시 히트율, API 지연시간, 비용 소모율 등 숨겨진 메트릭을 실시간으로 표시합니다.

```
Opus | Ctx: 42% ████░░░░░░ | Cache: 80% | API: 2.5s | $0.15 | +25/-10 | main ↑1
```

## 특징

- **숨겨진 메트릭 시각화**: 캐시 히트율, API 지연시간, 코드 변경량 등 Claude Code가 내부적으로 사용하지만 노출하지 않던 데이터
- **빠른 시작**: Go 기반으로 5ms 이내 cold startup
- **유연한 설정**: TOML 설정 파일과 TUI 편집기로 위젯 배치, 테마 커스터마이징
- **안정성**: 잘못된 입력에도 panic 없이 graceful fallback

## 설치

### 바이너리 다운로드 (권장)

Go 설치 없이 바로 사용할 수 있습니다.

```bash
# 1. 버전 설정 (https://github.com/namyoungkim/visor/releases 에서 최신 버전 확인)
VERSION=0.11.0

# 2. 플랫폼에 맞는 바이너리 다운로드
curl -sL "https://github.com/namyoungkim/visor/releases/download/v${VERSION}/visor_${VERSION}_darwin_arm64.tar.gz" | tar xz   # macOS Apple Silicon
curl -sL "https://github.com/namyoungkim/visor/releases/download/v${VERSION}/visor_${VERSION}_darwin_amd64.tar.gz" | tar xz   # macOS Intel
curl -sL "https://github.com/namyoungkim/visor/releases/download/v${VERSION}/visor_${VERSION}_linux_amd64.tar.gz" | tar xz    # Linux x64
curl -sL "https://github.com/namyoungkim/visor/releases/download/v${VERSION}/visor_${VERSION}_linux_arm64.tar.gz" | tar xz    # Linux ARM64

# 3. PATH에 설치
sudo mv visor /usr/local/bin/

# sudo 권한이 없다면:
mkdir -p ~/.local/bin && mv visor ~/.local/bin/
# ~/.local/bin이 PATH에 없다면 쉘 설정에 추가: export PATH="$HOME/.local/bin:$PATH"
```

### Go install

Go 1.22 이상이 설치되어 있다면:

```bash
go install github.com/namyoungkim/visor@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/namyoungkim/visor.git
cd visor
go build -o visor ./cmd/visor
```

## 빠른 시작

### 1. Claude Code에 연결

`~/.claude/settings.json`에 추가:

```json
{
  "statusline": {
    "command": "visor"
  }
}
```

또는 환경 변수로:

```bash
export CLAUDE_STATUSLINE_COMMAND="visor"
```

### 2. 설정 초기화

```bash
visor --init          # 기본 설정 생성
visor --init minimal  # 최소 설정 (4개 위젯)
visor --init help     # 프리셋 목록 보기
```

### 3. 설정 편집 (선택)

```bash
visor --tui  # 인터랙티브 설정 편집기
```

## 위젯

| 위젯 | 식별자 | 설명 | 예시 |
|------|--------|------|------|
| 모델명 | `model` | 현재 사용 중인 모델 | `Opus` |
| 컨텍스트 | `context` | 컨텍스트 윈도우 사용률 | `Ctx: 42% ████░░░░░░` |
| 캐시 히트율 | `cache_hit` | 캐시에서 읽은 토큰 비율 | `Cache: 80%` |
| API 지연시간 | `api_latency` | API 호출 응답 시간 | `API: 2.5s` |
| 비용 | `cost` | 세션 누적 비용 | `$0.15` |
| 코드 변경량 | `code_changes` | 추가/삭제된 라인 수 | `+25/-10` |
| Git | `git` | 브랜치와 상태 | `main ↑1` |
| 비용 소모율 | `burn_rate` | 분당 비용 | `64.0¢/min` |
| 컨텍스트 예측 | `compact_eta` | 80% 도달 예상 시간 | `~18m` |
| 컨텍스트 추이 | `context_spark` | 사용률 변화 그래프 | `▂▃▄▅▆` |
| 도구 상태 | `tools` | 최근 도구 호출 | `✓Read ✓Write ◐Bash` |
| 에이전트 상태 | `agents` | 서브 에이전트 상태 | `✓Plan ◐Explore` |
| 일별 비용 | `daily_cost` | 오늘 누적 비용 | `$2.34 today` |
| 주별 비용 | `weekly_cost` | 이번 주 누적 비용 | `$15.67 week` |
| 블록 비용 | `block_cost` | 5시간 블록 비용 | `$0.45 block` |
| 5시간 제한 | `block_limit` | 5시간 블록 사용률 | `5h: 42%` |
| 7일 제한 | `week_limit` | 주간 사용률 | `7d: 69%` |
| 세션 ID | `session_id` | 현재 세션 ID | `abc123de` |
| 세션 시간 | `duration` | 세션 경과 시간 | `⏱️ 5m` |
| 토큰 속도 | `token_speed` | 출력 토큰 생성 속도 | `42.1 tok/s` |
| 요금제 | `plan` | 구독/API 타입 | `Pro` |
| 작업 진행 | `todos` | 작업 진행 상황 | `⊙ Task (3/5)` |
| 설정 현황 | `config_counts` | Claude 설정 현황 | `2📄 3🔒 2🔌 1🪝` |

### 핵심 메트릭 해석

**캐시 히트율** — 높을수록 비용 효율적
```
cache_read_input_tokens / (cache_read_input_tokens + input_tokens) × 100
```
- 80% 이상: 🟢 효율적
- 50~80%: 🟡 보통
- 50% 미만: 🔴 비효율적

**API 지연시간** — 콜당 평균 응답 속도 지표
- 2초 미만: 🟢 빠름
- 2~5초: 🟡 보통
- 5초 초과: 🔴 느림

**코드 변경량** — 세션 중 변경된 코드
- 🟢 추가된 라인 (+)
- 🔴 삭제된 라인 (-)

> 전체 위젯 설명은 [위젯 레퍼런스](docs/08_WIDGET_REFERENCE.md) 참조

## 설정

설정 파일: `~/.config/visor/config.toml`

### 프리셋

| 프리셋 | 용도 | 위젯 수 |
|--------|------|---------|
| `minimal` | 필수 정보만 | 4개 |
| `default` | 균형 잡힌 기본값 | 6개 |
| `efficiency` | 비용 최적화 | 6개 |
| `developer` | 도구/에이전트 모니터링 | 7개 |
| `pro` | Claude Pro 사용량 추적 | 6개 |
| `full` | 모든 위젯 (멀티라인) | 22개 |

```bash
visor --init efficiency  # 원하는 프리셋으로 초기화
```

### 설정 예시

```toml
[general]
separator = " | "

[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"
  [line.widget.extra]
  show_bar = "true"
  bar_width = "10"

  [[line.widget]]
  name = "cost"
```

### 위젯 옵션

| 위젯 | 옵션 | 기본값 | 설명 |
|------|------|--------|------|
| `context` | `show_bar` | `true` | 프로그레스 바 표시 |
| `context` | `bar_width` | `10` | 바 너비 |
| `cache_hit` | `show_label` | `true` | "Cache:" 라벨 표시 |
| `cost` | `show_label` | `false` | "Cost:" 라벨 표시 |
| `block_limit` | `show_remaining` | `true` | 남은 시간 표시 |

## 테마

| 테마 | 설명 |
|------|------|
| `default` | 기본 ASCII |
| `powerline` | Powerline 글리프 |
| `gruvbox` | Gruvbox 색상 |
| `nord` | Nord 색상 |
| `gruvbox-powerline` | Gruvbox + Powerline |
| `nord-powerline` | Nord + Powerline |

```toml
[theme]
name = "gruvbox"
powerline = true

# 색상 커스터마이징 (선택)
[theme.colors]
warning = "#ff00ff"
critical = "red"
```

## TUI 편집기

```bash
visor --tui
```

| 키 | 동작 |
|----|------|
| `j/k` | 이동 |
| `a/d` | 위젯 추가/삭제 |
| `J/K` | 순서 변경 |
| `e` | 옵션 편집 |
| `t` | 테마 변경 |
| `s` | 저장 |
| `q` | 종료 |

## CLI 옵션

```bash
visor --version   # 버전 확인
visor --init      # 설정 파일 생성
visor --setup     # Claude Code 연동 가이드
visor --check     # 설정 유효성 검사
visor --tui       # 설정 편집기
visor --debug     # 디버그 모드
```

## 요구사항

- **실행**: 별도 의존성 없음 (바이너리 설치 시)
- **빌드**: Go 1.22 이상
- **Git 위젯**: git CLI

## 라이선스

MIT License
