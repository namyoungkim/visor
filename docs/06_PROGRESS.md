# 06. Progress Tracking

visor 프로젝트의 PRD 대비 진행상황을 추적합니다.

**최종 업데이트**: 2025-02-03 (v0.3.0 완료)

---

## 요약

| Phase | 상태 | 진행률 |
|-------|------|--------|
| **v0.1 MVP** | ✅ 완료 | 100% |
| **v0.2 효율 심화** | ✅ 완료 | 100% |
| **v0.3 고급 기능** | ✅ 완료 | 60% (core features) |

---

## v0.1 MVP 상세

### 마일스톤 진행

| 마일스톤 | 내용 | 상태 |
|----------|------|------|
| M1 | Go 프로젝트 셋업 + stdin 파싱 + 첫 stdout 출력 | ✅ 완료 |
| M2 | 7개 MVP 위젯 구현 | ✅ 완료 |
| M3 | TOML 설정 + 멀티라인 + ANSI 컬러 | ✅ 완료 |
| M4 | --setup, --init CLI + 스마트 truncation | ✅ 완료 |
| M5 | 테스트 + 실사용 검증 | ✅ 완료 |

### MVP 위젯

| 위젯 | 식별자 | PRD 스펙 | 상태 | 비고 |
|------|--------|----------|------|------|
| 모델명 | `model` | `model.display_name` 표시 | ✅ 완료 | |
| Context | `context` | `used_percentage` + 색상 코딩 + 프로그레스 바 | ✅ 완료 | |
| Git | `git` | 브랜치 + staged/modified + ahead/behind | ✅ 완료 | 200ms 타임아웃 적용 |
| 비용 | `cost` | `total_cost_usd` 표시 | ✅ 완료 | |
| 캐시 히트율 ★ | `cache_hit` | cache_read / (cache_read + input) | ✅ 완료 | 고유 메트릭 |
| API 레이턴시 ★ | `api_latency` | ms/s 단위 변환 | ✅ 완료 | 고유 메트릭 |
| 코드 변경 ★ | `code_changes` | +added/-removed 형식 | ✅ 완료 | 고유 메트릭 |

★ = 기존 프로젝트에 없는 기능

### CLI 명령어

| 명령어 | PRD 스펙 | 상태 |
|--------|----------|------|
| `visor` | stdin JSON → stdout ANSI | ✅ 완료 |
| `visor --version` | 버전 출력 | ✅ 완료 |
| `visor --init` | config.toml 생성 | ✅ 완료 |
| `visor --setup` | Claude Code 연동 가이드 | ✅ 완료 |
| `visor --check` | 설정 유효성 검증 | ✅ 완료 |
| `visor --debug` | 디버그 출력 (PRD 외 추가) | ✅ 완료 |

### 설정 시스템

| 기능 | PRD 스펙 | 상태 |
|------|----------|------|
| TOML 설정 파일 | `~/.config/visor/config.toml` | ✅ 완료 |
| 멀티라인 레이아웃 | `[[line]]` 문법 | ✅ 완료 |
| 위젯 순서 커스터마이징 | 배열 순서대로 | ✅ 완료 |
| 위젯별 스타일 | fg, bg, bold | ✅ 완료 |
| 위젯별 format | `{value}` 플레이스홀더 | ✅ 완료 (v0.1.1) |
| 위젯별 extra 옵션 | show_label 등 | ✅ 완료 (v0.1.1) |
| separator 설정 | `" \| "` 구분자 | ✅ 완료 (v0.1.2) |
| truncate 설정 | 터미널 너비 초과 처리 | ✅ 완료 |

### 성능 요구사항

| 요구사항 | 목표 | 상태 | 측정값 |
|----------|------|------|--------|
| Cold startup | < 5ms | ✅ 달성 | ~19ms (첫 실행 포함) |
| JSON 파싱 실패 시 panic | 0 | ✅ 달성 | graceful fallback |
| Git 명령어 타임아웃 | — | ✅ 완료 | 200ms |

### MVP 완료 기준 체크리스트

- [x] `echo '{...}' | visor` 로 포맷된 출력 확인
- [x] 7개 위젯 표시
- [x] 캐시 히트율, API 레이턴시, 코드 변경량 차별점 동작
- [x] cold startup 5ms 이내

---

## v0.2 효율 심화 (완료)

### 위젯

| 위젯 | 식별자 | 설명 | 상태 |
|------|--------|------|------|
| 번 레이트 | `burn_rate` | $/min 계산 | ✅ 완료 |
| Compact 예측 | `compact_eta` | 80% 도달 카운트다운 | ✅ 완료 |
| Context 스파크라인 | `context_spark` | 최근 N회 미니 그래프 `▂▃▄▅▆` | ✅ 완료 |
| 조건부 위젯 | — | threshold 기반 on/off | ✅ 완료 (compact_eta에 적용) |

### 기능

