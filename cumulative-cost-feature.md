# Feature: 누적 사용량/비용 추적 (Cumulative Usage & Cost Tracking)

## 개요

Claude Code statusline이 제공하는 `total_cost_usd`는 **현재 세션**의 비용만 포함한다. 일별, 주별, 월별 누적 사용량/비용을 확인하려면 Claude Code가 기록하는 JSONL 로그 파일을 직접 파싱하거나, OAuth API를 통해 사용률 한도를 조회해야 한다.

## 배경

### 사용 형태별 차이점

| 형태 | 과금 방식 | `/cost` 의미 | 핵심 추적 지표 |
|:-----|:---------|:------------|:--------------|
| **구독 (Pro/Max)** | 월정액 $20~$200 | 추정치 (실제 청구 아님) | 5h/7d 사용률 한도 |
| **Anthropic API** | 종량제 | 실제 비용 | 토큰 × Anthropic 가격 |
| **Vertex AI** | 종량제 | 부정확 (가격 다름) | 토큰 × GCP 가격 |
| **AWS Bedrock** | 종량제 | 부정확 (가격 다름) | 토큰 × AWS 가격 |

> ⚠️ **구독 사용자 주의**: `/cost`가 "$19.64"라고 표시해도 실제 청구는 $0. 구독자에게 중요한 건 **사용률 한도**임.

### 현재 상태

Claude Code는 statusline에 다음 JSON을 stdin으로 전달:

```json
{
  "cost": {
    "total_cost_usd": 0.01234,      // 현재 세션만 (API 가격 기준 추정)
    "total_duration_ms": 45000,
    "total_api_duration_ms": 2300
  }
}
```

### 사용자 니즈

**구독 사용자 (Pro/Max):**
- "5시간 블록 얼마나 썼지?" → 블록 사용률 + 잔여 시간
- "7일 한도 얼마나 남았지?" → 주간 사용률 + 리셋 시간

**종량제 사용자 (API/Vertex/Bedrock):**
- "오늘 얼마나 썼지?" → 일별 비용
- "이번 주 사용량은?" → 주별 비용
- "이번 달 총 비용은?" → 월별 비용

## 기술 분석

### 접근법 1: OAuth API (구독 사용자용)

구독자의 5시간/7일 사용률 한도는 Anthropic OAuth API로 조회 가능:

```bash
GET https://api.anthropic.com/api/oauth/usage
Authorization: Bearer {access_token}
anthropic-beta: oauth-2025-04-20
```

**응답:**
```json
{
  "five_hour": {
    "utilization": 0.42,        // 42% 사용
    "resets_at": "2025-01-15T15:00:00Z"
  },
  "seven_day": {
    "utilization": 0.69,        // 69% 사용
    "resets_at": "2025-01-20T00:00:00Z"
  }
}
```

**인증 토큰 위치:**
- macOS: Keychain → "Claude Code-credentials"
- Linux: `~/.config/claude/credentials.json` (추정)
- Windows: Credential Manager (추정)

### 접근법 2: JSONL 파싱 (종량제 사용자용)

#### Claude Code 로그 구조

```
~/.claude/projects/
├── {project-hash-1}/
│   ├── {session-id-1}.jsonl
│   ├── {session-id-2}.jsonl
│   └── ...
├── {project-hash-2}/
│   └── ...
└── ...
```

### JSONL 파일 포맷

각 라인은 하나의 API 호출 기록:

```json
{
  "timestamp": "2025-01-15T10:30:00Z",
  "model": "claude-sonnet-4-20250514",
  "usage": {
    "input_tokens": 1500,
    "output_tokens": 500,
    "cache_creation_input_tokens": 0,
    "cache_read_input_tokens": 1200
  }
}
```

### 비용 계산 공식

```
cost = (input_tokens * input_price) 
     + (output_tokens * output_price)
     + (cache_creation_tokens * cache_creation_price)
     + (cache_read_tokens * cache_read_price)
```

