# Impact Mapping

## Goal (Why)

### 비즈니스 목표

**Claude Code 사용 시 "지금 얼마나 효율적인가"를 실시간으로 답할 수 있는 statusline을 만든다.**

기존 statusline 프로젝트들은 "지금 뭘 쓰고 있나"(모델명, 토큰, 브랜치)에 집중한다.
그러나 Claude Code 헤비 유저에게 진짜 필요한 정보는 비용 효율, 캐시 최적화 상태, 
context 소진 예측 같은 **운영 지표(operational metrics)** 다.

### 성공 기준

| 기준 | 측정 방법 |
|------|----------|
| 본인이 매일 쓰는 도구가 된다 | Claude Code 세션 시작 시 자동 로드되어 글랜스만으로 상황 파악 |
| 기존 대비 차별화된 기능이 3개 이상 | 캐시 히트율, API 레이턴시, Compact 예측 등 다른 프로젝트에 없는 지표 |
| Go 기반 고성능 | cold startup 5ms 이내, 300ms 호출 주기에서 랙 0 |
| 포트폴리오 프로젝트로서 완성도 | GitHub README, 스크린샷, go install 배포 가능 |

### Anti-Goals (명시적으로 하지 않는 것)

- ❌ 커뮤니티 성장 / GitHub stars 경쟁 — 본인 사용이 최우선
- ❌ 범용 프레임워크 — ccstatusline 같은 위젯 마켓플레이스가 아님
- ❌ TUI 설정 도구 (MVP) — 설정은 TOML 직접 편집으로 충분
- ❌ 모든 기존 프로젝트 기능 대체 — Powerline 테마 등 미용 기능은 후순위
- ❌ Windows 지원 (MVP) — Linux/macOS 먼저

---

## Actors (Who)

### Primary Actor

| Actor | 설명 | 핵심 니즈 |
|-------|------|----------|
| **Claude Code 헤비 유저 (본인)** | 하루 수시간 Claude Code를 사용하는 개발자 | 비용 효율 모니터링, context 소진 예측, 프롬프트 최적화 피드백 |

### Secondary Actors

| Actor | 설명 | 핵심 니즈 |
|-------|------|----------|
| Claude Code 일반 사용자 | 기본적인 세션 정보를 원하는 사용자 | 모델명, git 상태, context % 등 기본 정보 |
| Go 도구 관심자 | Go로 만든 CLI 도구에 관심 있는 개발자 | 설치 간편, 단일 바이너리, 빠른 성능 |

---

## Impacts (How)

### 행동 변화

| Impact | As-Is (현재) | To-Be (목표) |
|--------|-------------|-------------|
| **비용 인식** | 세션 끝나고 대시보드에서 확인 | 실시간 번 레이트($/min)로 즉각 인지 |
| **캐시 최적화** | 캐시 효율을 체감할 방법 없음 | 캐시 히트율 %로 프롬프트 구조 피드백 |
| **Context 관리** | 80% 도달 후 갑작스러운 auto-compact | 소진 속도 기반 사전 예측 카운트다운 |
| **API 병목 파악** | 응답 느려도 원인 불명 | API 레이턴시 실시간 확인 |
| **생산성 추적** | 체감으로만 판단 | 코드 변경량 대비 비용으로 객관적 추적 |

### Pain Points

1. Context가 갑자기 auto-compact되면서 대화 맥락이 끊김 → 예측할 수 없음
2. 세션 비용이 얼마나 나가는지 사후에만 알 수 있음
3. prompt caching이 잘 되고 있는지 확인할 방법이 없음
4. 기존 statusline들이 Node/Bun 의존성 → 환경마다 설치 번거로움

---

## Deliverables (What)

### MVP (v0.1)

| Deliverable | Impact 연결 | 우선순위 |
|-------------|-----------|---------|
| Context % + 프로그레스 바 | Context 관리 | P0 |
| 모델명 표시 | 기본 정보 | P0 |
| Git 브랜치 + 변경사항 | 기본 정보 | P0 |
| 세션 비용 ($) | 비용 인식 | P0 |
| **캐시 히트율 (%)** | 캐시 최적화 | P0 ★ |
| **API 레이턴시** | API 병목 파악 | P0 ★ |
| **코드 변경량 (+/-)** | 생산성 추적 | P0 ★ |
| TOML 설정 파일 | 커스터마이징 | P0 |
| ANSI 컬러 출력 | 가독성 | P0 |

### Phase 2 (v0.2)

| Deliverable | Impact 연결 |
|-------------|-----------|
| **번 레이트 ($/min)** | 비용 인식 |
| **Context 스파크라인** | Context 관리 |
| **Compact 예측 카운트다운** | Context 관리 |
| 조건부 위젯 (threshold) | 정보 과부하 방지 |
| 좌/우 정렬 split 레이아웃 | 가독성 |

### Phase 3 (v0.3+)

| Deliverable | Impact 연결 |
|-------------|-----------|
| Transcript 파싱 (tool/agent) | 작업 추적 |
| 5시간 블록 타이머 | 사용량 관리 |
| 세션 히스토리 영속 저장 | 장기 추적 |
| Powerline 스타일 | 미용 |
| TUI 설정 도구 | 편의성 |

### Constraints (제약 조건)

| 제약 | 설명 |
|------|------|
| 언어: Go | 성능(1-2ms startup) + 개발 생산성 균형, 빈 포지션 선점 |
| 호출 주기: ~300ms | Claude Code가 매 업데이트마다 새 프로세스 spawn |
| 입력: stdin JSON only | Claude Code statusline API 규격 |
| 출력: stdout ANSI string | 터미널 하단 표시 |
| 의존성: 최소화 | 단일 바이너리 배포 목표 |
| 설정: ~/.config/visor/ | XDG 규격 준수 |

★ = 기존 프로젝트에 없는 차별화 기능
