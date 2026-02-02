# Contributing to visor

visor에 기여해 주셔서 감사합니다! 이 문서는 기여 방법을 안내합니다.

## 개발 환경 설정

### 요구사항

- Go 1.22 이상
- git

### 저장소 클론

```bash
git clone https://github.com/namyoungkim/visor.git
cd visor
```

### 빌드 및 테스트

```bash
# 빌드
go build -o visor ./cmd/visor

# 테스트
go test ./...

# 테스트 (상세)
go test -v ./...

# 커버리지
go test -cover ./...
```

## 기여 방법

### 1. Issue 확인

- 기존 이슈에서 작업할 항목 찾기
- 새로운 기능/버그는 먼저 이슈 생성

### 2. 브랜치 생성

```bash
git checkout -b feature/my-feature
# 또는
git checkout -b fix/my-bugfix
```

### 3. 변경 사항 작성

- 코드 스타일 가이드 준수
- 테스트 작성
- 문서 업데이트 (필요 시)

### 4. 테스트 실행

```bash
go test ./...
go build -o visor ./cmd/visor
echo '{}' | ./visor  # 수동 테스트
```

### 5. 커밋

```bash
git add .
git commit -m "feat: add new widget for X"
```

커밋 메시지 형식:
- `feat:` 새 기능
- `fix:` 버그 수정
- `docs:` 문서 변경
- `refactor:` 리팩토링
- `test:` 테스트 추가/수정
- `chore:` 빌드, 설정 등

### 6. Pull Request

- main 브랜치로 PR 생성
- PR 템플릿에 맞게 설명 작성
- 리뷰어 피드백 반영

## 코드 스타일

### Go 코드

```go
// 좋은 예
func (w *MyWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
    if session.Data == nil {
        return ""
    }
    return render.Colorize(session.Data.Value, "green")
}

// 피해야 할 예
func (w *MyWidget) Render(s *input.Session, c *config.WidgetConfig) string {
    // 축약된 변수명 사용 금지
    // nil 체크 누락
    return render.Colorize(s.Data.Value, "green")
}
```

### 원칙

1. **명확한 변수명**: `s` 대신 `session`, `cfg` 대신 `config`
2. **nil 안전**: 포인터 접근 전 nil 체크
3. **graceful fallback**: panic 대신 빈 값 반환
4. **최소 의존성**: 표준 라이브러리 우선

### 파일 구조

```
internal/widgets/
├── widget.go           # 인터페이스, Registry
├── model.go            # 개별 위젯
├── model_test.go       # 위젯 테스트
```

## 새 위젯 추가

### 체크리스트

- [ ] `internal/widgets/` 에 위젯 파일 생성
- [ ] `Widget` 인터페이스 구현 (`Name`, `Render`, `ShouldRender`)
- [ ] `widget.go`의 `init()`에 등록
- [ ] 테스트 파일 작성
- [ ] README.md 위젯 테이블 업데이트
- [ ] CHANGELOG.md에 추가

### 템플릿

```go
package widgets

import (
    "github.com/namyoungkim/visor/internal/config"
    "github.com/namyoungkim/visor/internal/input"
    "github.com/namyoungkim/visor/internal/render"
)

type MyWidget struct{}

func (w *MyWidget) Name() string {
    return "my_widget"
}

func (w *MyWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
    // 구현
}

func (w *MyWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
    // 구현
}
```

## 테스트 가이드

### 테스트 파일 위치

```
internal/widgets/cache_hit.go      # 구현
internal/widgets/cache_hit_test.go # 테스트
```

### 테스트 패턴

```go
func TestMyWidget_Scenario(t *testing.T) {
    w := &MyWidget{}
    session := &input.Session{
        // 테스트 데이터
    }

    result := w.Render(session, &config.WidgetConfig{})

    // ANSI 코드 포함 시 strings.Contains 사용
    if !strings.Contains(result, "expected") {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}
```

### 테스트 케이스

- 정상 데이터
- nil/빈 데이터
- 경계값 (0, 음수, 최대값)
- 색상 코드 확인

## 문서 업데이트

변경 사항에 따라 업데이트:

| 변경 유형 | 업데이트 문서 |
|-----------|---------------|
| 새 위젯 | README.md, 05_IMPLEMENTATION.md, CHANGELOG.md |
| CLI 옵션 | README.md, CHANGELOG.md |
| 설정 변경 | README.md, 05_IMPLEMENTATION.md |
| 버그 수정 | CHANGELOG.md |

## 질문 및 도움

- 이슈 생성
- PR 코멘트

감사합니다!