**모델별 가격 (per 1M tokens, 2025년 1월 기준):**

| Model | Input | Output | Cache Write | Cache Read |
|:------|------:|-------:|------------:|-----------:|
| Claude Sonnet 4 | $3.00 | $15.00 | $3.75 | $0.30 |
| Claude Opus 4 | $15.00 | $75.00 | $18.75 | $1.50 |
| Claude Haiku 3.5 | $0.80 | $4.00 | $1.00 | $0.08 |

> ⚠️ 가격은 변경될 수 있음. 설정 파일 또는 API로 최신 가격 반영 필요.

## 구현 계획

### Track A: 구독 사용자용 (OAuth API)

**Phase A1: 인증 토큰 추출**

```go
// internal/auth/credentials.go
type Credentials struct {
    AccessToken      string
    RefreshToken     string
    ExpiresAt        time.Time
    SubscriptionType string  // "pro", "max", "max20"
}

func GetCredentials() (*Credentials, error)  // OS별 분기
```

**Phase A2: 사용률 조회**

```go
// internal/usage/api.go
type UsageLimits struct {
    FiveHour struct {
        Utilization float64
        ResetsAt    time.Time
    }
    SevenDay struct {
        Utilization float64
        ResetsAt    time.Time
    }
}

func FetchUsageLimits(token string) (*UsageLimits, error)
```

**Phase A3: 위젯 구현**

| 위젯 | 표시 예시 | 설명 |
|:-----|:---------|:-----|
| `BlockLimitWidget` | `5h: 42% (2h30m left)` | 5시간 블록 사용률 |
| `WeekLimitWidget` | `7d: 69% (3d left)` | 7일 사용률 |

---

### Track B: 종량제 사용자용 (JSONL 파싱)

**Phase B1: 핵심 파서 구현**

**Phase B1: 핵심 파서 구현**

**파일:** `internal/cost/parser.go`

```go
type UsageRecord struct {
    Timestamp   time.Time
    Model       string
    InputTokens int64
    OutputTokens int64
    CacheCreate int64
    CacheRead   int64
}

type CostAggregator struct {
    Today     float64
    Week      float64
    Month     float64
    FiveHour  float64  // 현재 5시간 블록
}

func ParseJSONLFiles(projectsDir string) ([]UsageRecord, error)
func Aggregate(records []UsageRecord, now time.Time) CostAggregator
```

**Phase B2: 위젯 구현**

| 위젯 | 표시 예시 | 설명 |
|:-----|:---------|:-----|
| `DailyCostWidget` | `$2.34 today` | 오늘 누적 비용 |
| `WeeklyCostWidget` | `$15.67 week` | 이번 주 누적 비용 |
| `BlockCostWidget` | `$0.45 block (2h30m left)` | 5시간 블록 비용 |

**Phase B3: 성능 최적화**

대량 파일 처리 시 성능 고려:

1. **증분 파싱**: 마지막 파싱 이후 변경된 파일만 처리
2. **캐싱**: 계산 결과를 임시 파일에 저장
3. **병렬 처리**: 여러 프로젝트 디렉토리 동시 스캔

**목표 성능:** 100개 세션 파일 기준 < 50ms

## 설정

`config.toml`에서 활성화:

```toml
[usage]
# 사용 형태 자동 감지 또는 명시적 지정
# auto | subscription | api | vertex | bedrock
provider = "auto"

# 구독 사용자용 (Track A)
[widgets.block_limit]
enabled = true
format = "5h: %.0f%% (%s left)"
warn_threshold = 80  # 80% 이상이면 경고색

[widgets.week_limit]
enabled = true
format = "7d: %.0f%%"

# 종량제 사용자용 (Track B)
[widgets.daily_cost]
enabled = true
format = "$%.2f today"

[widgets.block_cost]
enabled = true
show_remaining_time = true

# 가격 설정 (종량제 전용)
[cost]
# Provider별 가격 오버라이드
# vertex와 bedrock은 regional endpoint 10% 프리미엄 자동 적용
[cost.pricing.anthropic]
# Anthropic API 기본 가격 사용

[cost.pricing.vertex]
# GCP 가격 (regional endpoint 기준)
premium_multiplier = 1.10

[cost.pricing.bedrock]
# AWS 가격 (기본값 사용 또는 오버라이드)
```

