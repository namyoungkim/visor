# User Story Mapping

## 1. Backbone

사용자가 visor를 사용하는 여정을 시간 순서로 정의한다.

```
설치(Install) → 설정(Configure) → 세션 시작(Session) → 모니터링(Monitor) → 최적화(Optimize)
```

---

## 2. Story Map

### MVP (v0.1)

| Install | Configure | Session | Monitor | Optimize |
|---------|-----------|---------|---------|----------|
| IN-1: go install로 설치 | CF-1: 기본 설정으로 즉시 작동 | SE-1: Claude Code 시작 시 자동 로드 | MO-1: 모델명 확인 | OP-1: 캐시 히트율로 프롬프트 구조 판단 |
| IN-2: settings.json 자동 연결 | CF-2: TOML로 위젯 선택/순서 조정 | SE-2: stdin JSON 수신 | MO-2: Context % + 바 확인 | OP-2: API 레이턴시로 병목 파악 |
| | CF-3: 멀티라인 구성 | | MO-3: Git 상태 확인 | OP-3: 코드 변경량 대비 비용 체크 |
| | | | MO-4: 세션 비용 확인 | |

### Phase 2 (v0.2)

| Install | Configure | Session | Monitor | Optimize |
|---------|-----------|---------|---------|----------|
| | CF-4: 조건부 위젯 threshold 설정 | SE-3: 세션 내 히스토리 버퍼 유지 | MO-5: 번 레이트 ($/min) | OP-4: Context 추이 스파크라인으로 패턴 감지 |
| | CF-5: 좌/우 split 레이아웃 설정 | | MO-6: Compact 예측 카운트다운 | OP-5: Compact 임박 시 사전 대응 |
| | | | MO-7: 조건부 경고 표시 | |

### Phase 3 (v0.3+)

| Install | Configure | Session | Monitor | Optimize |
|---------|-----------|---------|---------|----------|
| | CF-6: TUI 설정 도구 | SE-4: Transcript JSONL 파싱 | MO-8: Tool/Agent 활동 추적 | OP-6: 프로젝트별 누적 비용 비교 |
| | CF-7: 테마/Powerline 설정 | SE-5: 세션 히스토리 영속 저장 | MO-9: 5시간 블록 타이머 | |

---

## 3. Story Details (MVP)

### IN-1: go install로 설치

**설명**: `go install` 한 줄로 바이너리 설치. 외부 런타임(Node, Python, Ruby) 불필요.

**Acceptance Criteria**:
- [ ] `go install github.com/leo/visor@latest` 로 설치 가능
- [ ] 설치 후 `visor --version`으로 버전 확인
- [ ] Linux amd64, macOS arm64 빌드 동작

### IN-2: settings.json 자동 연결

**설명**: `visor --setup` 실행 시 `~/.claude/settings.json`에 statusLine 설정 자동 추가.

**Acceptance Criteria**:
- [ ] 기존 settings.json이 있으면 statusLine 키만 추가/갱신
- [ ] 없으면 새로 생성
- [ ] 설정 후 Claude Code 재시작하면 statusline 동작
- [ ] 결과: `{"statusLine": {"type": "command", "command": "visor", "padding": 0}}`

### CF-1: 기본 설정으로 즉시 작동

**설명**: 설정 파일 없이도 합리적인 기본값으로 동작.

**Acceptance Criteria**:
- [ ] `~/.config/visor/config.toml` 없으면 기본 레이아웃 사용
- [ ] 기본: `[Model] | [Git Branch] | Context: XX% ████░░ | $X.XX | Cache: XX% | API: X.Xs | +XX/-XX`
- [ ] 첫 실행 시 설정 파일 없어도 에러 없이 정상 출력

### CF-2: TOML로 위젯 선택/순서 조정

**설명**: TOML 설정 파일로 어떤 위젯을 어떤 순서로 표시할지 커스터마이징.

**Acceptance Criteria**:
- [ ] `visor --init`으로 기본 설정 파일 생성
- [ ] 위젯 on/off, 순서 변경 가능
- [ ] 유효하지 않은 설정 시 에러 메시지 + 기본값 fallback
- [ ] 설정 예시:
```toml
[[line]]
widgets = ["model", "git", "context", "cost", "cache_hit", "api_latency", "code_changes"]
separator = " | "
```

### CF-3: 멀티라인 구성

**설명**: 여러 줄에 걸쳐 위젯을 배치.

**Acceptance Criteria**:
- [ ] `[[line]]` 블록을 여러 개 정의하면 멀티라인 출력
- [ ] 각 라인별 독립적인 위젯 구성
- [ ] 기본값: 2줄 (Line 1: 상태, Line 2: 효율)

