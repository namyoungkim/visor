# Architecture Decision Records

## Decision Log

| ADR | Decision | Status |
|-----|----------|--------|
| 001 | 구현 언어: Go | Accepted |
| 002 | 설정 방식: TOML | Accepted |
| 003 | 아키텍처: Widget 인터페이스 패턴 | Accepted |
| 004 | Git 정보: 외부 CLI 호출 | Accepted |
| 005 | 차별화 전략: 효율성 메트릭스 | Accepted |
| 006 | 배포: go install + 단일 바이너리 | Accepted |
| 007 | 프로젝트 구조: internal 패키지 | Accepted |

---

## ADR-001: 구현 언어 — Go

### Status
Accepted

### Context
Claude Code statusline은 ~300ms마다 새 프로세스로 spawn되는 도구다.
기존 생태계 분석 결과:

- TypeScript/Node: 가장 큰 생태계 (ccstatusline ⭐1.9k), 하지만 30-50ms startup
- Rust: 성능 최적 (1-2ms startup), 하지만 4개+ 프로젝트가 이미 존재
- Bash: 0ms startup, 하지만 JSON 파싱/로직 복잡도 한계
- Go: 1-2ms startup, 제대로 된 statusline 구현체가 없음

요구사항:
- Cold startup 5ms 이내 (300ms 예산의 2% 미만)
- JSON 파싱 + 문자열 가공 + 외부 명령(git) 실행
- 단일 바이너리 배포 (런타임 의존성 없음)
- 빠른 개발 이터레이션 (수정 → 확인 반복)

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **Go** | 1-2ms startup, 즉시 컴파일, stdlib JSON, 단일 바이너리, 빈 포지션 | 새로 배워야 함 |
| **Rust** | 1-2ms startup, 최고 성능, 타입 안전성 | 이미 4+ 프로젝트 존재, 컴파일 대기 시간, 이 규모에서 over-engineering |
| **TypeScript** | 가장 큰 생태계, 빠른 개발 | 30-50ms startup, Node/Bun 런타임 의존성 |
| **Bash** | 0ms startup, 의존성 없음 | jq 필요, 복잡한 로직 유지보수 어려움 |

### Decision
**Go 선택**

### Rationale
1. **빈 포지션 선점**: Go로 만든 전용 statusline이 사실상 없음
2. **적절한 복잡도 매칭**: 이 프로젝트는 "JSON 읽기 → 문자열 가공 → 출력" 파이프라인으로, Rust의 소유권/라이프타임이 가치를 발휘할 장면이 없음
3. **개발 속도**: `go build`가 수백ms — ANSI 출력 미세 조정의 빠른 이터레이션에 필수적
4. **Startup 동등**: Go와 Rust 모두 ~1-2ms로 실질 차이 없음
5. **단일 바이너리**: `go build` → 즉시 배포 가능, 크로스 컴파일 trivial
6. **첫 Go 프로젝트로 적정 규모**: 고루틴/채널 불필요, struct + interface만으로 충분

### Consequences

**긍정적**:
- Go 생태계에서의 차별화 (유일한 전용 statusline)
- 포트폴리오에 Go CLI 경험 추가
- 런타임 의존성 제로 (사용자 설치 마찰 최소)

**부정적**:
- Go 학습 곡선 (첫 프로젝트)
- Rust 학습과 시간 분산 가능성

**중립적**:
- ccstatusline(TS)이나 CCometixLine(Rust)과 직접 경쟁하지 않음 — 다른 언어, 다른 포지셔닝

---

## ADR-002: 설정 방식 — TOML

### Status
Accepted

### Context
사용자 설정을 어떤 포맷으로 관리할지 결정 필요.
기존 프로젝트들:
- ccstatusline: React/Ink TUI로 JSON 생성
- CCometixLine: TOML + TUI
- rz1989s: TOML + 환경변수
- claude_monitor: 환경변수만

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **TOML** | 사람 읽기 좋음, Go에 좋은 라이브러리, CLI 도구 표준 | JSON보다 파서 추가 필요 |
| **JSON** | Go stdlib 지원, Claude Code와 동일 포맷 | 주석 불가, 사람이 편집하기 불편 |
| **YAML** | 가독성 좋음 | 들여쓰기 실수 위험, Go 외부 라이브러리 필요 |
| **환경변수** | 가장 단순 | 복잡한 설정(멀티라인, 위젯 옵션) 표현 한계 |

