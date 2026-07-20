# Memory Calculator

這個 context 定義 Memory Calculator 安全維護與發布驗收時使用的專案語言。

## Language

**Issue #67 範圍完成**:
issue #67 指定的 Go 模組已解析至安全版本、相容性驗證已通過，且下游維護者已完成產物重建、Harbor 重掃與結果回填，不再包含該 issue 中已有修復版本的 CVE。其他既有弱點必須明確記錄，但不阻擋此範圍完成。
_Avoid_: 整體安全乾淨、零弱點發布

**原始碼驗收完成**:
依賴版本、測試、靜態檢查與建置均符合 issue #67 的原始碼驗收條件，但尚未代表產物已發布或通過 Harbor 重掃。
_Avoid_: Issue 完成、發布完成

**契約層測試**:
在 Linux 上驗證 Memory Calculator 可控制的輸入、編排與可觀察輸出，不要求直接控制或斷言第三方 DNS／socket 套件的內部行為。
_Avoid_: 第三方內部測試、完整網路整合測試

**下游 Harbor 驗收**:
下游維護者使用本 repo 的已驗證交付物重建其映像、執行 Harbor 掃描，並將結果回填至 issue #67；本 repo 不直接發布或操作該 Harbor 產物。
_Avoid_: 本 repo Harbor 發布、本 repo Harbor 重掃

**修復版 1.2.6**:
本 repo 針對 issue #67 發布的 patch release；下游必須以此固定版本重建映像並回填 Harbor 結果。
_Avoid_: latest、main、未標記的 source commit

**JVM options 相容契約**:
以一個 canonical 案例固定完整 `JAVA_TOOL_OPTIONS` 字串與參數順序，其他案例則驗證必要參數、優先順序與輸出語意；debug log 文字不屬於契約。
_Avoid_: 只驗證非空輸出、鎖定所有 log 文字
