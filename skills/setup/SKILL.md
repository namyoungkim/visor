---
name: visor-setup
description: visor 설치 및 설정 자동화. 설치 여부 확인, 버전 업데이트, preset 설정. Triggers: "visor 설치", "visor 설정", "statusline 설정", "visor setup", "visor update".
allowed-tools: Bash, Read, Edit
---

# visor Setup

Claude Code efficiency dashboard 설치 및 설정.

## Workflow

### Step 1: 설치 상태 확인

```bash
which visor && visor --version
```

### Step 2: 설치 또는 업데이트

#### 미설치 시

```bash
# 플랫폼 감지
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"

# 최신 버전 확인
LATEST=$(curl -sI https://github.com/namyoungkim/visor/releases/latest | grep -i ^location | sed 's/.*tag\/v//' | tr -d '\r\n')

# 다운로드 및 설치
curl -sL "https://github.com/namyoungkim/visor/releases/download/v${LATEST}/visor_${LATEST}_${OS}_${ARCH}.tar.gz" | tar xz
sudo mv visor /usr/local/bin/ || { mkdir -p ~/.local/bin && mv visor ~/.local/bin/; }

visor --version
```

#### 설치됨 - 버전 확인

```bash
CURRENT=$(visor --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
LATEST=$(curl -sI https://github.com/namyoungkim/visor/releases/latest | grep -i ^location | sed 's/.*tag\/v//' | tr -d '\r\n')

if [[ "$CURRENT" != "$LATEST" ]]; then
  echo "Update available: $CURRENT → $LATEST"
  # 위 설치 스크립트로 업데이트
fi
```

### Step 3: 설정

#### 3a. Preset 선택

사용자에게 질문하여 preset 선택:

| Preset | 용도 | Widgets |
|--------|------|---------|
| `full` | **전체 위젯, 7줄 멀티라인 (권장)** | 24 |
| `minimal` | 필수 정보만 | 4 |
| `default` | 기본 | 6 |
| `efficiency` | 비용 최적화 | 6 |
| `developer` | Tool/agent 모니터링 | 7 |
| `pro` | Claude Pro 한도 | 6 |

```bash
visor --init <selected-preset>
```

#### 3b. Claude Code 연동

`~/.claude/settings.json`에 추가:

```json
{
  "statusline": {
    "command": "visor"
  }
}
```

### Step 4: 검증

```bash
visor --check
echo '{"model":{"display_name":"Opus"}}' | visor
```

## Checklist

- [ ] visor 설치 또는 업데이트
- [ ] `visor --init <preset>` 실행
- [ ] `~/.claude/settings.json` 연동
- [ ] `visor --check` 검증