### Provider 자동 감지 로직

```go
func DetectProvider() Provider {
    if os.Getenv("CLAUDE_CODE_USE_VERTEX") == "1" {
        return ProviderVertex
    }
    if os.Getenv("CLAUDE_CODE_USE_BEDROCK") == "1" {
        return ProviderBedrock
    }
    if os.Getenv("ANTHROPIC_API_KEY") != "" {
        return ProviderAPI
    }
    // OAuth credentials 존재 여부로 구독 감지
    if creds, _ := GetCredentials(); creds != nil {
        return ProviderSubscription
    }
    return ProviderUnknown
}
```

## 의존성

**순수 Go 표준 라이브러리:**
- `encoding/json`: JSONL 파싱, API 응답 파싱
- `filepath.WalkDir`: 디렉토리 순회
- `time`: 시간 기반 집계
- `net/http`: OAuth API 호출
- `os/exec`: Keychain 접근 (macOS)

**외부 의존성 없음** — 단일 바이너리 유지

## 리스크 및 고려사항

| 리스크 | 대응 |
|:-------|:-----|
| OAuth API 변경/제한 | 버전 헤더 관리, fallback to JSONL |
| 인증 토큰 접근 (OS별 차이) | OS별 credential 추출 로직 분기 |
| JSONL 포맷 변경 | 버전별 파서 분기, graceful degradation |
| Vertex/Bedrock 가격 차이 | Provider별 가격 테이블, 설정 오버라이드 |
| 대량 파일로 인한 지연 | 캐싱, 증분 파싱, 백그라운드 처리 |
| 권한 문제 | 에러 핸들링, 부분 결과 반환 |

## 수용 기준 (Acceptance Criteria)

### Track A (구독 사용자)
- [ ] OAuth API로 5시간/7일 사용률 조회
- [ ] macOS Keychain에서 인증 토큰 추출
- [ ] Linux/Windows credential 지원 (best effort)
- [ ] 사용률 기반 위젯 (BlockLimitWidget, WeekLimitWidget)
- [ ] API 실패 시 graceful degradation

### Track B (종량제 사용자)
- [ ] `~/.claude/projects/` 하위 모든 JSONL 파일 파싱
- [ ] 일별/주별/월별/5시간 블록별 비용 집계
- [ ] Provider별 가격 적용 (Anthropic/Vertex/Bedrock)
- [ ] 100개 세션 기준 50ms 이내 처리
- [ ] 파싱 실패 시 graceful degradation

### 공통
- [ ] Provider 자동 감지
- [ ] 단위 테스트 커버리지 80% 이상

## 우선순위

**Priority:** P2 (MVP 이후)

MVP에서는 stdin의 `total_cost_usd` (세션 비용)만 사용하고, 이 기능은 v0.2에서 구현.

## 참고 자료

- [OAuth Usage API 구현 예시](https://codelynx.dev/posts/claude-code-usage-limits-statusline) - TypeScript 기반 사용률 조회
- [ccusage 구현](https://github.com/ryoppippi/ccusage) - TypeScript 기반 비용 추적
- [Claude Code 공식 문서](https://docs.anthropic.com/en/docs/claude-code) - statusline JSON 스펙
- [Anthropic 가격 정책](https://www.anthropic.com/pricing)
- [Vertex AI 가격](https://cloud.google.com/vertex-ai/generative-ai/pricing)
- [AWS Bedrock 가격](https://aws.amazon.com/bedrock/pricing/)
