# Changelog

모든 주요 변경사항을 이 파일에 기록합니다.

형식은 [Keep a Changelog](https://keepachangelog.com/ko/1.0.0/)를 따르며,
버전은 [Semantic Versioning](https://semver.org/lang/ko/)을 따릅니다.

## [Unreleased]

### Added

- `[general]` 섹션의 `separator` 설정 - 위젯 간 구분자 커스터마이징 (기본값: `" | "`)
- Context 위젯 프로그레스 바 - `Ctx: 42% ████░░░░░░` 형식 (기본 활성화)
  - `show_bar`: "true"/"false" - 프로그레스 바 표시 여부 (기본: true)
  - `bar_width`: "10" - 프로그레스 바 너비 (기본: 10)

### Changed

- 테스트 커버리지 대폭 개선
  - `internal/git`: 0% → 80.9%
  - `internal/render`: 74.7% → 90.8%
  - `internal/widgets`: 58.9% → 83.6%

## [0.1.1] - 2025-02-02

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

## [0.1.0] - 2025-02-02

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

[Unreleased]: https://github.com/namyoungkim/visor/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/namyoungkim/visor/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/namyoungkim/visor/releases/tag/v0.1.0
