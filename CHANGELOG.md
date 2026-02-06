# Changelog

모든 주요 변경사항을 이 파일에 기록합니다.

형식은 [Keep a Changelog](https://keepachangelog.com/ko/1.0.0/)를 따르며,
버전은 [Semantic Versioning](https://semver.org/lang/ko/)을 따릅니다.

## [Unreleased]

## [0.11.0] - 2026-02-06

### Added

- **로컬 사용량 추정** - API 호출 없이 JSONL 트랜스크립트 기반 사용량 추정
  - `internal/usage/local.go` - 로컬 추정 로직
  - `block_limit`/`week_limit` 위젯에서 OAuth 실패 시 fallback으로 사용
  - `EstimateLimits()` 함수로 tier별 한도 추정
- **Credential 파싱 개선**
  - nested envelope 구조 지원 (이중 JSON 래핑)
  - multi-keychain 탐색 (여러 Keychain 항목 시도)
  - `internal/auth/credentials_test.go` - credential 파싱 테스트 추가
- **per-call API 지연시간** - `api_latency` 위젯이 총 지연시간 대신 콜당 평균 표시
- **plan 위젯 구독 감지** - OAuth credential 기반 구독 타입 자동 감지
- **`cost.BlockDuration` 상수** - 5시간 블록 기간 상수 export
- **`cost.StartOfWeek` export** - 주간 시작 시점 계산 함수 export

### Changed

- `full` 프리셋: 22 → 24 위젯 (모든 위젯 포함, 7개 라인)
- `EstimateLimits()` 시그니처에 tier 파라미터 추가

### Fixed

- API key 사용자에 잘못된 Pro 한도 적용 방지
- `startOfWeek` 함수 중복 제거 (`cost.StartOfWeek`으로 통합)

## [0.10.0] - 2026-02-05

### Added

- **세션 정보 위젯 6종**
  - `session_id` - 현재 세션 ID (truncated, 기본 8자)
  - `duration` - 세션 경과 시간 (`⏱️ 5m`, `1h23m`)
  - `token_speed` - 출력 토큰 생성 속도 (`42.1 tok/s`)
  - `plan` - 요금제 타입 (`Pro`, `API`, `Bedrock`)
  - `todos` - TaskCreate/TaskUpdate 작업 진행 (`⊙ Task (3/5)`)
  - `config_counts` - Claude 설정 현황 (`2📄 3🔒 2🔌 1🪝`)
- **Transcript 파서 개선**
  - TaskCreate/TaskUpdate 파싱으로 Todo 추적
  - ISO 8601 타임스탬프 파싱 (time.RFC3339Nano)
  - 디버그 출력 (`config.Debug` 연동)
  - `VISOR_TRANSCRIPT_MAX_LINES` 환경변수로 파싱 라인 수 오버라이드
- **internal/claudeconfig/** - Claude 설정 파일 파싱 패키지
  - CLAUDE.md 파일 카운트 (cwd → root)
  - ~/.claude/settings.json 파싱 (rules, MCPs, hooks)

### Changed

- `developer` 프리셋: 6 → 7 위젯 (`todos` 추가)
- `full` 프리셋: 18 → 22 위젯 (모든 새 위젯 추가)
- `tools`/`agents` 위젯 `max_display` 기본값: 3 → 0 (무제한)
- Transcript 파싱 라인 수: 100 → 500

## [0.9.0] - 2026-02-04

### Added

- **설정 프리셋 시스템** - `visor --init <preset>` 명령으로 용도별 설정 생성
  - `minimal`: 필수 4개 위젯 (model, context, cost, git)
  - `default`: 균형 잡힌 6개 위젯 (visor 고유 메트릭 포함)
  - `efficiency`: 비용 최적화 6개 위젯
  - `developer`: 도구/에이전트 모니터링 6개 위젯
  - `pro`: Claude Pro 사용량 제한 6개 위젯
  - `full`: 모든 18개 위젯 (5개 라인)
- **`visor --init help`** - 사용 가능한 프리셋 목록 표시
- **`block_limit` 위젯 프로그레스 바 옵션**
  - `show_bar`: 프로그레스 바 표시 (기본: false)
  - `bar_width`: 프로그레스 바 너비 (기본: 10)

### Changed

- `Init()` 함수가 `InitWithPreset("default", path)` 위임으로 변경 (deprecated)
- `InitWithPreset()` 함수 추가 - 프리셋 기반 설정 생성

## [0.8.0] - 2026-02-04

### Added

- **커스텀 테마 오버라이드** - 프리셋 테마의 색상/구분자를 개별 커스터마이징
  - `[theme.colors]` - 색상 오버라이드 (normal, warning, critical, good, primary, secondary, muted, backgrounds)
  - `[theme.separators]` - 구분자 오버라이드 (left, right, left_soft, right_soft, left_hard, right_hard)
  - `theme.powerline = true` - 프리셋에 Powerline 스타일 적용
- **테마 해석 함수** - `theme.Resolve()` 프리셋 + 오버라이드 병합
- **색상 검증** - `theme.ValidateColor()` hex/named color 검증
- **config.Validate() 색상 검증** - `--check` 시 잘못된 색상 에러 반환
- 잘못된 테마 프리셋명 처리 - default로 자동 폴백

### Changed

- `config.ThemeConfig` 구조체에 `Colors`, `Separators`, `Powerline` 필드 추가
- `config.Validate()` 함수에서 테마 색상 유효성 검증

## [0.7.0] - 2026-02-04

### Added

- **도구 사용 횟수** - `tools` 위젯에 호출 횟수 표시 (`✓Bash ×7 | ✓Edit ×4`)
  - 같은 이름의 도구를 그룹화하여 Count 표시
  - `show_count` 옵션으로 횟수 표시 on/off (기본: true)
- **에이전트 상세 정보** - `agents` 위젯에 description과 실행 시간 표시
  - Task description 표시 (`Explore: Analyze widgets`)
  - 실행 시간 표시 (`(42s)`, `(2m)`, `(1h5m)`)
  - Running 상태는 실시간 경과시간 표시 (`(42s...)`)
  - `show_description`, `show_duration`, `max_description_len` 옵션

### Changed

- `Tool` 구조체에 `Count` 필드 추가
- `Agent` 구조체에 `Description`, `StartTime`, `EndTime` 필드 추가
- Parser가 도구를 ID 대신 Name으로 그룹화

### Fixed

- `truncateString` 함수가 멀티바이트 문자를 올바르게 처리하도록 수정

## [0.6.0] - 2026-02-03

### Added

- **테마 시스템** - Powerline 및 색상 테마 지원
  - 프리셋 테마: `default`, `powerline`, `gruvbox`, `nord`, `gruvbox-powerline`, `nord-powerline`
  - Powerline 글리프 구분자 지원 (`, `)
  - Hex 색상 코드 지원 (`#RRGGBB`)
  - TUI 테마 피커 (`t` 키)
- **누적 비용 추적** - JSONL 트랜스크립트 파싱으로 비용 집계
  - `daily_cost` 위젯 - 오늘 누적 비용
  - `weekly_cost` 위젯 - 이번 주 누적 비용
  - `block_cost` 위젯 - 5시간 블록 비용
  - Provider별 가격 적용 (Anthropic/Vertex/Bedrock)
  - 증분 파싱 캐시 시스템
- **사용량 제한 위젯** - Claude Pro OAuth API 연동
  - `block_limit` 위젯 - 5시간 블록 사용률 (`5h: 42%`)
  - `week_limit` 위젯 - 7일 사용률 (`7d: 69%`)
  - macOS Keychain credential provider
- **internal/theme/** - 테마 관리 패키지
- **internal/cost/** - JSONL 파싱 및 비용 집계 패키지
- **internal/auth/** - OAuth credential provider 패키지
- **internal/usage/** - 사용량 API 클라이언트 패키지

### Changed

- TUI에 테마 피커 추가 (`t` 키)
- render 패키지에 Powerline 레이아웃 및 Hex 색상 지원 추가
- 기본 위젯 11개 → 17개 (cost/usage 위젯 6개 추가)

### Fixed

- 테마 리스트 정렬 (TUI에서 일관된 순서)
- Go 1.22+ builtin min 함수 사용 (중복 제거)
- 캐시 히트 시 파싱 스킵으로 성능 개선

## [0.5.0] - 2026-02-03

### Added

- **TUI 설정 편집기** - `visor --tui`로 인터랙티브 설정 편집
  - Charm 생태계 사용 (bubbletea, bubbles, lipgloss)
  - 위젯 추가/삭제/순서변경
  - 위젯별 옵션 편집 (threshold 등)
  - 레이아웃 변경 (single/split)
  - 실시간 미리보기
  - Vim 스타일 키바인딩 (`j/k`, `J/K`, `a`, `d`, `e`)
- **Config 저장 기능** - `config.Save()`, `config.DeepCopy()` 함수 추가
- **위젯 메타데이터** - 모든 위젯의 옵션 정의 (`internal/tui/widget_options.go`)

### Dependencies

- `github.com/charmbracelet/bubbletea v1.2.4` - TUI 프레임워크
- `github.com/charmbracelet/bubbles v0.20.0` - TUI 컴포넌트
- `github.com/charmbracelet/lipgloss v1.0.0` - 스타일링

## [0.4.0] - 2026-02-03

### Added

- **위젯별 Threshold 커스터마이징** - Extra 옵션으로 색상 임계값 설정 가능
  - `context`: `warn_threshold` (60%), `critical_threshold` (80%)
  - `cost`: `warn_threshold` ($0.50), `critical_threshold` ($1.00)
  - `cache_hit`: `good_threshold` (80%), `warn_threshold` (50%)
  - `api_latency`: `warn_threshold` (2000ms), `critical_threshold` (5000ms)
  - `burn_rate`: `warn_threshold` (10¢/min), `critical_threshold` (25¢/min)
- **`block_timer` 위젯** - 5시간 Claude Pro 사용량 블록 남은 시간 표시
  - `Block: 4h23m` 형식으로 출력
  - 80% 경과 시 노란색, 95% 경과 시 빨간색
  - 블록 만료 시 자동 갱신
- **`GetExtraFloat()` 헬퍼** - Extra 맵에서 float64 값 파싱
- **GitHub Actions 자동 릴리즈**
  - `.goreleaser.yml` - 멀티 플랫폼 빌드 설정
  - `.github/workflows/release.yml` - 태그 푸시 시 자동 릴리즈
  - Linux/macOS (amd64/arm64) 바이너리 자동 생성

### Changed

- 기본 위젯 10개 → 11개 (`block_timer` 추가)
- `version` 변수가 ldflags로 주입 가능하게 변경 (`-X main.version=X.Y.Z`)
- History 구조체에 `BlockStartTime` 필드 추가

## [0.3.1] - 2026-02-03

### Fixed

- **tailLines 성능 최적화** - EOF에서 역방향으로 읽어 대용량 파일 처리 개선 (#16)
  - 전체 파일을 읽지 않고 필요한 만큼만 청크 단위로 읽음
  - 4KB * N 청크로 시작, 최대 1MB까지 증가
  - 대용량 트랜스크립트에서 메모리 사용량 감소

## [0.3.0] - 2026-02-03

### Added

- **Transcript 파싱** - Claude Code JSONL 트랜스크립트 파일에서 tool/agent 데이터 추출
  - `internal/transcript/` 패키지 추가
  - 마지막 100줄 파싱 (메모리 효율적)
  - 잘못된 JSON 라인 graceful skip
- **새 위젯 2종**
  - `tools` - 최근 도구 호출 상태 표시 (`✓Read ✓Write ◐Bash`)
  - `agents` - 서브 에이전트 상태 표시 (`◐ 1 agent`, `✓ 2 done`)
- Session struct에 `transcript_path` 필드 추가

### Changed

- tools 위젯이 약어 대신 풀 네임 표시 (R → Read, W → Write)

## [0.2.0] - 2026-02-03

### Added

- **새 위젯 3종**
  - `burn_rate` - 비용 번 레이트 (¢/min 또는 $/min)
  - `compact_eta` - 80% context 도달 예측 시간
  - `context_spark` - 히스토리 기반 스파크라인 (`▂▃▄▅▆`)
- **Split 레이아웃** - 좌/우 정렬 지원 (`[[line.left]]`, `[[line.right]]`)
- **세션 히스토리 버퍼** - `~/.cache/visor/history_<session>.json`에 세션별 히스토리 저장
- **조건부 위젯 렌더링** - `show_when_above` 옵션으로 threshold 기반 표시/숨김
- `[general]` 섹션의 `separator` 설정 - 위젯 간 구분자 커스터마이징 (기본값: `" | "`)
- Context 위젯 프로그레스 바 - `Ctx: 42% ████░░░░░░` 형식
- Session struct에 `total_duration_ms`, `session_id` 필드 추가

### Changed

- 기본 위젯 7개 → 10개 (burn_rate, compact_eta, context_spark 추가)
- 테스트 커버리지 대폭 개선
  - `internal/git`: 0% → 80.9%
  - `internal/render`: 74.7% → 90.8%
  - `internal/widgets`: 58.9% → 83.6%

### Security

- Session ID sanitization 추가 - path traversal 방지 (#14)
  - 영문, 숫자, `-`, `_`만 허용
  - 최대 64자 제한

## [0.1.1] - 2026-02-02

### Fixed

- Git 명령어에 200ms 타임아웃 추가 - 대형 저장소에서 statusline 멈춤 방지 (#1)
- `parseInt()` 함수 버그 수정 - `strconv.Atoi()` 사용 (#2)
- cost 위젯 중복 코드 제거 (#3)

### Added

- `--debug` 플래그 - config 에러 등 디버깅 정보 출력 (#4)
- `format` 필드 - 위젯 출력 포맷 커스터마이징 (예: `format = "Context: {value}"`) (#7)
- `extra` 필드 - 위젯별 추가 옵션 (예: `show_label = "false"`) (#7)
- 테스트 커버리지 대폭 개선 (#6)
  - `internal/config`: 0% → 82.8%
  - `internal/render`: ~50% → 74.7%
  - `internal/widgets`: ~30% → 47.5%

### Changed

- 임계값을 상수로 추출하여 코드 가독성 향상 (#5)
  - `ContextWarningPct`, `ContextDangerPct`
  - `CostWarningUSD`, `CostDangerUSD`
  - `CacheHitGoodPct`, `CacheHitWarningPct`
  - `LatencyWarningMs`, `LatencyDangerMs`
- `ColorByThreshold()`, `ColorByThresholdInverse()` 헬퍼 함수 추가

## [0.1.0] - 2026-02-02

### Added

- 초기 릴리스
- **Core**
  - stdin JSON 파싱 (`internal/input`)
  - TOML 설정 시스템 (`internal/config`)
  - ANSI 컬러 렌더링 (`internal/render`)
  - Widget 인터페이스 및 Registry (`internal/widgets`)

- **Widgets**
  - `model` - 현재 모델명 표시
  - `context` - 컨텍스트 윈도우 사용률 (색상 코딩)
  - `cost` - 세션 총 비용
  - `git` - 브랜치, staged/modified, ahead/behind 상태
  - `cache_hit` - 캐시 히트율 (고유 메트릭)
  - `api_latency` - API 총 지연시간 (고유 메트릭)
  - `code_changes` - 추가/삭제 라인 수 (고유 메트릭)

- **CLI**
  - `--version` - 버전 정보 출력
  - `--init` - 기본 설정 파일 생성
  - `--setup` - Claude Code 연동 가이드
  - `--check` - 설정 파일 유효성 검사

- **설정**
  - `~/.config/visor/config.toml` 지원
  - 멀티라인 레이아웃 (`[[line]]`)
  - 위젯 순서 커스터마이징
  - 위젯별 스타일 설정 (fg, bg, bold)

### Performance

- Cold startup < 20ms
- 잘못된 JSON에서 panic 없이 graceful fallback

### Dependencies

- `github.com/BurntSushi/toml v1.3.2` - TOML 파싱

---

## 버전 가이드

- **MAJOR** (X.0.0): 하위 호환되지 않는 변경
- **MINOR** (0.X.0): 하위 호환되는 기능 추가
- **PATCH** (0.0.X): 하위 호환되는 버그 수정

## 링크

[Unreleased]: https://github.com/namyoungkim/visor/compare/v0.11.0...HEAD
[0.11.0]: https://github.com/namyoungkim/visor/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/namyoungkim/visor/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/namyoungkim/visor/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/namyoungkim/visor/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/namyoungkim/visor/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/namyoungkim/visor/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/namyoungkim/visor/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/namyoungkim/visor/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/namyoungkim/visor/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/namyoungkim/visor/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/namyoungkim/visor/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/namyoungkim/visor/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/namyoungkim/visor/releases/tag/v0.1.0
