# 05. Implementation Guide

visor 코드베이스의 구조, 핵심 API, 확장 방법을 설명합니다.

## 목차

1. [프로젝트 구조](#프로젝트-구조)
2. [데이터 흐름](#데이터-흐름)
3. [핵심 패키지](#핵심-패키지)
4. [Widget 시스템](#widget-시스템)
5. [새 위젯 추가하기](#새-위젯-추가하기)
6. [설정 시스템](#설정-시스템)
7. [테스트](#테스트)
8. [성능 고려사항](#성능-고려사항)

---

## 프로젝트 구조

```
visor/
├── cmd/visor/
│   └── main.go              # CLI 엔트리포인트
├── internal/
│   ├── input/               # stdin JSON 파싱
│   │   ├── session.go       # Session 구조체 정의
│   │   ├── reader.go        # JSON 파서
│   │   └── reader_test.go
│   ├── config/              # TOML 설정 관리
│   │   ├── types.go         # Config 구조체
│   │   ├── defaults.go      # 기본 설정값
│   │   └── loader.go        # 파일 로딩/저장
│   ├── widgets/             # 위젯 구현
│   │   ├── widget.go        # Widget 인터페이스 + Registry
│   │   ├── model.go         # 모델명 위젯
│   │   ├── context.go       # 컨텍스트 위젯
│   │   ├── context_spark.go # 컨텍스트 스파크라인 위젯 (v0.2)
│   │   ├── compact_eta.go   # Compact 예측 위젯 (v0.2)
│   │   ├── burn_rate.go     # 번 레이트 위젯 (v0.2)
│   │   ├── git.go           # Git 상태 위젯
│   │   ├── cost.go          # 비용 위젯
│   │   ├── cache_hit.go     # 캐시 히트율 위젯 (고유)
│   │   ├── api_latency.go   # API 지연시간 위젯 (고유)
│   │   ├── code_changes.go  # 코드 변경량 위젯 (고유)
│   │   └── *_test.go        # 위젯별 테스트
│   ├── render/              # 출력 렌더링
│   │   ├── ansi.go          # ANSI 컬러 코드
│   │   ├── truncate.go      # 문자열 자르기
│   │   └── layout.go        # 위젯 조합 + Split 레이아웃
│   ├── history/             # 세션 히스토리 (v0.2)
│   │   ├── history.go       # 히스토리 관리
│   │   └── history_test.go
│   └── git/                 # git CLI 래퍼
│       └── status.go        # git status 파싱
├── go.mod
├── go.sum
└── docs/                    # 설계 문서 (00_PRD ~ 06_PROGRESS)
```

### 디렉토리 역할

| 디렉토리 | 역할 | 의존성 |
|----------|------|--------|
| `cmd/visor` | CLI 진입점, 플래그 처리 | 모든 internal 패키지 |
| `internal/input` | JSON 파싱 | 없음 |
| `internal/config` | 설정 로딩 | BurntSushi/toml |
| `internal/widgets` | 위젯 로직 | input, config, render, git, history |
| `internal/render` | ANSI 출력 | 없음 |
| `internal/git` | git 명령 실행 | 없음 (외부 git 바이너리) |
| `internal/history` | 세션 히스토리 버퍼 | 없음 |

---

## 데이터 흐름

```
┌─────────────────────────────────────────────────────────────────┐
│                        main.go                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  input.Parse()  │  │  config.Load()  │  │  git.GetStatus()│
│  stdin → Session│  │  TOML → Config  │  │  CLI → Status   │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              ▼
                 ┌─────────────────────────┐
                 │  widgets.RenderAll()    │
                 │  Session + Config →     │
                 │  []string               │
                 └─────────────────────────┘
                              │
                              ▼
                 ┌─────────────────────────┐
                 │  render.MultiLine()     │
                 │  [][]string → ANSI      │
                 └─────────────────────────┘
                              │
                              ▼
                          stdout
```

### 단계별 설명

1. **입력 파싱** (`internal/input`)
   - `os.Stdin`에서 JSON 읽기
   - `Session` 구조체로 역직렬화
   - 파싱 실패 시 빈 `Session` 반환 (panic 방지)

2. **설정 로딩** (`internal/config`)
   - `~/.config/visor/config.toml` 읽기
   - 파일 없으면 기본 설정 사용
   - 위젯 순서, 스타일 정보 포함

3. **위젯 렌더링** (`internal/widgets`)
   - 설정된 순서대로 각 위젯 호출
   - `ShouldRender()` → `Render()` 순서
   - 빈 문자열 반환 시 스킵

4. **레이아웃 조합** (`internal/render`)
   - 위젯 출력을 공백으로 연결
   - 터미널 너비에 맞게 자르기
   - 멀티라인 지원

---

## 핵심 패키지

### internal/input

**Session 구조체** - Claude Code가 stdin으로 전달하는 JSON의 Go 표현:

```go
type Session struct {
    SessionID     string        `json:"session_id"`      // v0.2: 세션 식별자
    Model         Model         `json:"model"`
    Cost          Cost          `json:"cost"`
    ContextWindow ContextWindow `json:"context_window"`
    Workspace     Workspace     `json:"workspace"`
    CurrentUsage  *CurrentUsage `json:"current_usage"`   // nullable
}

type Model struct {
    DisplayName string `json:"display_name"`
    ID          string `json:"id"`
}

type Cost struct {
    TotalCostUSD          float64 `json:"total_cost_usd"`
    TotalDurationMs       int64   `json:"total_duration_ms"`     // v0.2: 전체 세션 시간
    TotalAPIDurationMs    int64   `json:"total_api_duration_ms"`
    TotalAPICalls         int     `json:"total_api_calls"`
    TotalInputTokens      int     `json:"total_input_tokens"`
    TotalOutputTokens     int     `json:"total_output_tokens"`
    TotalCacheReadTokens  int     `json:"total_cache_read_tokens"`
    TotalCacheWriteTokens int     `json:"total_cache_write_tokens"`
}

type ContextWindow struct {
    UsedPercentage float64 `json:"used_percentage"`
    UsedTokens     int     `json:"used_tokens"`
    MaxTokens      int     `json:"max_tokens"`
}

type Workspace struct {
    LinesAdded   int `json:"lines_added"`
    LinesRemoved int `json:"lines_removed"`
    FilesChanged int `json:"files_changed"`
}

type CurrentUsage struct {
    InputTokens     int `json:"input_tokens"`
    CacheReadTokens int `json:"cache_read_tokens"`
}
```

**Parse 함수**:

```go
func Parse(r io.Reader) *Session
```

- `io.Reader` 인터페이스 사용으로 테스트 용이
- JSON 파싱 실패 시 빈 `Session` 반환
- 절대 panic하지 않음

### internal/config

**Config 구조체**:

```go
type Config struct {
    General GeneralConfig `toml:"general"`
    Lines   []Line        `toml:"line"`
}

type GeneralConfig struct {
    Separator string `toml:"separator"`  // 위젯 간 구분자 (기본: " | ")
}

type Line struct {
    Widgets []WidgetConfig `toml:"widget"`  // 단일 레이아웃
    Left    []WidgetConfig `toml:"left"`    // v0.2: Split 레이아웃 (좌측)
    Right   []WidgetConfig `toml:"right"`   // v0.2: Split 레이아웃 (우측)
}

type WidgetConfig struct {
    Name   string            `toml:"name"`
    Format string            `toml:"format"`
    Style  StyleConfig       `toml:"style"`
    Extra  map[string]string `toml:"extra"`
}

type StyleConfig struct {
    Fg   string `toml:"fg"`
    Bg   string `toml:"bg"`
    Bold bool   `toml:"bold"`
}
```

**주요 함수**:

```go
func Load(path string) (*Config, error)     // 설정 로딩
func Init(path string) error                 // 기본 설정 생성
func Validate(path string) error             // 설정 검증
func DefaultConfigPath() string              // ~/.config/visor/config.toml
```

### internal/render

**ANSI 색상**:

```go
// 사용 가능한 색상명
var ColorMap = map[string]string{
    "black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
    "bright_black", "bright_red", "bright_green", "bright_yellow",
    "bright_blue", "bright_magenta", "bright_cyan", "bright_white",
    "gray", "grey",
}

func Colorize(text, fg string) string           // 단일 색상
func Style(text, fg, bg string, bold bool) string  // 복합 스타일
```

**문자열 처리**:

```go
func TerminalWidth() int                    // 터미널 너비 (기본 80)
func Truncate(s string, maxWidth int) string  // ANSI 인식 자르기
func VisibleLength(s string) int            // ANSI 제외 길이
```

**레이아웃**:

```go
func Layout(widgets []string, separator string) string                    // 단일 라인 조합
func SplitLayout(left, right []string, separator string) string           // v0.2: 좌/우 정렬
func MultiLine(lines [][]string, separator string) string                 // 멀티라인 조합
func JoinLines(lines []string) string                                     // v0.2: 라인 결합
```

### internal/git

```go
const commandTimeout = 200 * time.Millisecond  // Git 명령어 타임아웃

type Status struct {
    Branch   string
    IsRepo   bool
    IsDirty  bool
    Ahead    int
    Behind   int
    Staged   int
    Modified int
}

func GetStatus() Status
func gitCommand(args ...string) ([]byte, error)     // 타임아웃 적용된 git 실행
func gitCommandRun(args ...string) error            // 출력 없는 git 실행
```

- 외부 `git` 바이너리 호출 (200ms 타임아웃 적용)
- 비 git 디렉토리에서 빈 Status 반환
- 대형 저장소에서도 statusline 멈춤 방지

### internal/history (v0.2)

```go
const MaxEntries = 20  // 세션당 최대 히스토리 수

type Entry struct {
    Timestamp      int64   `json:"ts"`
    ContextPct     float64 `json:"ctx_pct"`
    CostUSD        float64 `json:"cost"`
    DurationMs     int64   `json:"dur_ms"`
    CacheHitPct    float64 `json:"cache_pct"`
    APILatencyMs   int64   `json:"api_ms"`
}

type History struct {
    SessionID string  `json:"session_id"`
    Entries   []Entry `json:"entries"`
}

func Load(sessionID string) (*History, error)   // 히스토리 로딩
func (h *History) Save() error                  // 히스토리 저장
func (h *History) Add(entry Entry)              // 엔트리 추가
func (h *History) GetContextHistory(n int) []float64  // 최근 n개 컨텍스트 값
func (h *History) Latest() *Entry               // 최신 엔트리
func (h *History) Count() int                   // 엔트리 수
```

- 히스토리 저장 경로: `~/.cache/visor/history_<session_id>.json`
- 세션별로 독립적인 히스토리 관리
- 최대 20개 엔트리 유지 (FIFO)

---

## Widget 시스템

### Widget 인터페이스

모든 위젯이 구현해야 하는 인터페이스:

```go
type Widget interface {
    // Name returns the widget identifier used in config
    Name() string

    // Render returns the formatted output string
    Render(session *input.Session, cfg *config.WidgetConfig) string

    // ShouldRender returns whether the widget should be displayed
    ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool
}
```

### Registry

위젯은 init()에서 자동 등록:

```go
var Registry = make(map[string]Widget)

func Register(w Widget)
func Get(name string) (Widget, bool)
func RenderAll(session *input.Session, widgets []config.WidgetConfig) []string
```

### 헬퍼 함수

**임계값 상수**:

```go
const (
    // Context window thresholds
    ContextWarningPct = 60.0   // Context 경고 임계값
    ContextDangerPct  = 80.0   // Context 위험 임계값

    // Cost thresholds (USD)
    CostWarningUSD    = 0.5    // 비용 경고 임계값
    CostDangerUSD     = 1.0    // 비용 위험 임계값

    // Cache hit rate thresholds (inverse: higher is better)
    CacheHitGoodPct    = 80.0  // 캐시 양호 임계값
    CacheHitWarningPct = 50.0  // 캐시 경고 임계값

    // API latency thresholds (ms)
    LatencyWarningMs  = 2000   // 지연시간 경고 임계값
    LatencyDangerMs   = 5000   // 지연시간 위험 임계값

    // v0.2: Burn rate thresholds (cents per minute)
    BurnRateWarningCents = 10.0  // 번 레이트 경고 (10¢/min)
    BurnRateDangerCents  = 25.0  // 번 레이트 위험 (25¢/min)

    // v0.2: Compact ETA thresholds (minutes)
    CompactETAWarningMin = 10.0  // Compact 예측 경고 (<10분)
    CompactETADangerMin  = 5.0   // Compact 예측 위험 (<5분)
    CompactThresholdPct  = 80.0  // Compact 트리거 임계값
)
```

**색상 결정 함수**:

```go
// 값이 클수록 나쁜 경우 (cost, latency, context)
func ColorByThreshold(value, warning, danger float64) string

// 값이 클수록 좋은 경우 (cache hit rate)
func ColorByThresholdInverse(value, good, warning float64) string
```

**포맷/옵션 함수**:

```go
// 커스텀 포맷 적용. {value} 플레이스홀더 지원.
func FormatOutput(cfg *config.WidgetConfig, defaultFormat, value string) string

// Extra 맵에서 값 조회
func GetExtra(cfg *config.WidgetConfig, key, defaultValue string) string

// Extra 맵에서 bool 값 조회
func GetExtraBool(cfg *config.WidgetConfig, key string, defaultValue bool) bool

// Extra 맵에서 int 값 조회
func GetExtraInt(cfg *config.WidgetConfig, key string, defaultValue int) int

// 렌더링
func RenderAll(session *input.Session, widgets []config.WidgetConfig) []string
```

**프로그레스 바**:

```go
const (
    BarFilled       = "█"
    BarEmpty        = "░"
    DefaultBarWidth = 10
)

// 프로그레스 바 문자열 생성
// 예: ProgressBar(42.0, 10) → "████░░░░░░"
func ProgressBar(pct float64, width int) string
```

### 기존 위젯 구현 패턴

```go
// internal/widgets/model.go
type ModelWidget struct{}

func (w *ModelWidget) Name() string {
    return "model"  // config.toml에서 사용하는 식별자
}

func (w *ModelWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
    name := session.Model.DisplayName
    if name == "" {
        return ""
    }
    return render.Colorize(name, "cyan")
}

func (w *ModelWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
    return session.Model.DisplayName != "" || session.Model.ID != ""
}
```

---

## 새 위젯 추가하기

### Step 1: 위젯 파일 생성

`internal/widgets/mywidget.go`:

```go
package widgets

import (
    "fmt"
    "github.com/namyoungkim/visor/internal/config"
    "github.com/namyoungkim/visor/internal/input"
    "github.com/namyoungkim/visor/internal/render"
)

type MyWidget struct{}

func (w *MyWidget) Name() string {
    return "my_widget"
}

func (w *MyWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
    // 데이터 추출
    value := session.SomeField

    // 조건부 색상
    color := "green"
    if value > threshold {
        color = "red"
    }

    // 포맷팅
    text := fmt.Sprintf("Label: %v", value)
    return render.Colorize(text, color)
}

func (w *MyWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
    return session.SomeField != nil  // 데이터가 있을 때만 표시
}
```

### Step 2: Registry에 등록

`internal/widgets/widget.go`의 init() 함수에 추가:

```go
func init() {
    Register(&ModelWidget{})
    Register(&ContextWidget{})
    // ... 기존 위젯들
    Register(&MyWidget{})  // 추가
}
```

### Step 3: 테스트 작성

`internal/widgets/mywidget_test.go`:

```go
package widgets

import (
    "testing"
    "github.com/namyoungkim/visor/internal/config"
    "github.com/namyoungkim/visor/internal/input"
)

func TestMyWidget_Render(t *testing.T) {
    w := &MyWidget{}
    session := &input.Session{
        // 테스트 데이터
    }

    result := w.Render(session, &config.WidgetConfig{})

    if !strings.Contains(result, "expected") {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}
```

### Step 4: 설정에서 사용

```toml
[[line]]
  [[line.widget]]
  name = "my_widget"
```

---

## 설정 시스템

### TOML 구조

```toml
# ~/.config/visor/config.toml

[general]
separator = " | "

# 단일 레이아웃 (기본)
[[line]]
  [[line.widget]]
  name = "model"

  [[line.widget]]
  name = "context"
  [line.widget.style]
  fg = "yellow"
  bold = true

# v0.2: Split 레이아웃 (좌/우 정렬)
[[line]]
  [[line.left]]
  name = "model"

  [[line.left]]
  name = "git"

  [[line.right]]
  name = "cost"

  [[line.right]]
  name = "cache_hit"
```

### v0.2 위젯 Extra 옵션

```toml
[[line.widget]]
name = "compact_eta"
[line.widget.extra]
show_when_above = "40"   # context 40% 이상일 때만 표시

[[line.widget]]
name = "context_spark"
[line.widget.extra]
width = "8"              # 스파크라인 너비

[[line.widget]]
name = "burn_rate"
[line.widget.extra]
show_label = "true"      # "Burn:" 접두사 표시
```

### 스타일 옵션

```toml
[[line.widget]]
name = "cost"
[line.widget.style]
fg = "green"      # 전경색
bg = "black"      # 배경색
bold = true       # 굵게
```

**사용 가능한 색상**:
- 기본: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- 밝은: `bright_black`, `bright_red`, `bright_green`, `bright_yellow`, `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`
- 별칭: `gray`, `grey` (= `bright_black`)

---

## 테스트

### 테스트 실행

```bash
# 모든 테스트
go test ./...

# 특정 패키지
go test ./internal/widgets/...

# 상세 출력
go test -v ./...

# 커버리지
go test -cover ./...
```

### 테스트 작성 가이드

```go
func TestWidgetName_Scenario(t *testing.T) {
    w := &SomeWidget{}
    session := &input.Session{
        // 최소한의 필요한 데이터만
    }

    result := w.Render(session, &config.WidgetConfig{})

    // ANSI 코드가 포함되므로 strings.Contains 사용
    if !strings.Contains(result, "expected") {
        t.Errorf("Expected to contain 'expected', got '%s'", result)
    }
}
```

### 수동 테스트

```bash
# 전체 JSON
echo '{"model":{"display_name":"Opus"},"context_window":{"used_percentage":42.5},"cost":{"total_cost_usd":0.15,"total_api_duration_ms":2500},"current_usage":{"input_tokens":100,"cache_read_tokens":400},"workspace":{"lines_added":25,"lines_removed":10}}' | ./visor

# 최소 JSON
echo '{}' | ./visor

# 잘못된 JSON (graceful fallback 확인)
echo 'invalid' | ./visor
```

---

## 성능 고려사항

### 목표

- Cold startup: < 5ms
- 메모리: 최소화 (statusline은 자주 호출됨)

### 최적화 포인트

1. **의존성 최소화**
   - 유일한 외부 의존성: `BurntSushi/toml`
   - 표준 라이브러리 우선 사용

2. **Lazy 로딩**
   - Git 상태는 git 위젯이 설정에 있을 때만 조회
   - 불필요한 계산 회피

3. **문자열 처리**
   - `strings.Builder` 대신 단순 연결 (짧은 문자열)
   - ANSI 코드 상수 사용

4. **에러 처리**
   - panic 대신 빈 값 반환
   - 로깅 없음 (stdout만 사용)

### 벤치마크

```bash
# 성능 측정
time (echo '{}' | ./visor)

# 반복 테스트
for i in {1..100}; do echo '{}' | ./visor > /dev/null; done
```

---

## 향후 개선 방향 (v0.3+)

1. **Transcript 파싱**
   - transcript JSONL 파싱
   - Tool/Agent 활동 추적
   - `tools`, `agents` 위젯

2. **설정 확장**
   - 위젯별 threshold 커스터마이징
   - 색상 테마 프리셋
   - Powerline 스타일 지원

3. **배포 자동화**
   - GitHub Actions 자동 릴리즈
   - goreleaser 통합

4. **성능 개선**
   - Git 정보 캐싱
   - 설정 파일 변경 감지
