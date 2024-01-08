# 記憶體計算工具（Memory Calculator）

**記憶體計算工具**是為了協助 Java
虛擬機（JVM）在運行時計算記憶體設定而開發的工具，基於 [paketo-buildpacks/libjvm](https://github.com/paketo-buildpacks/libjvm/)。

## 先備條件：

- Golang: v1.20+
- Jib ContainerTool
- Linux 基礎映像

## 配置項目：

- `$BPL_JVM_HEAD_ROOM`：記憶體計算工具分配的預留空間百分比，預設為 `0`。
- `$BPL_JVM_LOADED_CLASS_COUNT`：運行時將加載的類的數量，預設為總類數的35%。
- `$BPL_JVM_THREAD_COUNT`：運行時的用戶線程數，預設為 `200`。
- `$JAVA_HOME`：JRE的安裝位置。
- `$JAVA_OPTS`：Java啟動選項。
- `$JAVA_TOOL_OPTIONS`：JVM啟動選項(由JVM提供)。

## 使用方法：

1. 將 **記憶體計算工具** 和 `entrypoint.sh` 放置於 `/tmp` 目錄。
2. 自定義入口點為 `/tmp/entrypoint.sh`。
3. 避免使用 `jvmFlags`，改用 `JAVA_OPTS` 環境變數。
