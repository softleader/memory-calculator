# Spec: 修正 issue #67 的 Go 套件風險

## 狀態

- 階段：Plan
- 對應議題：[softleader/memory-calculator#67](https://github.com/softleader/memory-calculator/issues/67)
- 實作閘門：計畫需經使用者確認後，才能進入 Implement。

## 假設

1. 主要目標是消除已經有修復版本的 Harbor 掃描風險，同時維持既有 CLI、JVM 記憶體計算、Spring Boot 最佳化與產物行為相容。
2. 主要使用者是本專案維護者、發布人員，以及使用 Linux／macOS 發布產物的下游團隊。
3. 為降低變更風險，採用可解析的最低安全版本集合：`golang.org/x/crypto v0.52.0`、`golang.org/x/net v0.55.0`、`golang.org/x/sys v0.45.0`。`x/sys v0.45.0` 是前兩個目標模組的 solver 必要版本，高於 issue 所列下限 `v0.44.0`。
4. 保持 `go 1.25.0` 與 `toolchain go1.25.9`，不因本次修正主動升級 Go toolchain。
5. `GO-2026-5932` 目前沒有上游修復版本，因此本次必須明確記錄為已知殘餘風險並持續追蹤，不得宣稱已修復或隱藏掃描結果。
6. 實作範圍包含相依版本更新、必要的測試補強及驗證；本 repo 不直接發布 Harbor 產物或觸發 Harbor 重掃，下游維護者負責重建、掃描與結果回填。
7. 完成定義採「Issue #67 範圍完成」，不要求本次一併修復既有 Go 標準庫與 `go-pkcs12` 弱點，也不得將此結果描述為整體安全掃描乾淨。
8. 本 repo 的下游交付物固定為 GitHub patch release `1.2.6`；下游不得以浮動的 `latest`、`main` 或未標記 commit 作為 Harbor 回填依據。

## Objective

將下列間接 Go 依賴升級到 issue #67 指定的最低安全版本，並以針對實際使用路徑的測試證明升級沒有破壞本專案的核心功能：

| 模組 | 現行版本 | 目標版本下限 | 本專案主要使用路徑 |
|---|---:|---:|---|
| `golang.org/x/crypto` | `v0.45.0` | `v0.52.0` | `libjvm` → `go-pkcs12` → `pbkdf2` 憑證載入 |
| `golang.org/x/net` | `v0.47.0` | `v0.55.0` | `miekg/dns` → Linux IPv4／IPv6 DNS 與 socket |
| `golang.org/x/sys` | `v0.41.0` | `v0.45.0` | `miekg/dns`／`x/net` → Linux `unix` 系統呼叫 |

### Scope

- 更新 `go.mod` 與 `go.sum` 中三個目標模組及 solver 必要的間接依賴。
- 補強完整 CLI 編排、憑證載入及 Linux 契約層測試。
- 讓目前未被 Go test runner 執行的 `WebApplicationType` specs 正常執行。
- 執行單元／整合測試、race detector、vet、跨平台 build、coverage 與弱點檢查。
- 建立 GitHub patch release `1.2.6` 作為下游重建的固定交付物，並記錄下游回填的 Harbor 重掃結果。

### Out of Scope

- 與 issue #67 無關的依賴全面升級、重構或功能調整。
- Go toolchain `1.25.9` 與 `software.sslmate.com/src/go-pkcs12 v0.6.0` 的既有弱點修正；其掃描結果需保留為後續事項。
- 在上游尚無修復版本前，自行修改或 fork `GO-2026-5932` 所屬模組。
- 直接發布、推送或修改下游 Harbor 產物；本 repo 只接收下游驗收結果。
- 僅為提高 coverage 數字而測試第三方套件的內部實作；測試應驗證本專案可觀察行為。

## Commands

以下命令皆從專案根目錄執行。

### Baseline

```sh
GOTOOLCHAIN=go1.25.9 go mod download
GOTOOLCHAIN=go1.25.9 go test -count=1 ./...
GOTOOLCHAIN=go1.25.9 go test -count=1 -race ./...
GOTOOLCHAIN=go1.25.9 go vet ./...
```

### Coverage

預設 package-local coverage 僅有 `49.1%`；跨 package 合併基準為 `53.3%`，本次以後者作為相同平台的回歸比較基準。

```sh
GOTOOLCHAIN=go1.25.9 go test -count=1 -covermode=atomic -coverpkg='./...' -coverprofile=/tmp/memory-calculator-coverage.out ./...
GOTOOLCHAIN=go1.25.9 go tool cover -func=/tmp/memory-calculator-coverage.out
```

### Dependency Update

```sh
GOTOOLCHAIN=go1.25.9 go get golang.org/x/crypto@v0.52.0 golang.org/x/net@v0.55.0 golang.org/x/sys@v0.45.0
GOTOOLCHAIN=go1.25.9 go mod tidy
GOTOOLCHAIN=go1.25.9 go mod tidy -diff
GOTOOLCHAIN=go1.25.9 go list -m all | rg '^golang.org/x/(crypto|net|sys) '
```

### Build Matrix

```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.25.9 go build ./...
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOTOOLCHAIN=go1.25.9 go build ./...
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.25.9 go build ./...
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 GOTOOLCHAIN=go1.25.9 go build ./...
```

### Vulnerability Verification

```sh
GOTOOLCHAIN=go1.25.9 go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
```

`govulncheck` 是原始碼可達性檢查，不能取代發布產物的 Harbor package inventory 掃描。兩者結果都必須保留。本次 baseline 已有 6 個 symbol-level findings；Issue #67 範圍完成要求升級後不得新增 findings，但不以清除這 6 個既有 findings 為完成條件。

## Project Structure

- `go.mod`／`go.sum`：目標 module closure。
- `boot/helper/*_test.go`：啟用既有 specs。
- `main*_test.go`、`calc/*_test.go`、`calc/testdata/`：最小 production-path 契約測試與 test-only fixtures。

## Code Style

- 所有 Go 程式碼必須通過 `gofmt` 與 `go vet ./...`。
- 沿用標準 library-first、明確 error handling、早期 return，以及 `%w` 包裝底層錯誤的寫法。
- 測試名稱採 `Test<Subject>_<Scenario>`，使用 Arrange／Act／Assert，驗證輸出或狀態，不只驗證 helper 是否存在。
- JVM options 採混合相容契約：canonical 案例鎖定完整字串與順序，其餘案例驗證必要參數、優先順序與語意；debug log 不屬於契約。
- 測試環境變數優先使用 `t.Setenv`，暫存資料優先使用 `t.TempDir`，避免跨測試洩漏狀態。
- 新測試優先使用標準 `testing`；既有 spec 檔可繼續使用 `sclevine/spec` 與 Gomega，不新增測試框架。
- 不手動編輯 `go.sum`；只由 Go module commands 產生。

新增的整合測試必須斷言可觀察結果，例如精確的 `JAVA_TOOL_OPTIONS`、輸出檔內容或憑證載入結果。

## Testing Strategy

### Baseline Evidence

- 現有 74 個 Go tests 全部通過。
- 同一平台的跨 package statement coverage 基準為 `53.3%`。
- `calc.Calculator.Execute` 為 `77.8%`，但 `main`／`run`／`out`、`SpringOptimizer.Execute` 目前為 `0%`。
- 將 `x/crypto`、`x/net`、`x/sys` 納入 instrumentation 時，現有測試對其執行 coverage 為 `0%`。

### Test Levels

1. **Unit tests**
   - 保留現有參數、環境變數、contributor 與 preparer 測試。
   - 為錯誤分支及依賴升級可能影響的可觀察結果補上精確斷言。
   - 為 `boot/helper/web_application_type_test.go` 新增頂層 `TestWebApplicationType`，確保現有 specs 確實被執行。

2. **CLI／orchestration integration test**
   - 在可控制的暫存目錄與記憶體限制下執行完整 `prep → boot → calc → out`。
   - 至少涵蓋 `--loaded-class-count` 路徑、輸出到 stdout 或檔案，以及非 preview 的預設流程。
   - 一個 canonical 案例斷言完整、順序固定的 `JAVA_TOOL_OPTIONS` 字串。
   - 其他案例斷言必要參數、重複參數的優先順序與輸出語意；只有「輸出非空」不算通過。
   - 不鎖定 debug log 或其他非契約文字。

3. **Certificate integration test (`x/crypto`)**
   - 使用明確標示為測試用途、不得用於正式環境的 fixture。
   - 設定有效的 JVM cacerts／certificate input，實際執行 `Calculator.Execute` 或等價公開流程。
   - 確認 `OpenSSLCertificateLoader` 被執行、沒有錯誤，且可觀察輸出符合預期。
   - 若實作需要閱讀第三方原始碼，必須先取得與 runtime 完全相同的 dependency tag，不能用 mirror 的 `main` 推測。

4. **Linux contract integration test (`x/net`, `x/sys`)**
   - 在 Linux amd64 環境執行完整 `Calculator.Execute` 與既有 Ubuntu CI 流程。
   - 驗證目前 resolver 設定不造成錯誤，並斷言 Memory Calculator 可控制的最終 JVM options。
   - 不新增 resolver injection seam，不要求特權網路環境，也不直接斷言第三方 IPv4／IPv6 socket 內部行為。

5. **Static／runtime verification**
   - `go vet ./...`、`go test -race ./...`、四個發布平台 build 均需通過。
   - 同平台 `-coverpkg='./...'` coverage 不得低於 `53.3%`；新增測試應提高關鍵路徑覆蓋，不以追求任意全域百分比取代行為斷言。
   - 以固定的 `govulncheck v1.6.0` 比較升級前後結果，不得新增 symbol-level findings；既有 6 個 findings 明列為本次非阻斷的後續事項。

6. **Downstream artifact security verification**
   - 下游維護者取得 GitHub Release `1.2.6` 後，重建其實際映像並執行 Harbor 重掃。
   - issue #67 中已有修復版本的 CVE 不得再次出現。
   - 下游維護者將 image tag／digest、掃描時間與結果回填至 issue #67。
   - `GO-2026-5932` 影響 `x/crypto/openpgp`，目前不在本專案可達路徑，但 Harbor 若仍以 module-level finding 回報，需保留為已知殘餘風險。
   - GitHub Release 驗收需記錄實際 binary build info；現有 workflow 使用 Go `1.25` 範圍，本次不修改 workflow 以固定 patch。

## Boundaries

### Always Do

- 使用 `go1.25.9` 重現 baseline，並在升級後重跑相同命令。
- 採最小必要依賴變更，檢查完整 `go.mod`／`go.sum` diff 與 module graph。
- 先新增或啟用能證明目標 production path 的測試，再將依賴升級視為完成。
- 保持 CLI flags、環境變數、輸出格式及既有發布平台相容。
- 記錄 `govulncheck` 與 Harbor 結果的差異及 `GO-2026-5932` 狀態。
- 若需要閱讀 dependency source，先驗證並使用完全相同的 runtime version tag。
- 若規格或範圍改變，先更新本文件並重新取得確認。

### Ask First

- 升級 Go language version 或 toolchain。
- 升級直接依賴，或升級與三個目標模組無 solver 必要關係的套件。
- 修改 production logic、公開 CLI、環境變數、輸出格式或 error semantics。
- 新增第三方測試／建置依賴，或修改 GitHub Actions、Dockerfile、GoReleaser 與發布流程。
- 加入需要網路、特權容器或正式憑證的測試。
- 建立 commit、push、PR、tag 或 GitHub Release。

### Never Do

- 手動竄改 `go.sum`、略過 module checksum 或關閉安全驗證。
- 移除、跳過或弱化失敗測試來讓升級通過。
- 隱藏、抑制或宣稱已修復尚無修復版本的 `GO-2026-5932`。
- 將真實私鑰、正式憑證、token、registry credential 或其他秘密提交到 repository。
- 為本次修正順手進行無關重構、全面 dependency refresh 或破壞性 Git 操作。
- 直接操作下游 Harbor，或將「原始碼測試通過」誤報為「Harbor 已驗收」。

## Success Criteria

- [ ] `go list -m all` 顯示 `golang.org/x/crypto >= v0.52.0`。
- [ ] `go list -m all` 顯示 `golang.org/x/net >= v0.55.0`。
- [ ] `go list -m all` 顯示 `golang.org/x/sys >= v0.45.0`。
- [ ] 除 solver 必要調整外，沒有無關依賴、Go toolchain 或 production behavior 變更。
- [ ] 既有 74 個 tests 與新增 tests 全部通過。
- [ ] `WebApplicationType` 現有 specs 已由 Go test runner 執行。
- [ ] CLI／orchestration canonical 測試鎖定完整 JVM options 字串與順序，其他案例驗證必要參數、優先順序與輸出語意。
- [ ] 有效憑證案例實際執行 certificate loader 並通過。
- [ ] Linux 契約層測試已執行完整 `Calculator.Execute`，並驗證無 resolver error 與穩定的 JVM options。
- [ ] 相同平台的跨 package coverage 不低於 baseline `53.3%`，且關鍵 0% 路徑已有行為測試。
- [ ] `go vet ./...`、`go test -race ./...` 與 Linux／macOS amd64／arm64 builds 全部通過。
- [ ] 固定版 `govulncheck v1.6.0` 的前後結果已保存，且沒有因本次升級新增 symbol-level findings。
- [ ] Go `1.25.9` 與 `go-pkcs12 v0.6.0` 的 6 個既有 symbol-level findings 已記錄為 Issue #67 範圍外的後續事項。
- [ ] `GO-2026-5932` 的無修復狀態與追蹤方式已明確記錄。
- [ ] GitHub patch release `1.2.6` 已建立，四個 GoReleaser 平台產物與 checksums 均存在。
- [ ] `1.2.6` binary 的實際 Go build info 已記錄於驗收證據。
- [ ] 下游維護者已使用 `1.2.6` 完成映像重建與 Harbor 重掃，回填結果不再包含 issue #67 中已有修復版本的 CVE。
