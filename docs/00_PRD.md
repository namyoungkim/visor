# visor — Product Requirements Document

> Claude Code Efficiency Dashboard: "지금 얼마나 효율적인가"를 답하는 statusline

---

## Executive Summary

**visor**는 Go 기반의 고성능 Claude Code statusline으로, 기존 프로젝트들이 집중하는
"현재 상태 표시"를 넘어 **운영 효율성 지표(캐시 히트율, API 레이턴시, 비용 번 레이트,
Compact 예측)**를 실시간으로 제공한다.

```
┌─────────────────────────────────────────────────────────────────┐
│ Opus | my-project (main +3~2) | Ctx: 42% ████░░░░░░ | ~18m     │
│ $0.48 | Cache: 73% | API: 1.2s | +156/-23                      │
└─────────────────────────────────────────────────────────────────┘
```

### Why Now

- Claude Code statusline 생태계에 **Go 구현체가 사실상 없음** (유일한 빈 포지션)
- stdin JSON에 `cache_read_input_tokens`, `total_api_duration_ms` 등
  **효율성 데이터가 이미 존재**하지만 어떤 프로젝트도 활용하지 않음
- 기존 10+ 프로젝트와 정면 경쟁이 아닌 **보완적 포지셔닝** 가능

---

## 1. 목표 (Impact Mapping 요약)

### Goal
Claude Code 세션의 효율성을 실시간으로 모니터링하여, 비용 최적화와
context 관리에 대한 즉각적 피드백 루프를 제공한다.

### 핵심 행동 변화

| As-Is | To-Be |
|-------|-------|
| 세션 끝나고 비용 확인 | 실시간 번 레이트로 즉시 인지 |
| 캐시 효율 확인 불가 | 캐시 히트율 %로 프롬프트 구조 피드백 |
| 80% 도달 후 갑작스러운 compact | 소진 속도 기반 사전 예측 (Phase 2) |
| 응답 느려도 원인 불명 | API 레이턴시 실시간 표시 |

### Anti-Goals
- 커뮤니티 성장 / GitHub stars 경쟁
- ccstatusline급 위젯 마켓플레이스
- TUI 설정 도구 (MVP에서 제외)
- Windows 지원 (MVP에서 제외)

---

## 2. 생태계 분석 요약

### 경쟁 현황

| 프로젝트 | 언어 | ⭐ | 포지셔닝 |
|---------|------|---|---------|
| ccstatusline | TypeScript | 1.9k | 위젯 생태계 + TUI |
| CCometixLine | Rust | 933 | 고성능 + 테마 + CC 패치 |
| claude-hud | TypeScript | 16 | 플러그인 + transcript 파싱 |
| claude-powerline-rust | Rust | — | 성능 벤치마크 |
| rz1989s/statusline | Bash | — | 227 설정 + 캐싱 |
| **visor (본 프로젝트)** | **Go** | — | **효율성 메트릭스** |

### 기능 갭 — 아무도 안 하는 것

| 기능 | 데이터 소스 | 가치 |
|------|-----------|------|
| **캐시 히트율** | `current_usage.cache_read_input_tokens` | 프롬프트 최적화 피드백 |
| **API 레이턴시** | `cost.total_api_duration_ms` | 병목 원인 파악 |
| **코드 변경량** | `cost.total_lines_added/removed` | 생산성 추적 |
| **번 레이트** | `cost.total_cost_usd / duration` | 비용 속도 인지 |
| **Compact 예측** | context 소진 속도 기반 계산 | 사전 대응 |
| **Context 스파크라인** | 호출 간 히스토리 | 추이 패턴 감지 |

---

## 3. 기능 명세

### 3.1 stdin JSON 스키마

visor가 소비하는 입력 데이터 (Claude Code 제공):

