# visor

Claude Code용 효율성 대시보드 statusline. 캐시 히트율, API 지연시간, 코드 변경량 등 다른 statusline에서 제공하지 않는 고유 메트릭을 실시간으로 표시합니다.

```
Opus | Ctx: 42% ████░░░░░░ | Cache: 80% | API: 2.5s | $0.15 | +25/-10 | main ↑1
```

## 특징

- **고유 메트릭**: 캐시 히트율, API 지연시간, 코드 변경량 - stdin JSON에 있지만 아무도 활용하지 않던 데이터
- **빠른 시작**: Go로 작성되어 < 5ms cold startup
- **설정 가능**: TOML 설정으로 위젯 순서, 색상 커스터마이징
- **안정성**: 잘못된 JSON에서도 panic 없이 graceful fallback

## 설치

### Go install (권장)

```bash
go install github.com/namyoungkim/visor@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/namyoungkim/visor.git
cd visor
go build -o visor ./cmd/visor
```

## Claude Code 연동

### 방법 1: settings.json

`~/.claude/settings.json`에 추가:

```json
{
  "statusline": {
    "command": "visor"
  }
}
```

### 방법 2: 환경 변수

```bash
export CLAUDE_STATUSLINE_COMMAND="visor"
```

## 위젯

| 위젯 | 식별자 | 설명 | 예시 |
|------|--------|------|------|
| 모델명 | `model` | 현재 사용 중인 모델 | `Opus` |
| 컨텍스트 | `context` | 컨텍스트 윈도우 사용률 + 프로그레스 바 | `Ctx: 42% ████░░░░░░` |
| 캐시 히트율 | `cache_hit` | 캐시에서 읽은 토큰 비율 | `Cache: 80%` |
| API 지연시간 | `api_latency` | 총 API 호출 시간 | `API: 2.5s` |
| 비용 | `cost` | 세션 총 비용 | `$0.15` |
| 코드 변경 | `code_changes` | 추가/삭제된 라인 수 | `+25/-10` |
| Git | `git` | 브랜치, 상태 | `main ↑1` |
| 번 레이트 | `burn_rate` | 분당 비용 소모율 | `64.0¢/min` |
| Compact 예측 | `compact_eta` | 80% context 도달 예측 | `~18m` |
| Context 스파크라인 | `context_spark` | 히스토리 기반 미니 그래프 | `▂▃▄▅▆` |
| 도구 상태 | `tools` | 최근 도구 호출 상태 | `✓Read ✓Write ◐Bash` |
| 에이전트 상태 | `agents` | 서브 에이전트 상태 | `✓Plan ◐Explore` |
| 일별 비용 | `daily_cost` | 오늘 누적 비용 | `$2.34 today` |
| 주별 비용 | `weekly_cost` | 이번 주 누적 비용 | `$15.67 week` |
| 블록 비용 | `block_cost` | 5시간 블록 비용 | `$0.45 block` |
| 5시간 제한 | `block_limit` | 5시간 블록 사용률 | `5h: 42%` |
| 7일 제한 | `week_limit` | 주간 사용률 | `7d: 69%` |

### 고유 메트릭 상세

**Cache Hit Rate** - 캐시 효율성 지표
```
rate = cache_read_tokens / (cache_read_tokens + input_tokens) × 100
```
- 80% 이상: 초록색 (효율적)
- 50-80%: 노란색 (보통)
- 50% 미만: 빨간색 (비효율적)

**API Latency** - 응답 시간 모니터링
- < 2초: 초록색
- 2-5초: 노란색
- > 5초: 빨간색

**Code Changes** - 세션 중 코드 변경량
- 초록색: 추가된 라인 (+)
- 빨간색: 삭제된 라인 (-)

## 설정

### 기본 설정 생성

```bash
visor --init
```

`~/.config/visor/config.toml` 생성:

