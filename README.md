# 記憶體計算工具（Memory Calculator）

**記憶體計算工具**是為了協助 Java
虛擬機（JVM）在運行時計算記憶體設定而開發的工具，基於 [paketo-buildpacks/libjvm](https://github.com/paketo-buildpacks/libjvm/)。

## 先備條件：

- Golang: v1.20+
- Jib ContainerTool
- Linux 基礎映像

## 配置項目：

在計算過程中, 需要依照特定的順序和邏輯使用多個參數, 以下是參數的取得順序及其說明：

| 參數說明 | 優先判斷 args 傳入  | 其次判斷 OS Variable | 最後的預設值或行為 |
|---|---|---|---|
| 記憶體計算工具分配的預留空間百分比 | `--head-room` | `$BPL_JVM_HEAD_ROOM` | `0` |
| 運行時將加載的 class 數量 | `--loaded-class-count ` | `$BPL_JVM_LOADED_CLASS_COUNT ` | 動態計算 App 目錄下 class 總數量的 35% |
| 運行時的用戶線程數 | `--thread-count` | `$BPL_JVM_THREAD_COUNT` | `200` |
| App 目錄 | `--application-path` | | `/app` |
| VM 建立參數 | `--jvm-options` | `$JAVA_TOOL_OPTIONS` | |
| Java啟動參數 |   | `$JAVA_OPTS ` | |
| Java Home |   | `$JAVA_HOME ` | |

## 使用方法：

1. 將 **記憶體計算工具** 和 `entrypoint.sh` 放置於 `/tmp` 目錄。
2. 自定義入口點為 `/tmp/entrypoint.sh`。
3. 避免使用 `jvmFlags`，改用 `JAVA_OPTS` 環境變數。
4. 執行 `memory-calculator --help` 閱讀完整說明

## Reference

- [Java Buildpack - Memory Calculator](https://paketo.io/docs/reference/java-reference/#memory-calculator)