```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/transcript.jsonl",
  "model": { "id": "claude-opus-4-1", "display_name": "Opus" },
  "workspace": {
    "current_dir": "/current/dir",
    "project_dir": "/project/root"
  },
  "version": "1.0.80",
  "cost": {
    "total_cost_usd": 0.48,
    "total_duration_ms": 45000,
    "total_api_duration_ms": 2300,
    "total_lines_added": 156,
    "total_lines_removed": 23
  },
  "context_window": {
    "total_input_tokens": 15234,
    "total_output_tokens": 4521,
    "context_window_size": 200000,
    "used_percentage": 42.5,
    "remaining_percentage": 57.5,
    "current_usage": {
      "input_tokens": 8500,
      "output_tokens": 1200,
      "cache_creation_input_tokens": 5000,
      "cache_read_input_tokens": 2000
    }
  }
}
```

### 3.2 위젯 명세

#### MVP 위젯 (v0.1)

| 위젯 | 식별자 | 입력 | 출력 예시 | 색상 규칙 |
|-------|--------|------|----------|----------|
| 모델명 | `model` | `model.display_name` | `Opus` | 볼드 시안 |
| Context | `context` | `context_window.used_percentage` | `42% ████░░░░░░` | 0-60% 초록, 60-80% 노랑, 80%+ 빨강 |
| Git | `git` | 외부: `git` CLI | `main +3~2?1` | 브랜치 초록, 변경 노랑 |
| 비용 | `cost` | `cost.total_cost_usd` | `$0.48` | 기본 |
| 캐시 히트율 ★ | `cache_hit` | `current_usage.*` | `Cache: 73%` | 50%+ 초록, 30-50% 노랑, 30%- 빨강 |
| API 레이턴시 ★ | `api_latency` | `cost.total_api_duration_ms` | `API: 1.2s` | <1s 초록, 1-2s 노랑, 2s+ 빨강 |
| 코드 변경 ★ | `code_changes` | `cost.total_lines_*` | `+156/-23` | 추가 초록, 삭제 빨강 |

★ = 기존 프로젝트에 없는 기능

**캐시 히트율 계산**:
```
cache_hit_rate = cache_read_input_tokens / (cache_read_input_tokens + input_tokens) × 100
```
`current_usage`가 null이면 "—" 표시.

**API 레이턴시 포맷**:
```
total_api_duration_ms >= 1000 → "X.Xs" (예: "1.2s")
total_api_duration_ms < 1000  → "XXXms" (예: "850ms")
```

#### Phase 2 위젯 (v0.2)

| 위젯 | 식별자 | 로직 | 출력 예시 |
|-------|--------|------|----------|
| 번 레이트 | `burn_rate` | total_cost_usd / (duration_ms / 60000) | `8.2¢/min` |
| Context 스파크라인 | `context_spark` | 최근 N회 호출의 used_percentage | `▁▃▅▇█` |
| Compact 예측 | `compact_eta` | (80 - current%) / 분당 소진률 | `~18m to compact` |
| 조건부 위젯 | — | threshold 기반 ShouldRender | context > 70%일 때만 경고 |

#### Phase 3 위젯 (v0.3) ✅ 완료

| 위젯 | 식별자 | 데이터 소스 | 상태 |
|-------|--------|-----------|------|
| Tool 활동 | `tools` | transcript JSONL 파싱 | ✅ `✓Read ✓Write ◐Bash` |
| Agent 상태 | `agents` | transcript JSONL 파싱 | ✅ `✓Plan ◐Explore` |

#### Phase 4 위젯 (v0.4) ✅ 완료

| 위젯 | 식별자 | 데이터 소스 | 상태 |
|-------|--------|-----------|------|
| 5시간 블록 타이머 | `block_timer` | 세션 히스토리 기반 | ✅ `Block: 4h23m` |

### 3.3 레이아웃 시스템

**MVP: 순차 나열**
```
[위젯1] | [위젯2] | [위젯3] | ...
```

**Phase 2: 좌/우 Split**
```
[좌측 위젯들]                                    [우측 위젯들]
```

**스마트 Truncation**: 터미널 너비 초과 시 우선순위 낮은 위젯부터 숨기고,
남은 위젯이 너비를 초과하면 `...`으로 자름.

### 3.4 설정 파일

**경로**: `~/.config/visor/config.toml`