### SE-1: Claude Code 시작 시 자동 로드

**설명**: Claude Code의 statusline hook으로 자동 실행.

**Acceptance Criteria**:
- [ ] Claude Code 시작 시 visor가 자동 호출됨
- [ ] ~300ms 주기로 반복 호출되어도 안정적
- [ ] 프로세스 즉시 종료 (데몬 아님)

### SE-2: stdin JSON 수신

**설명**: Claude Code가 전달하는 JSON을 stdin으로 읽어 파싱.

**Acceptance Criteria**:
- [ ] Session 구조체로 완전 파싱
- [ ] 누락 필드 시 zero value fallback (panic 없음)
- [ ] 파싱 실패 시 빈 출력 (에러 메시지를 stdout에 쓰지 않음)
- [ ] mock JSON으로 수동 테스트 가능: `echo '{"model":...}' | visor`

### MO-1: 모델명 확인

**설명**: 현재 사용 중인 Claude 모델명 표시.

**Acceptance Criteria**:
- [ ] `model.display_name` 값 그대로 출력
- [ ] ANSI 볼드/컬러 적용 가능

### MO-2: Context % + 프로그레스 바 확인

**설명**: Context window 사용률을 숫자 + 시각적 바로 표시.

**Acceptance Criteria**:
- [ ] `context_window.used_percentage` 기반
- [ ] 프로그레스 바: `████░░░░░░` 스타일 (10칸)
- [ ] 0-60%: 초록, 60-80%: 노랑, 80-100%: 빨강

### MO-3: Git 상태 확인

**설명**: 현재 Git 브랜치와 변경사항 표시.

**Acceptance Criteria**:
- [ ] 브랜치명 표시 (예: `main`, `feature/hud`)
- [ ] staged/unstaged/untracked 카운트 (예: `+3 ~2 ?1`)
- [ ] Git 저장소가 아니면 표시하지 않음 (위젯 숨김)
- [ ] `workspace.current_dir` 기준으로 git 명령 실행

### MO-4: 세션 비용 확인

**설명**: 현재 세션의 누적 비용 표시.

**Acceptance Criteria**:
- [ ] `cost.total_cost_usd` 기반
- [ ] $0.01 이하: 3자리, $1 이상: 2자리 (예: `$0.003`, `$1.24`)

### OP-1: 캐시 히트율 ★

**설명**: prompt caching 효율을 실시간 %로 표시.

**Acceptance Criteria**:
- [ ] 계산: `cache_read_input_tokens / (cache_read_input_tokens + input_tokens) * 100`
- [ ] `current_usage`가 null이면 "—" 표시
- [ ] 50% 이상: 초록, 30-50%: 노랑, 30% 미만: 빨강

### OP-2: API 레이턴시 ★

**설명**: Claude API 응답 시간을 표시.

**Acceptance Criteria**:
- [ ] `cost.total_api_duration_ms / (호출 횟수 추정)` 또는 그냥 total 표시
- [ ] 포맷: `1.2s`, `850ms`
- [ ] 2s 이상: 빨강, 1-2s: 노랑, 1s 미만: 초록

### OP-3: 코드 변경량 ★

**설명**: 세션 중 추가/삭제된 코드 라인 수 표시.

**Acceptance Criteria**:
- [ ] `cost.total_lines_added` / `cost.total_lines_removed` 기반
- [ ] 포맷: `+156/-23`
- [ ] 초록(추가) / 빨강(삭제) 컬러

---

## 4. MVP Validation Scenario

### 시나리오: "일상적인 Claude Code 코딩 세션"

```
프로젝트: visor 자체 개발
태스크: Go 코드 작성 + 디버깅 + 리팩토링

성공 기준:
- statusline 글랜스만으로 5가지 핵심 정보 즉시 파악 가능
- context 60% 넘어갈 때 자연스럽게 인지
- 캐시 히트율이 낮으면 프롬프트 구조를 조정하려는 행동 변화 발생
- cold startup이 체감되지 않음

검증 흐름:
1. [Install] go install로 visor 설치 → settings.json 연결
2. [Configure] 기본 설정으로 시작 (설정 파일 없이)
3. [Session] Claude Code 실행 → 하단에 statusline 자동 표시
4. [Monitor] 코딩하면서 context %, 비용, git 상태 글랜스 확인
5. [Optimize] 캐시 히트율 40% 확인 → 대화 구조 조정 → 65%로 개선 확인
```
