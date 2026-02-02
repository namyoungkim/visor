# C4 Model

## C1: System Context

```
                              ┌──────────────────────┐
                              │     Claude Code      │
                              │    (Host Process)    │
                              └──────────┬───────────┘
                                         │ stdin JSON (~300ms 간격)
                                         ▼
┌───────────┐              ┌──────────────────────────┐              ┌───────────┐
│           │  settings    │                          │   stdout     │           │
│  User     │─────────────▶│         visor            │─────────────▶│ Terminal  │
│ (개발자)  │  .json 설정  │   (Go CLI Binary)        │  ANSI text   │  (하단)   │
│           │              │                          │              │           │
└───────────┘              └──────────┬───────────────┘              └───────────┘
                                      │
                           ┌──────────┼──────────┐
                           ▼          ▼          ▼
                    ┌──────────┐ ┌────────┐ ┌─────────┐
                    │  Git     │ │ Config │ │Transcript│
                    │Repository│ │ (TOML) │ │ (JSONL)  │
                    └──────────┘ └────────┘ └─────────┘
```

### External Entities

| Entity | 유형 | 상호작용 |
|--------|------|----------|
| **Claude Code** | Host Process | stdin으로 세션 JSON 전달, stdout으로 포맷된 문자열 수신 |
| **User (개발자)** | Person | 설정 파일 편집, 터미널 하단 statusline 확인 |
| **Terminal** | Display | ANSI escape sequence가 포함된 문자열을 렌더링 |
| **Git Repository** | Data Source | 브랜치, 변경사항 등 VCS 정보 제공 |
| **Config (TOML)** | Data Store | `~/.config/visor/config.toml` — 사용자 설정 |
| **Transcript (JSONL)** | Data Source | Claude Code 세션 로그 (Phase 3에서 활용) |

---

## C2: Container Diagram

visor는 단일 바이너리이므로 Container는 하나다.
내부를 논리적 레이어로 구분한다.

```
┌─────────────────────────────────────────────────────────────────────┐
│                           visor (Go Binary)                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌───────────┐    ┌───────────┐    ┌───────────┐    ┌───────────┐ │
│   │   main    │───▶│  config   │    │  widgets  │    │  render   │ │
│   │ (진입점)  │    │ (설정 로드)│    │ (데이터   │    │ (출력     │ │
│   │           │───▶│           │───▶│  가공)    │───▶│  포매팅)  │ │
│   └───────────┘    └───────────┘    └───────────┘    └───────────┘ │
│        │                                  │                         │
│        ▼                                  ▼                         │
│   ┌───────────┐                    ┌───────────┐                   │
│   │  input    │                    │   git     │                   │
│   │ (stdin    │                    │ (외부 명령│                   │
│   │  JSON)    │                    │  실행)    │                   │
│   └───────────┘                    └───────────┘                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Containers (논리 레이어)

| Layer | 기술 | 책임 |
|-------|------|------|
| **main** | Go main package | CLI 진입점, 플래그 파싱 (`--setup`, `--init`, `--version`) |
| **input** | Go `encoding/json` | stdin에서 JSON 읽기 → `Session` 구조체 파싱 |
| **config** | Go + TOML parser | TOML 설정 파일 로드, 기본값 머징, 유효성 검증 |
| **widgets** | Go interfaces | 각 위젯이 Session 데이터에서 표시 문자열 생성 |
| **render** | Go + ANSI codes | 위젯 출력 합성, 컬러링, 레이아웃, truncation |
| **git** | Go `os/exec` | `git` CLI 호출로 브랜치/상태 정보 수집 |

---

## C3: Component Diagram

### input 패키지

```
┌─────────────────────────────────────────────────────────────────┐
│                        input package                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────┐         ┌─────────────────┐              │
│   │    Reader        │────────▶│    Session       │              │
│   │                  │         │   (struct)       │              │
│   │ io.ReadAll(stdin)│         │                  │              │
│   │ json.Unmarshal   │         │ .Model           │              │
│   └─────────────────┘         │ .Workspace       │              │
│                                │ .Cost            │              │
│                                │ .ContextWindow   │              │
│                                │ .Version         │              │
│                                │ .TranscriptPath  │              │
│                                └─────────────────┘              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Component | 책임 |
|-----------|------|
| **Reader** | stdin에서 raw bytes 읽기 → JSON unmarshal |
| **Session** | 파싱된 세션 데이터 구조체, 누락 필드 zero value 보장 |

### config 패키지