### Decision
**TOML 선택**

### Rationale
1. CLI 도구의 사실상 표준 (Cargo.toml, pyproject.toml 등)
2. `[[line]]` 배열 문법이 멀티라인 위젯 구성에 자연스러움
3. 주석 가능 → 설정 파일 자체가 문서 역할
4. Go용 `BurntSushi/toml` 라이브러리 성숙

### Consequences

**긍정적**:
- 설정 파일이 자기문서화 (주석으로 옵션 설명)
- `[[line]]` 배열이 멀티라인 구성과 1:1 매핑

**부정적**:
- 외부 의존성 1개 추가 (`BurntSushi/toml`)
- JSON보다 약간 덜 보편적

---

## ADR-003: 아키텍처 — Widget 인터페이스 패턴

### Status
Accepted

### Context
statusline에 표시할 정보(모델, context, git 등)를 어떤 구조로 관리할지 결정 필요.

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **Widget interface** | 각 위젯 독립적 추가/제거, 테스트 용이, 설정으로 on/off | interface 오버헤드 |
| **단일 함수** | 가장 단순, 빠르게 구현 | 기능 추가 시 함수가 거대해짐 |
| **Template 기반** | 포맷 문자열로 유연 | 로직 표현 한계, 조건부 렌더링 어려움 |

### Decision
**Widget interface 선택**

### Rationale
1. 새 위젯 추가가 파일 하나 + Registry 등록으로 완결
2. `ShouldRender()` 메서드로 조건부 표시 자연스럽게 지원
3. 각 위젯 단위 테스트 가능
4. 설정 파일의 위젯 이름과 Registry가 1:1 매핑

### Consequences

**긍정적**:
- 확장성: Phase 2/3 위젯 추가 시 기존 코드 수정 불필요
- 테스트: 위젯별 독립 테스트 + mock Session으로 검증

**부정적**:
- MVP 7개 위젯에 interface는 약간 과잉일 수 있으나, Phase 2 이후를 고려하면 합리적

---

## ADR-004: Git 정보 — 외부 CLI 호출

### Status
Accepted

### Context
Git 브랜치/상태 정보를 얻는 방법.

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **외부: `git` CLI 호출** | 구현 간단, 100% 호환, 의존성 0 | exec 오버헤드 (~5-10ms) |
| **라이브러리: go-git** | exec 없음, 순수 Go | 큰 의존성, 일부 기능 미지원 |
| **.git 디렉토리 직접 파싱** | 의존성 0, 빠름 | 구현 복잡, edge case 많음 |

### Decision
**외부 `git` CLI 호출 선택**

### Rationale
1. Claude Code 사용자는 반드시 git이 설치되어 있음
2. 5-10ms 오버헤드는 300ms 예산 내 충분
3. `git rev-parse --abbrev-ref HEAD`, `git status --porcelain` 2개 명령이면 충분
4. go-git 의존성(수 MB)은 단일 바이너리 사이즈에 불필요한 부담

### Consequences

**긍정적**:
- 의존성 최소화
- git의 모든 기능(worktree, sparse checkout 등) 자동 호환

**부정적**:
- git 미설치 환경에서 위젯 비활성화 필요 (ShouldRender에서 처리)
- subprocess 생성 비용 (단, 캐싱으로 완화 가능)

---

## ADR-005: 차별화 전략 — 효율성 메트릭스

### Status
Accepted

### Context
이미 10+ 프로젝트가 존재하는 Claude Code statusline 생태계에서 어떻게 차별화할지.

기존 프로젝트 분석 결과, 전부 "현재 상태 표시"에 집중:
- 모델명, 브랜치, context %, 토큰 수, 비용 절대값