```toml
# visor 설정 파일
# 생성: visor --init

[general]
separator = " | "        # 위젯 간 구분자
truncate = true           # 터미널 너비 초과 시 자름

# Line 1: 상태
[[line]]
widgets = ["model", "git", "context", "compact_eta"]

# Line 2: 효율
[[line]]
widgets = ["cost", "cache_hit", "api_latency", "code_changes"]

# 위젯별 설정
[widget.context]
bar_width = 10            # 프로그레스 바 칸 수
warn_threshold = 60       # 노란색 전환 (%)
critical_threshold = 80   # 빨간색 전환 (%)

[widget.cache_hit]
warn_threshold = 30
good_threshold = 50

[widget.api_latency]
warn_threshold = 1000     # ms, 노란색 전환
critical_threshold = 2000 # ms, 빨간색 전환

[widget.compact_eta]
show_when_above = 40      # context 40% 이상일 때만 표시
```

### 3.5 CLI 인터페이스

```
visor                  # stdin JSON 읽어서 statusline 출력 (기본 동작)
visor --version        # 버전 출력
visor --init           # ~/.config/visor/config.toml 생성
visor --setup          # ~/.claude/settings.json에 statusLine 설정 추가
visor --check          # 설정 파일 유효성 검증
```

---

## 4. 아키텍처 (C4 Model 요약)

### System Context

```
Claude Code ──stdin JSON──▶ visor ──stdout ANSI──▶ Terminal
                              │
                    ┌─────────┼──────────┐
                    ▼         ▼          ▼
                 Git Repo   Config     Transcript
                            (TOML)     (JSONL, Phase 3)
```

### 내부 구조

```
cmd/visor/main.go          ← 진입점
internal/input/             ← stdin JSON → Session 구조체
internal/config/            ← TOML 설정 로드
internal/widgets/           ← Widget 인터페이스 + 7개 구현체
internal/render/            ← 레이아웃, ANSI, truncation
internal/git/               ← git CLI 래퍼
```

### Widget 인터페이스

```go
type Widget interface {
    Name() string
    Render(session *Session, cfg *WidgetConfig) string
    ShouldRender(session *Session, cfg *WidgetConfig) bool
}
```

### 데이터 흐름

```
stdin → input.Parse() → Session
                           │
config.Load() → Config     │
                  │        │
                  ▼        ▼
            widgets.Registry.RenderAll(session, config)
                           │
                           ▼
                  render.Layout() → stdout
```

---

## 5. 기술 결정 (ADR 요약)

| 결정 | 선택 | 핵심 근거 |
|------|------|----------|
| 언어 | Go | 1-2ms startup + 빈 포지션 + 개발 속도 |
| 설정 | TOML | CLI 도구 표준 + `[[line]]` 문법 |
| 아키텍처 | Widget interface | 확장성 + 조건부 렌더링 + 테스트 용이 |
| Git | 외부 CLI 호출 | 의존성 0 + 5-10ms 허용 범위 |
| 차별화 | 효율성 메트릭스 | stdin 데이터 활용 갭 + 유일한 포지셔닝 |
| 배포 | go install + GitHub Releases | 최단 설치 경로 |
| 구조 | cmd/ + internal/ | Go 표준 + 캡슐화 |

---

## 6. 릴리즈 계획

### v0.1 — MVP (Week 1-2)

**목표**: Claude Code에서 실제로 동작하는 효율성 대시보드

| 마일스톤 | 내용 | 기간 |
|---------|------|------|
| M1 | Go 프로젝트 셋업 + stdin 파싱 + 첫 stdout 출력 | Day 1-2 |
| M2 | 7개 MVP 위젯 구현 (model, context, git, cost, cache_hit, latency, code_changes) | Day 3-5 |
| M3 | TOML 설정 + 멀티라인 + ANSI 컬러 | Day 6-8 |
| M4 | --setup, --init CLI + 스마트 truncation | Day 9-10 |
| M5 | 테스트 + 실사용 검증 + mock JSON 테스트 스크립트 | Day 11-12 |