```toml
[general]
separator = " | "  # 위젯 간 구분자 (기본값)

[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"

  [[line.widget]]
  name = "cache_hit"

  [[line.widget]]
  name = "api_latency"

  [[line.widget]]
  name = "cost"

  [[line.widget]]
  name = "code_changes"

  [[line.widget]]
  name = "git"
```

### 위젯 순서 변경

원하는 순서로 위젯 재배열:

```toml
[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "cost"

  [[line.widget]]
  name = "cache_hit"
```

### 멀티라인 설정

```toml
[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"

[[line]]
  [[line.widget]]
  name = "cache_hit"

  [[line.widget]]
  name = "api_latency"
```

### 위젯 커스터마이징

`format` 필드로 출력 포맷을 변경할 수 있습니다:

```toml
[[line.widget]]
name = "context"
format = "Context: {value}"  # "Ctx: 42%" 대신 "Context: 42%"
```

`extra` 필드로 위젯별 옵션을 설정할 수 있습니다:

```toml
[[line.widget]]
name = "context"
[line.widget.extra]
show_label = "false"  # "Ctx:" 접두사 숨기기 → "42%"만 표시

[[line.widget]]
name = "cost"
[line.widget.extra]
show_label = "true"   # "Cost:" 접두사 표시 → "Cost: $0.15"
```

**지원되는 extra 옵션**:

| 위젯 | 옵션 | 기본값 | 설명 |
|------|------|--------|------|
| `context` | `show_label` | `true` | "Ctx:" 접두사 표시 |
| `context` | `show_bar` | `true` | 프로그레스 바 표시 |
| `context` | `bar_width` | `10` | 프로그레스 바 너비 |
| `cache_hit` | `show_label` | `true` | "Cache:" 접두사 표시 |
| `cost` | `show_label` | `false` | "Cost:" 접두사 표시 |

### 구분자 설정

위젯 간 구분자를 변경할 수 있습니다:

```toml
[general]
separator = " :: "  # 기본값: " | "
```

출력 예시:
- `" | "` → `Opus | Ctx: 42% | $0.15`
- `" :: "` → `Opus :: Ctx: 42% :: $0.15`
- `" "` → `Opus Ctx: 42% $0.15`

## 테마

visor는 여러 테마 프리셋을 지원합니다:

| 테마 | 설명 |
|------|------|
| `default` | 기본 ASCII 구분자 |
| `powerline` | Powerline 글리프 (, ) |
| `gruvbox` | Gruvbox 색상 팔레트 |
| `nord` | Nord 색상 팔레트 |
| `gruvbox-powerline` | Gruvbox + Powerline |
| `nord-powerline` | Nord + Powerline |

TUI에서 `t` 키로 테마를 변경할 수 있습니다.

## TUI 설정 편집기

인터랙티브 TUI로 설정을 편집할 수 있습니다:

```bash
visor --tui
```

**주요 키바인딩:**

| 키 | 동작 |
|----|------|
| `j/k` | 커서 이동 |
| `a` | 위젯 추가 |
| `d` | 위젯 삭제 |
| `e` | 옵션 편집 |
| `J/K` | 위젯 순서 변경 |
| `L` | 레이아웃 변경 (single/split) |
| `t` | 테마 변경 |
| `s` | 저장 |
| `q` | 종료 |

**기능:**
- 위젯 추가/삭제/순서변경
- 위젯별 옵션 편집 (threshold 등)
- 레이아웃 변경 (single/split)
- 실시간 미리보기

## CLI 옵션

```bash
visor --version   # 버전 출력
visor --init      # 기본 설정 파일 생성
visor --setup     # Claude Code 연동 가이드
visor --check     # 설정 파일 유효성 검사
visor --debug     # 디버그 정보 출력 (stderr)
visor --tui       # 인터랙티브 설정 편집기
```

## 수동 테스트

```bash
echo '{"model":{"display_name":"Opus"},"context_window":{"used_percentage":42.5}}' | visor
```

## 요구사항

- Go 1.22 이상 (빌드 시)
- git (git 위젯 사용 시)

## 라이선스

MIT License
