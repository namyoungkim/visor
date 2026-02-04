# Widget Reference

이 문서는 visor statusline의 모든 위젯에 대한 상세 레퍼런스입니다.

## 목차

1. [색상 규칙](#색상-규칙)
2. [상태 아이콘](#상태-아이콘)
3. [Core Widgets](#core-widgets)
4. [Efficiency Widgets](#efficiency-widgets)
5. [Tool/Agent Widgets](#toolagent-widgets)
6. [Rate Limit Widgets](#rate-limit-widgets)
7. [Cost Tracking Widgets](#cost-tracking-widgets)
8. [Session Info Widgets](#session-info-widgets)
9. [추천 레이아웃](#추천-레이아웃)

---

## 색상 규칙

visor는 상태를 직관적으로 파악할 수 있도록 일관된 색상 체계를 사용합니다.

| 색상 | 의미 | 사용 예시 |
|------|------|----------|
| 🟢 **Green** | 양호 - 안전 범위 | context <60%, cache_hit >80%, cost <$0.50 |
| 🟡 **Yellow** | 경고 - 주의 필요 | context 60-80%, cache_hit 50-80%, cost $0.50-1.00 |
| 🔴 **Red** | 위험 - 조치 필요 | context >80%, cache_hit <50%, cost >$1.00 |
| ⚪ **Gray/Dim** | 데이터 없음 | `—` 표시 |
| 🔵 **Cyan** | 정보/강조 | model, git ahead |
| 🟣 **Magenta** | 브랜치명 | git branch |

### 임계값 동작 방식

**일반 임계값** (높을수록 나쁨):
```
value < warning  → Green
warning ≤ value < critical → Yellow
value ≥ critical → Red
```

**역순 임계값** (높을수록 좋음 - cache_hit만 해당):
```
value ≥ good    → Green
warning ≤ value < good → Yellow
value < warning → Red
```

---

## 상태 아이콘

도구 및 에이전트 위젯에서 사용되는 상태 아이콘입니다.

| 아이콘 | 의미 | 색상 | 설명 |
|--------|------|------|------|
| `✓` | 완료 (Success) | Green | 작업이 성공적으로 완료됨 |
| `✗` | 에러 (Error) | Red | 작업이 실패함 |
| `◐` | 실행 중 (Running) | Yellow | 작업이 진행 중 |

---

## Core Widgets

기본 제공되는 핵심 위젯들입니다.

### `model`

현재 사용 중인 Claude 모델을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `Opus`, `Sonnet`, `Haiku` |
| **색상** | Cyan (고정) |
| **표시 조건** | 모델 정보가 있을 때 |

**설정 옵션**: 없음

---

### `context`

컨텍스트 윈도우 사용률을 프로그레스 바와 함께 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `Ctx: 42% ████░░░░░░` |
| **색상** | <60% Green, 60-80% Yellow, >80% Red |
| **기본 임계값** | warn=60%, critical=80% |

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `true` | "Ctx:" 접두사 표시 |
| `show_bar` | `true` | 프로그레스 바 표시 |
| `bar_width` | `10` | 프로그레스 바 너비 (문자 수) |
| `warn_threshold` | `60` | 경고 색상 임계값 (%) |
| `critical_threshold` | `80` | 위험 색상 임계값 (%) |

**설정 예시**:
```toml
[[line.widget]]
name = "context"
[line.widget.extra]
show_bar = "true"
bar_width = "15"
warn_threshold = "70"
critical_threshold = "90"
```

---

### `git`

Git 저장소의 브랜치와 상태를 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `main +3~2↑1↓2` |
| **표시 조건** | Git 저장소 내에서만 표시 |

**상태 표시자**:

| 기호 | 의미 | 색상 |
|------|------|------|
| `+N` | Staged 파일 수 | Green |
| `~N` | Modified 파일 수 | Yellow |
| `↑N` | 리모트보다 앞선 커밋 수 | Cyan |
| `↓N` | 리모트보다 뒤쳐진 커밋 수 | Red |

**설정 옵션**: 없음

---

### `cost`

현재 세션의 총 API 비용을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `$0.15`, `$1.50`, `$0.003` |
| **색상** | <$0.50 Green, $0.50-1.00 Yellow, >$1.00 Red |
| **기본 임계값** | warn=$0.50, critical=$1.00 |

**출력 포맷**:
- `cost ≥ $0.01`: `$X.XX` (소수점 2자리)
- `cost < $0.01`: `$X.XXX` (소수점 3자리)

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Cost:" 접두사 표시 |
| `warn_threshold` | `0.5` | 경고 색상 임계값 (USD) |
| `critical_threshold` | `1.0` | 위험 색상 임계값 (USD) |

---

### `cache_hit`

캐시 히트율을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `Cache: 80%`, `Cache: —` |
| **색상** | >80% Green, 50-80% Yellow, <50% Red |
| **기본 임계값** | good=80%, warn=50% |
| **데이터 없음** | `Cache: —` (Gray) |

**계산 공식**:
```
rate = cache_read_tokens / (cache_read_tokens + input_tokens) × 100
```

**의미**: 높을수록 비용 효율적입니다. 캐시된 토큰은 새로 처리하는 토큰보다 저렴합니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `true` | "Cache:" 접두사 표시 |
| `good_threshold` | `80` | 양호 색상 임계값 (%) |
| `warn_threshold` | `50` | 경고 색상 임계값 (%) |

---

### `api_latency`

API 응답 지연시간을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `API: 2.5s`, `API: 500ms`, `API: —` |
| **색상** | <2s Green, 2-5s Yellow, >5s Red |
| **기본 임계값** | warn=2000ms, critical=5000ms |
| **데이터 없음** | `API: —` (Gray) |

**출력 포맷**:
- `latency ≥ 1000ms`: 초 단위 (`X.Xs`)
- `latency < 1000ms`: 밀리초 단위 (`Xms`)

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `warn_threshold` | `2000` | 경고 색상 임계값 (ms) |
| `critical_threshold` | `5000` | 위험 색상 임계값 (ms) |

---

### `code_changes`

세션 중 코드 변경량을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `+25/-10` |
| **색상** | 추가(+) Green, 삭제(-) Red |
| **표시 조건** | 변경이 있을 때만 표시 |

**설정 옵션**: 없음

---

## Efficiency Widgets

효율성 모니터링을 위한 고급 위젯들입니다.

### `burn_rate`

분당 비용 소모율을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `64.0¢/min`, `$1.5/min`, `—` |
| **색상** | <10¢/min Green, 10-25¢/min Yellow, >25¢/min Red |
| **기본 임계값** | warn=10¢/min, critical=25¢/min |
| **표시 조건** | duration 데이터가 있을 때 |

**계산 공식**:
```
burn_rate = total_cost_usd / (total_duration_ms / 60000)
```

**출력 포맷**:
- `rate ≥ $1/min`: 달러 단위 (`$X.X/min`)
- `rate < $1/min`: 센트 단위 (`X.X¢/min`)

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Burn:" 접두사 표시 |
| `warn_threshold` | `10` | 경고 색상 임계값 (¢/min) |
| `critical_threshold` | `25` | 위험 색상 임계값 (¢/min) |

---

### `compact_eta`

80% 컨텍스트 도달까지 예상 시간을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `~18m`, `~2h`, `~1h30m`, `<1m`, `compact soon` |
| **색상** | >10m Green, 5-10m Yellow, <5m Red |
| **기본 임계값** | warn=10분, critical=5분 |
| **표시 조건** | context >40%이고 duration 데이터가 있을 때 |

**계산 공식**:
```
burn_rate_pct_per_min = current_pct / (duration_ms / 60000)
eta_minutes = (80 - current_pct) / burn_rate_pct_per_min
```

**의미**: 현재 컨텍스트 소비 속도를 기준으로 80%(compact 트리거 지점)에 도달하는 예상 시간입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "ETA:" 접두사 표시 |
| `show_when_above` | `40` | 표시 시작 context % 임계값 |

---

### `context_spark`

컨텍스트 사용률의 히스토리를 스파크라인으로 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `▂▃▄▅▆`, `▁▁▂▃▅▆▇█` |
| **문자셋** | `▁▂▃▄▅▆▇█` (8단계) |
| **색상** | 상승 추세 Red, 하락 추세 Green, 유지 Yellow |
| **표시 조건** | 히스토리 데이터가 2개 이상일 때 |

**의미**: 최근 N회의 컨텍스트 사용률 변화를 시각화합니다. 급격한 상승은 빠른 컨텍스트 소비를 의미합니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Ctx:" 접두사 표시 |
| `width` | `8` | 스파크라인 너비 (문자 수) |

---

## Tool/Agent Widgets

도구 및 서브에이전트 상태를 표시하는 위젯들입니다.

### `tools`

최근 도구 호출 상태를 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `✓Bash ×7 \| ✓Edit ×4 \| ✓Read ×6` |
| **아이콘** | ✓완료(Green), ✗에러(Red), ◐실행중(Yellow) |
| **표시 조건** | 도구 호출이 있을 때 |

**상태 표시**:
- `✓Bash ×7`: Bash 도구가 성공적으로 7회 호출됨
- `✗Edit`: Edit 도구 호출이 실패함
- `◐Read`: Read 도구가 현재 실행 중

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Tools:" 접두사 표시 |
| `max_display` | `0` | 표시할 최대 도구 수 (0=무제한) |
| `show_count` | `true` | 호출 횟수 표시 (×N) |

---

### `agents`

서브에이전트의 상태를 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `✓Explore: Analyze... (42s)`, `◐Plan: Impl... (5s...)` |
| **아이콘** | ✓완료(Green), ◐실행중(Yellow) |
| **표시 조건** | 에이전트 호출이 있을 때 |

**출력 포맷**:
- 완료된 에이전트: `✓Type: Description (Ns)` - 총 소요 시간
- 실행 중인 에이전트: `◐Type: Description (Ns...)` - 현재 경과 시간 (실시간)

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Agents:" 접두사 표시 |
| `max_display` | `0` | 표시할 최대 에이전트 수 (0=무제한) |
| `show_description` | `true` | 작업 설명 표시 |
| `show_duration` | `true` | 소요/경과 시간 표시 |
| `max_description_len` | `20` | 설명 최대 길이 (초과시 `...` 처리) |

---

## Rate Limit Widgets

Claude Pro 요금제의 사용량 제한을 모니터링하는 위젯들입니다.

### `block_timer`

5시간 블록의 남은 시간을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `Block: 4h23m`, `Block: 45m` |
| **색상** | <80% 경과 Green, 80-95% Yellow, >95% Red |
| **기본 임계값** | warn=80% 경과, critical=95% 경과 |
| **표시 조건** | 블록 시작 시간이 기록되어 있을 때 |

**의미**: Claude Pro의 5시간 사용량 제한 블록에서 남은 시간입니다. 블록이 리셋되면 사용량이 초기화됩니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `true` | "Block:" 접두사 표시 |
| `warn_threshold` | `80` | 경고 색상 임계값 (경과 %) |
| `critical_threshold` | `95` | 위험 색상 임계값 (경과 %) |

---

### `block_limit`

5시간 블록의 사용률을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `5h: 42%`, `5h: 42% (4h18m)`, `5h: 42% ████░░░░░░ (4h18m)` |
| **색상** | <70% Green, 70-90% Yellow, >90% Red |
| **기본 임계값** | warn=70%, critical=90% |
| **표시 조건** | 사용량 데이터가 있을 때 |

**의미**: 현재 5시간 블록에서 사용한 양의 비율입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `true` | "5h:" 접두사 표시 |
| `show_remaining` | `true` | 남은 시간 표시 |
| `show_bar` | `false` | 프로그레스 바 표시 |
| `bar_width` | `10` | 프로그레스 바 너비 (문자 수) |
| `warn_threshold` | `70` | 경고 색상 임계값 (%) |
| `critical_threshold` | `90` | 위험 색상 임계값 (%) |

**설정 예시 (프로그레스 바 활성화)**:
```toml
[[line.widget]]
name = "block_limit"
[line.widget.extra]
show_bar = "true"
bar_width = "10"
```

---

### `week_limit`

7일 사용률을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `7d: 69%`, `7d: 69% (3d12h)` |
| **색상** | <70% Green, 70-90% Yellow, >90% Red |
| **기본 임계값** | warn=70%, critical=90% |
| **표시 조건** | 사용량 데이터가 있을 때 |

**의미**: 7일 윈도우에서 사용한 양의 비율입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `true` | "7d:" 접두사 표시 |
| `show_remaining` | `false` | 남은 시간 표시 |
| `warn_threshold` | `70` | 경고 색상 임계값 (%) |
| `critical_threshold` | `90` | 위험 색상 임계값 (%) |

---

## Cost Tracking Widgets

비용 추적을 위한 위젯들입니다.

### `daily_cost`

오늘 누적 비용을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `$2.34`, `$15`, `$0.45` |
| **색상** | <$5 Green, $5-10 Yellow, >$10 Red |
| **기본 임계값** | warn=$5, critical=$10 |

**의미**: 오늘(00:00~현재) 사용한 총 비용입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Today:" 접두사 표시 |
| `warn_threshold` | `5.0` | 경고 색상 임계값 (USD) |
| `critical_threshold` | `10.0` | 위험 색상 임계값 (USD) |

---

### `weekly_cost`

이번 주 누적 비용을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `$15.67`, `$50`, `$5.2` |
| **색상** | <$25 Green, $25-50 Yellow, >$50 Red |
| **기본 임계값** | warn=$25, critical=$50 |

**의미**: 이번 주 사용한 총 비용입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Week:" 접두사 표시 |
| `warn_threshold` | `25.0` | 경고 색상 임계값 (USD) |
| `critical_threshold` | `50.0` | 위험 색상 임계값 (USD) |

---

### `block_cost`

현재 5시간 블록 비용을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `$0.45`, `$2.5`, `$0.08` |
| **색상** | <$2 Green, $2-5 Yellow, >$5 Red |
| **기본 임계값** | warn=$2, critical=$5 |
| **표시 조건** | 블록 시작 시간이 기록되어 있을 때 |

**의미**: 현재 5시간 블록 동안 사용한 비용입니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Block$:" 접두사 표시 |
| `warn_threshold` | `2.0` | 경고 색상 임계값 (USD) |
| `critical_threshold` | `5.0` | 위험 색상 임계값 (USD) |

---

## Session Info Widgets

세션 정보 및 메타데이터를 표시하는 위젯들입니다.

### `session_id`

현재 세션 ID를 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `a1b2c3d4`, `Session: a1b2c3d4` |
| **색상** | Gray (고정) |
| **표시 조건** | 세션 ID가 있을 때 |

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Session:" 접두사 표시 |
| `max_length` | `8` | 최대 표시 길이 (0=전체) |

---

### `duration`

현재 세션의 경과 시간을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `⏱️ 5m`, `⏱️ 1h23m`, `⏱️ 45s` |
| **색상** | Gray (고정) |
| **표시 조건** | duration 데이터가 있을 때 |

**출력 포맷**:
- `< 1분`: 초 단위 (`Xs`)
- `1분 ~ 1시간`: 분 단위 (`Xm`)
- `>= 1시간`: 시+분 단위 (`XhYm`)

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_icon` | `true` | "⏱️" 아이콘 접두사 표시 |

---

### `token_speed`

출력 토큰 생성 속도를 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `42.1 tok/s`, `85.3 tok/s`, `—` |
| **색상** | >20 tok/s Green, 10-20 Yellow, <10 Red |
| **기본 임계값** | warn=20 tok/s, critical=10 tok/s |
| **데이터 없음** | `—` (Gray) |

**계산 공식**:
```
speed = total_output_tokens / (total_api_duration_ms / 1000)
```

**의미**: API가 초당 생성하는 출력 토큰 수입니다. 높을수록 빠른 응답을 의미합니다.

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "out:" 접두사 표시 |
| `warn_threshold` | `20` | 경고 색상 임계값 (tok/s 미만) |
| `critical_threshold` | `10` | 위험 색상 임계값 (tok/s 미만) |

---

### `plan`

현재 Claude 요금제 타입을 표시합니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `Pro`, `API`, `Bedrock` |
| **색상** | Pro/Max/Team=Cyan, API=Magenta, Bedrock=Yellow |
| **표시 조건** | 항상 표시 |

**감지 로직**:
1. 모델 ID에 "bedrock" 포함 → `Bedrock`
2. OAuth 자격 증명 존재 (`~/.claude/auth.json`) → `Pro`
3. 그 외 → `API`

**참고**: Pro/Max/Team은 현재 구분할 수 없어 모두 `Pro`로 표시됩니다.

**설정 옵션**: 없음

---

### `todos`

TaskCreate/TaskUpdate 도구로 생성된 작업 진행 상황을 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `✓ All complete (5/5)`, `⊙ Implement feature (3/5)` |
| **아이콘** | ✓완료(Green), ⊙진행중(Yellow) |
| **표시 조건** | 작업이 있을 때 |

**출력 포맷**:
- 모든 작업 완료: `✓ All complete (N/N)`
- 진행 중: `⊙ {현재 작업 제목} (완료/전체)`

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_label` | `false` | "Tasks:" 접두사 표시 |
| `max_subject_len` | `30` | 작업 제목 최대 길이 |

---

### `config_counts`

Claude 설정 파일들의 구성 항목 수를 표시합니다. **visor 고유 메트릭**입니다.

| 항목 | 값 |
|------|-----|
| **출력 예시** | `2📄 3🔒 2🔌 1🪝`, `1📄 0🔒 0🔌 0🪝` |
| **색상** | Gray (고정) |
| **표시 조건** | 설정 데이터가 있을 때 |

**표시 항목**:
| 기호 | 의미 | 소스 |
|------|------|------|
| 📄 | CLAUDE.md 파일 수 | cwd부터 루트까지 |
| 🔒 | 권한 규칙 수 | `~/.claude/settings.json` permissions |
| 🔌 | MCP 플러그인 수 | `~/.claude/settings.json` mcpServers |
| 🪝 | 훅 수 | `~/.claude/settings.json` hooks |

**설정 옵션**:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `show_claude_md` | `true` | CLAUDE.md 수 표시 |
| `show_rules` | `true` | 권한 규칙 수 표시 |
| `show_mcps` | `true` | MCP 플러그인 수 표시 |
| `show_hooks` | `true` | 훅 수 표시 |

---

## 추천 레이아웃

용도별 추천 위젯 구성입니다.

### 기본 (효율성 중심)

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "context"
  [[line.widget]]
  name = "cache_hit"
  [[line.widget]]
  name = "api_latency"
  [[line.widget]]
  name = "cost"
  [[line.widget]]
  name = "git"
```

**출력 예시**: `Opus | Ctx: 42% ████░░░░░░ | Cache: 80% | API: 2.5s | $0.15 | main ↑1`

### 비용 모니터링 중심

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "context"
  [[line.widget]]
  name = "burn_rate"
  [[line.widget]]
  name = "daily_cost"
  [[line.widget]]
  name = "block_limit"
```

**출력 예시**: `Opus | Ctx: 42% ████░░░░░░ | 12.5¢/min | $2.34 | 5h: 42%`

### 개발 중심

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "context"
  [[line.widget]]
  name = "tools"
  [[line.widget]]
  name = "agents"
  [[line.widget]]
  name = "code_changes"
  [[line.widget]]
  name = "git"
```

**출력 예시**: `Opus | Ctx: 42% ████░░░░░░ | ✓Bash ×7 | ✓Edit ×4 | ◐Explore: Anal... (5s...) | +25/-10 | main +3~2`

### 미니멀

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "context"
  [line.widget.extra]
  show_bar = "false"
  [[line.widget]]
  name = "cost"
```

**출력 예시**: `Opus | Ctx: 42% | $0.15`

### 컨텍스트 추적

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "context"
  [[line.widget]]
  name = "context_spark"
  [[line.widget]]
  name = "compact_eta"
  [[line.widget]]
  name = "burn_rate"
```

**출력 예시**: `Opus | Ctx: 65% ██████░░░░ | ▂▃▄▅▆ | ~18m | 15.2¢/min`

### 세션 정보 중심

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "plan"
  [[line.widget]]
  name = "session_id"
  [[line.widget]]
  name = "duration"
  [[line.widget]]
  name = "token_speed"
  [[line.widget]]
  name = "todos"
```

**출력 예시**: `Opus | Pro | a1b2c3d4 | ⏱️ 5m | 42.1 tok/s | ⊙ Implement feature (3/5)`

### 풀 모니터링 (멀티라인)

```toml
[[line]]
  [[line.widget]]
  name = "model"
  [[line.widget]]
  name = "plan"
  [[line.widget]]
  name = "context"
  [[line.widget]]
  name = "duration"
  [[line.widget]]
  name = "cost"
  [[line.widget]]
  name = "git"

[[line]]
  [[line.widget]]
  name = "tools"
  [[line.widget]]
  name = "agents"
  [[line.widget]]
  name = "todos"
  [[line.widget]]
  name = "config_counts"
```

**출력 예시**:
- Line 1: `Opus | Pro | Ctx: 42% ████░░░░░░ | ⏱️ 15m | $0.45 | main +2~1`
- Line 2: `✓Bash ×7 | ✓Edit ×4 | ◐Explore: Analyzing (5s...) | ⊙ Task (3/5) | 2📄 3🔒 2🔌`

---

## 위젯 요약표

| 위젯 | 식별자 | 고유 | 카테고리 |
|------|--------|------|----------|
| 모델명 | `model` | | Core |
| 컨텍스트 | `context` | | Core |
| Git 상태 | `git` | | Core |
| 비용 | `cost` | | Core |
| 캐시 히트율 | `cache_hit` | ✓ | Core |
| API 지연시간 | `api_latency` | ✓ | Core |
| 코드 변경 | `code_changes` | ✓ | Core |
| 번 레이트 | `burn_rate` | ✓ | Efficiency |
| Compact ETA | `compact_eta` | ✓ | Efficiency |
| 스파크라인 | `context_spark` | ✓ | Efficiency |
| 도구 상태 | `tools` | ✓ | Tool/Agent |
| 에이전트 상태 | `agents` | ✓ | Tool/Agent |
| 블록 타이머 | `block_timer` | ✓ | Rate Limit |
| 5시간 제한 | `block_limit` | | Rate Limit |
| 7일 제한 | `week_limit` | | Rate Limit |
| 일별 비용 | `daily_cost` | | Cost Tracking |
| 주별 비용 | `weekly_cost` | | Cost Tracking |
| 블록 비용 | `block_cost` | | Cost Tracking |
| 세션 ID | `session_id` | | Session Info |
| 세션 시간 | `duration` | | Session Info |
| 토큰 속도 | `token_speed` | ✓ | Session Info |
| 요금제 | `plan` | | Session Info |
| 작업 진행 | `todos` | ✓ | Session Info |
| 설정 현황 | `config_counts` | ✓ | Session Info |

**고유(✓)**: visor만의 고유 메트릭으로, 다른 statusline에서는 제공하지 않는 정보입니다.

---

## 버전 히스토리

| 버전 | 추가된 위젯 |
|------|-------------|
| v0.1 | `model`, `context`, `git`, `cost`, `cache_hit`, `api_latency`, `code_changes` |
| v0.2 | `burn_rate`, `compact_eta`, `context_spark` |
| v0.3 | `tools`, `agents` |
| v0.4 | `block_timer` |
| v0.6 | `daily_cost`, `weekly_cost`, `block_cost`, `block_limit`, `week_limit` |
| v0.10 | `session_id`, `duration`, `token_speed`, `plan`, `todos`, `config_counts` |