```
┌─────────────────────────────────────────────────────────────────┐
│                       config package                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────┐    ┌─────────────────┐                   │
│   │    Loader        │───▶│    Config        │                   │
│   │                  │    │   (struct)       │                   │
│   │ TOML 파일 읽기   │    │                  │                   │
│   │ 기본값 머징      │    │ .Lines[]         │                   │
│   │ 유효성 검증      │    │ .General         │                   │
│   └─────────────────┘    │ .Widgets map     │                   │
│                           └─────────────────┘                   │
│                                                                  │
│   ┌─────────────────┐                                           │
│   │    Defaults      │                                           │
│   │                  │                                           │
│   │ 설정 없을 때     │                                           │
│   │ 기본 레이아웃    │                                           │
│   └─────────────────┘                                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Component | 책임 |
|-----------|------|
| **Loader** | `~/.config/visor/config.toml` 파일 읽기, 파싱, 기본값 머징 |
| **Config** | 파싱된 설정 구조체 (라인 구성, 위젯 옵션) |
| **Defaults** | 설정 파일 없을 때의 기본 레이아웃 정의 |

### widgets 패키지

```
┌─────────────────────────────────────────────────────────────────┐
│                       widgets package                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────────────────────────────────────────┐       │
│   │  Widget (interface)                                  │       │
│   │                                                      │       │
│   │  Name() string                                       │       │
│   │  Render(session, config) string                      │       │
│   │  ShouldRender(session, config) bool                  │       │
│   └─────────────────────────────────────────────────────┘       │
│        △         △         △         △         △         △      │
│        │         │         │         │         │         │      │
│   ┌────┴──┐ ┌───┴───┐ ┌──┴────┐ ┌──┴──┐ ┌───┴───┐ ┌───┴──┐  │
│   │ Model │ │Context│ │  Git  │ │Cost │ │CacheHit│ │Latency│  │
│   │Widget │ │Widget │ │Widget │ │Widg.│ │ Widget │ │Widget │  │
│   └───────┘ └───────┘ └───────┘ └─────┘ └───────┘ └──────┘  │
│                                                    ┌───────┐   │
│                                                    │CodeChg│   │
│                                                    │Widget │   │
│                                                    └───────┘   │
│                                                                  │
│   ┌─────────────────┐                                           │
│   │    Registry      │                                           │
│   │                  │                                           │
│   │ name → Widget    │                                           │
│   │ map 관리         │                                           │
│   └─────────────────┘                                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Component | 입력 필드 | 출력 예시 |
|-----------|----------|----------|
| **ModelWidget** | `model.display_name` | `Opus` |
| **ContextWidget** | `context_window.used_percentage` | `42% ████░░░░░░` |
| **GitWidget** | 외부: `git rev-parse`, `git status` | `main +3~2?1` |
| **CostWidget** | `cost.total_cost_usd` | `$0.48` |
| **CacheHitWidget** | `context_window.current_usage.*` | `Cache: 73%` |
| **LatencyWidget** | `cost.total_api_duration_ms` | `API: 1.2s` |
| **CodeChangesWidget** | `cost.total_lines_added/removed` | `+156/-23` |
| **Registry** | — | name→Widget 매핑, 설정에서 이름으로 위젯 조회 |

### render 패키지

```
┌─────────────────────────────────────────────────────────────────┐
│                       render package                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────┐    ┌─────────────────┐                   │
│   │    Layout        │───▶│    ANSI          │                   │
│   │                  │    │                  │                   │
│   │ 위젯 배치       │    │ 색상 코드 래핑   │                   │
│   │ separator 삽입   │    │ 볼드/밑줄       │                   │
│   │ 멀티라인 합성   │    │ 리셋            │                   │
│   └─────────────────┘    └─────────────────┘                   │
│                                                                  │
│   ┌─────────────────┐                                           │
│   │   Truncate       │                                           │
│   │                  │                                           │
│   │ 터미널 너비 감지 │                                           │
│   │ 초과 시 ...으로  │                                           │
│   │ 스마트 자름      │                                           │
│   └─────────────────┘                                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

| Component | 책임 |
|-----------|------|
| **Layout** | 위젯 출력을 separator로 연결, 멀티라인 합성, (Phase 2) 좌/우 split |
| **ANSI** | ANSI escape code 유틸리티: 색상, 볼드, 리셋 |
| **Truncate** | 터미널 너비 감지 + 초과 시 ellipsis 처리 |

---

## Data Flow

```
┌─────────────┐
│ Claude Code  │
│ (host)       │
└──────┬──────┘
       │ stdin: JSON bytes
       ▼
┌──────────────┐     ┌──────────────┐
│ input.Reader │────▶│input.Session │
│              │     │  (struct)    │
└──────────────┘     └──────┬───────┘
                            │
       ┌────────────────────┼────────────────────┐
       ▼                    ▼                    ▼
┌──────────────┐   ┌──────────────┐     ┌──────────────┐
│config.Loader │   │widget.Registry│     │  git (exec)  │
│ → Config     │   │ → []Widget   │     │ → GitInfo    │
└──────┬───────┘   └──────┬───────┘     └──────┬───────┘
       │                  │                     │
       ▼                  ▼                     │
┌──────────────────────────────────────────────────────┐
│                   render.Layout                       │
│                                                      │
│  for each line in config.Lines:                      │
│    for each widget_name in line.Widgets:             │
│      widget = registry.Get(widget_name)              │
│      if widget.ShouldRender(session, widgetCfg):     │
│        parts = append(parts, widget.Render(...))     │
│    output = join(parts, separator)                   │
│    truncate(output, terminalWidth)                   │
│                                                      │
└──────────────────────────┬───────────────────────────┘
                           │
                           ▼ stdout: ANSI string
                    ┌──────────────┐
                    │   Terminal    │
                    └──────────────┘
```

---

## Directory Structure

```
visor/
├── cmd/
│   └── visor/
│       └── main.go              # 진입점, CLI 플래그
├── internal/
│   ├── input/
│   │   ├── reader.go            # stdin 읽기
│   │   └── session.go           # Session 구조체 정의
│   ├── config/
│   │   ├── loader.go            # TOML 로드 + 머징
│   │   ├── types.go             # Config 구조체
│   │   └── defaults.go          # 기본값
│   ├── widgets/
│   │   ├── widget.go            # Widget 인터페이스 + Registry
│   │   ├── model.go
│   │   ├── context.go
│   │   ├── git.go
│   │   ├── cost.go
│   │   ├── cache_hit.go
│   │   ├── latency.go
│   │   └── code_changes.go
│   ├── render/
│   │   ├── layout.go            # 위젯 배치 + 합성
│   │   ├── ansi.go              # ANSI 유틸리티
│   │   └── truncate.go          # 너비 감지 + 자르기
│   └── git/
│       └── status.go            # git CLI 래퍼
├── config.example.toml          # 예시 설정 파일
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```