**완료 기준**: ✅ 완료
- [x] `echo '{"model":...}' | visor` 로 포맷된 출력 확인
- [x] Claude Code에서 실제 동작하여 7개 위젯 표시
- [x] 캐시 히트율, API 레이턴시, 코드 변경량이 기존 프로젝트와의 차별점으로 동작
- [x] cold startup 5ms 이내

### v0.2 — 효율 심화 ✅ 완료

| 기능 | 설명 | 상태 |
|------|------|------|
| 번 레이트 | $/min 계산 | ✅ `burn_rate` |
| Context 스파크라인 | 최근 N회 호출 미니 그래프 | ✅ `context_spark` |
| Compact 예측 | 80% 도달 카운트다운 | ✅ `compact_eta` |
| 조건부 위젯 | threshold 기반 on/off | ✅ `show_when_above` |
| Split 레이아웃 | 좌/우 정렬 | ✅ `[[line.left/right]]` |
| 세션 히스토리 | 호출 간 데이터 유지 | ✅ `~/.cache/visor/` |

### v0.3 — 고급 기능 ✅ 완료 (core)

| 기능 | 설명 | 상태 |
|------|------|------|
| Transcript 파싱 | tool/agent 활동 추적 | ✅ `internal/transcript/` |
| `tools` 위젯 | 도구 호출 상태 표시 | ✅ `✓Read ✓Write ◐Bash` |
| `agents` 위젯 | 에이전트 상태 표시 | ✅ `✓Plan ◐Explore` |

### v0.4 — 커스터마이징 & 자동화 ✅ 완료

| 기능 | 설명 | 상태 |
|------|------|------|
| 위젯별 threshold 설정 | Extra 옵션으로 임계값 커스터마이징 | ✅ `warn_threshold`, `critical_threshold` |
| 5시간 블록 타이머 | 사용량 블록 모니터링 | ✅ `block_timer` 위젯 |
| GitHub 릴리즈 자동화 | goreleaser + GitHub Actions | ✅ `.goreleaser.yml` |

---

## 7. 성공 지표

| 지표 | 목표 | 상태 |
|------|------|------|
| **일일 사용** | 매일 Claude Code 세션에서 활성화 | ✅ 달성 |
| **Cold startup** | < 5ms | ✅ ~19ms (첫 실행 포함) |
| **차별 기능** | 3개 이상 유니크 위젯 | ✅ 9개 (cache_hit, api_latency, code_changes, burn_rate, compact_eta, context_spark, tools, agents, block_timer) |
| **설치 경험** | 2분 이내 완료 | ✅ go install → --setup → 동작 |
| **안정성** | JSON 파싱 실패 시 panic 0 | ✅ graceful fallback |

---

## 부록

### A. 관련 문서

- [01_IMPACT_MAPPING.md](01_IMPACT_MAPPING.md) — 목표, 액터, 임팩트, 산출물
- [02_USER_STORY_MAPPING.md](02_USER_STORY_MAPPING.md) — 사용자 여정, 스토리 상세, MVP 검증 시나리오
- [03_C4_MODEL.md](03_C4_MODEL.md) — 시스템 컨텍스트, 컨테이너, 컴포넌트, 데이터 플로우
- [04_ADR.md](04_ADR.md) — 아키텍처 결정 기록
- [05_IMPLEMENTATION.md](05_IMPLEMENTATION.md) — 코드 구조, API, 확장 가이드
- [06_PROGRESS.md](06_PROGRESS.md) — 진행 상황 추적

### B. 참고한 기존 프로젝트

| 프로젝트 | URL |
|---------|-----|
| ccstatusline | https://github.com/sirmalloc/ccstatusline |
| CCometixLine | https://github.com/Haleclipse/CCometixLine |
| claude-hud | https://github.com/jarrodwatts/claude-hud |
| claude-powerline-rust | https://github.com/david-strejc/claude-powerline-rust |
| rz1989s/statusline | https://github.com/rz1989s/claude-code-statusline |
| cc-statusline | https://github.com/chongdashu/cc-statusline |
| CCstatus | https://github.com/MaurUppi/CCstatus |
| Claude Code statusline docs | https://code.claude.com/docs/en/statusline |