stdin JSON에 있지만 아무도 활용하지 않는 데이터:
- `current_usage.cache_read_input_tokens` → 캐시 히트율
- `cost.total_api_duration_ms` → API 레이턴시
- `cost.total_lines_added/removed` → 생산성 지표

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **효율성 메트릭스 특화** | 명확한 차별점, 실용적 가치, 이미 있는 데이터 활용 | 일반 사용자에겐 과할 수 있음 |
| **미용/테마 특화** | 시각적 임팩트 | ccstatusline이 이미 지배적 |
| **기능 전부 구현** | 완성도 | 차별점 없음, 개발 기간 길어짐 |

### Decision
**효율성 메트릭스 특화**

### Rationale
1. "지금 얼마나 효율적인가"는 아무도 답하지 않는 질문
2. 캐시 히트율은 프롬프트 최적화의 핵심 피드백 루프
3. Compact 예측은 갑작스러운 context 초기화를 사전에 방지
4. 모든 데이터가 이미 stdin JSON에 있어 구현 복잡도 낮음

### Consequences

**긍정적**:
- 명확한 포지셔닝: "효율성 대시보드"
- 기존 프로젝트와 보완적 관계 (경쟁이 아닌 대안)

**부정적**:
- Powerline 테마, TUI 설정 등 "보이는" 기능 부재로 첫인상이 소박할 수 있음

---

## ADR-006: 배포 — go install + 단일 바이너리

### Status
Accepted

### Context
사용자가 visor를 어떻게 설치할지.

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **go install** | Go 사용자에게 가장 자연스러움, 소스에서 빌드 | Go 설치 필요 |
| **npm publish** | 가장 넓은 도달 (Node 생태계) | Go 바이너리를 npm으로 감싸는 복잡성 |
| **GitHub Releases (바이너리)** | 런타임 불필요 | 수동 다운로드, 플랫폼별 빌드 |
| **Homebrew** | macOS 표준 | 유지보수 부담 |

### Decision
**go install (주) + GitHub Releases (보조)**

### Rationale
1. `go install github.com/leo/visor@latest` 한 줄로 완결
2. GitHub Actions로 릴리즈 시 자동 바이너리 빌드 (linux/mac, amd64/arm64)
3. Go 미설치자는 GitHub Releases에서 바이너리 다운로드
4. npm 래핑은 불필요한 복잡성

### Consequences

**긍정적**:
- 설치 경험 최단 경로
- CI/CD 자동화 용이

**부정적**:
- Go 미설치 사용자에게 한 단계 추가 (GitHub Releases로 우회)

---

## ADR-007: 프로젝트 구조 — internal 패키지

### Status
Accepted

### Context
Go 프로젝트의 패키지 구조를 어떻게 잡을지.

### Options Considered

| Option | 장점 | 단점 |
|--------|------|------|
| **cmd/ + internal/** | Go 표준 레이아웃, 캡슐화 명확 | 디렉토리 깊이 |
| **플랫 구조 (모두 main)** | 가장 단순 | 파일 수 늘어나면 혼란 |
| **pkg/ 패턴** | 외부 재사용 가능 | 라이브러리 의도 없음에 과잉 |

### Decision
**cmd/visor/ + internal/ 구조**

### Rationale
1. `internal/`은 Go에서 패키지 외부 접근을 컴파일러 수준에서 차단 — 의도치 않은 API 노출 방지
2. `cmd/visor/main.go`는 진입점만 담당, 로직은 internal에
3. 위젯 추가 시 `internal/widgets/` 아래 파일 하나 추가로 완결
4. Go 커뮤니티 관례와 일치

### Consequences

**긍정적**:
- 패키지 경계 명확
- 위젯/렌더러 독립 테스트 가능
- Phase 2/3 확장 시 구조 변경 불필요

**부정적**:
- MVP 7개 파일에 디렉토리 구조가 과할 수 있으나, 성장 여지를 고려하면 합리적