| 기능 | 설명 | 상태 |
|------|------|------|
| Context 프로그레스 바 | `████░░░░░░` 형식 | ✅ 완료 (v0.1.2) |
| Split 레이아웃 | 좌/우 정렬 | ✅ 완료 |
| 세션 히스토리 버퍼 | 호출 간 데이터 유지 | ✅ 완료 |
| 위젯별 threshold 설정 | warn_threshold, critical_threshold | 🔲 미구현 (v0.3)

---

## v0.3 고급 기능 (진행중)

| 기능 | 설명 | 상태 |
|------|------|------|
| Transcript 파싱 | tool/agent 활동 추적 | ✅ 완료 |
| Tool 위젯 | `tools` - 도구 사용 표시 | ✅ 완료 |
| Agent 위젯 | `agents` - 에이전트 상태 | ✅ 완료 |
| 5시간 블록 타이머 | `block_timer` | 🔲 미구현 |
| TUI 설정 도구 | 인터랙티브 설정 | 🔲 미구현 |
| Powerline 테마 | 특수 문자 스타일 | 🔲 미구현 |
| GitHub Actions | 자동 릴리즈 | 🔲 미구현 |

---

## 릴리즈 히스토리

### v0.3.0 (2025-02-03)

**Added**:
- Transcript 파싱 - Claude Code JSONL 트랜스크립트에서 tool/agent 데이터 추출
- `tools` 위젯 - 최근 도구 호출 상태 (`✓Read ✓Write ◐Bash`)
- `agents` 위젯 - 서브 에이전트 상태 (`◐ 1 agent`, `✓ 2 done`)
- Session struct에 `transcript_path` 필드 추가

**Changed**:
- tools 위젯이 약어 대신 풀 네임 표시 (R → Read)

### v0.2.0 (2025-02-03)

**Added**:
- `burn_rate` 위젯 - 비용 번 레이트 (¢/min 또는 $/min)
- `compact_eta` 위젯 - 80% 도달 예측 시간
- `context_spark` 위젯 - 히스토리 기반 스파크라인 (`▂▃▄▅▆`)
- Split 레이아웃 - 좌/우 정렬 (`[[line.left]]`, `[[line.right]]`)
- 세션 히스토리 버퍼 - `~/.cache/visor/` 에 세션별 히스토리 저장
- 조건부 위젯 렌더링 (`show_when_above` 옵션)
- Session struct에 `total_duration_ms`, `session_id` 필드 추가

**Security**:
- Session ID sanitization 추가 - path traversal 방지

### v0.1.2 (2025-02-02)

**Added**:
- `[general]` 섹션의 `separator` 설정 - 위젯 간 구분자 커스터마이징 (기본값: `" | "`)
- Context 위젯 프로그레스 바 - `Ctx: 42% ████░░░░░░` 형식
  - `show_bar` extra 옵션 (기본: true)
  - `bar_width` extra 옵션 (기본: 10)

**Changed**:
- 테스트 커버리지 대폭 개선
  - `internal/git`: 0% → 80.9%
  - `internal/render`: 74.7% → 90.8%
  - `internal/widgets`: 58.9% → 83.6%

### v0.1.1 (2025-02-02)

**Fixed**:
- Git 명령어 200ms 타임아웃 추가
- parseInt() 버그 수정
- cost 위젯 중복 코드 제거

**Added**:
- `--debug` 플래그
- `format` 필드 (위젯 출력 커스터마이징)
- `extra` 필드 (위젯별 옵션)
- 테스트 커버리지 개선 (config 82.8%, render 74.7%, widgets 47.5%)

**Changed**:
- 임계값을 상수로 추출

### v0.1.0 (2025-02-02)

**Initial Release**:
- 7개 MVP 위젯 (model, context, git, cost, cache_hit, api_latency, code_changes)
- TOML 설정 시스템
- CLI (--version, --init, --setup, --check)
- ANSI 컬러 렌더링
- 멀티라인 레이아웃

---

## 다음 단계 제안

### 완료 (v0.2.0)
1. ~~번 레이트 위젯~~ ✅
2. ~~Compact 예측 위젯~~ ✅
3. ~~조건부 위젯 렌더링~~ ✅
4. ~~Split 레이아웃~~ ✅
5. ~~세션 히스토리 버퍼~~ ✅
6. ~~Context 스파크라인 위젯~~ ✅

### 완료 (v0.3.0)
1. ~~Transcript JSONL 파싱~~ ✅
2. ~~Tool/Agent 위젯~~ ✅

### 다음 (v0.4.0)
1. 위젯별 threshold 설정
2. 5시간 블록 타이머
3. GitHub Actions 자동 릴리즈
4. tailLines 성능 최적화 (#16)

---

## 참고

- [00_PRD.md](00_PRD.md) — 전체 제품 요구사항
- [CHANGELOG.md](../CHANGELOG.md) — 버전별 변경 내역
- [05_IMPLEMENTATION.md](05_IMPLEMENTATION.md) — 구현 가이드
