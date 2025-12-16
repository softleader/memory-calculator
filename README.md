[![Go Report Cart](https://goreportcard.com/badge/github.com/softleader/memory-calculator)](https://goreportcard.com/report/github.com/softleader/memory-calculator)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/softleader/memory-calculator)
[![GitHub License](https://img.shields.io/github/license/softleader/memory-calculator)](./LICENSE)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/softleader/memory-calculator/test.yaml)](https://github.com/softleader/memory-calculator/actions)
[![release](https://img.shields.io/github/release/softleader/memory-calculator.svg)](https://github.com/softleader/memory-calculator/releases)
![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/softleader/memory-calculator/latest/total)

# 記憶體計算工具（Memory Calculator）

**記憶體計算工具** 是為了協助 Java
虛擬機（JVM）在運行時計算記憶體設定而開發的工具，基於 [paketo-buildpacks/libjvm](https://github.com/paketo-buildpacks/libjvm/)。

## Calculation Algorithm

### 基本計算公式

JVM 記憶體的計算公式如下：

```markdown
Heap = Total Container Memory - Non-Heap - Headroom
```

- **Total Container Memory**: 應用程式可用的總記憶體，通常是 Container 設置的記憶體上限值
- **Non-Heap Memory**: Non-Heap 記憶體，詳細於後面說明
- **Headroom**: Container 總記憶體的一定比例，用於非 JVM 的操作。預設值為 0。

### Non-Heap 記憶體計算

Non-Heap Memory 的計算公式如下：

```markdown
Non-Heap = Direct Memory + Metaspace + Reserved Code Cache + (Thread Stack * Thread Count)
```

下表列出了各個參數對應的 JVM Falg 及其預設值：

| 記憶體區域 | JVM Flag | Default |
|-----------|----------|---------|
| Direct Memory | `-XX:MaxDirectMemorySize` | 10MB（JVM Default） |
| Metaspace | `-XX:MaxMetaspaceSize` | 基於 Class 計數自動計算 |
| Reserved Code Cache | `-XX:ReservedCodeCacheSize` | 240MB（JVM Default） |
| Thread Stack | `-Xss` | 1M * 250（JVM Default Thread Stack 大小 * Tomcat 預設的最佳 Thread Count） |

### 結果與應用

- 在計算完上述各值後，剩餘的記憶體將被分配給 `-Xmx` Flag，作為 JVM 的 Heap 記憶體
- 所有 Flag 和值將被放到 `JAVA_TOOL_OPTIONS` Flag 中，讓應用程式使用

若你想獲得更深入或更完整的說明，請點選以下連結：[Java Buildpack - Memory Calculator](https://paketo.io/docs/reference/java-reference/#memory-calculator)

## Input Variables

在計算過程中，需要依照特定的順序和邏輯使用多個參數，以下是參數的取得順序及其說明：

| 參數說明 | 優先判斷 args 傳入  | 其次判斷 OS Variable | 最後的預設值或行為 |
|---|---|---|---|
| 記憶體計算工具分配的預留空間百分比 | `--head-room` | `$BPL_JVM_HEAD_ROOM` | `0` |
| 運行時的用戶線程數 | `--thread-count` | `$BPL_JVM_THREAD_COUNT` | `200` |
| 運行時將加載的 class 數量 | `--loaded-class-count` | `$BPL_JVM_LOADED_CLASS_COUNT` | 若沒提供，則以 App 目錄，JVM class 數量，JVM class 數量調整，於啟動時動態的計算出建議值 |
| App 目錄 | `--app-path` | `$BPI_APPLICATION_PATH` | `/app` |
| App classes 目錄 | `--app-classes-path` | `$APPLICATION_CLASSES_PATH` | `/app/classes` |
| App lib 目錄 | `--app-lib-path` | `$APPLICATION_LIB_PATH` | `/app/libs` |
| JVM class 數量 | `--jvm-class-count` | `$BPI_JVM_CLASS_COUNT` | 若沒提供，則動態計算 `$JAVA_HOME` 下的 class 數量 |
| JVM class 數量調整 | `--jvm-class-adj` | `$BPL_JVM_CLASS_ADJUSTMENT` | 無，可接受數字 (`1000`) 或百分比 (`150%`) |
| JVM CA 目錄 | `--jvm-cacerts` | `$BPI_JVM_CACERTS` | 若沒提供，則試著使用 `$JAVA_HOME/lib/security/cacerts` |
| Java 啟動參數 | `--jvm-options` | `$JAVA_OPTS` | |
| 是否啟用 [JDWP](https://docs.oracle.com/javase/8/docs/technotes/guides/troubleshoot/introclientissues005.html) | `--enable-jdwp` | `$BPL_DEBUG_ENABLED` | `true` |
| 是否啟用 [NMT](https://docs.oracle.com/javase/8/docs/technotes/guides/troubleshoot/tooldescr007.html) | `--enable-nmt` | `$BPL_JAVA_NMT_ENABLED` | `false` |
| 是否啟用 [JFR](https://docs.oracle.com/javacomponents/jmc-5-4/jfr-runtime-guide/about.htm) | `--enable-jfr` | `$BPL_JFR_ENABLED` | `false` |
| 是否啟用 [JMX](https://www.oracle.com/java/technologies/javase/javamanagement.html) | `--enable-jmx` | `$BPL_JMX_ENABLED` | `false` |

執行以下指令以查看完整的 args 參數說明:

```sh
memory-calculator -h
```

## Entrypoint

[`entrypoint.sh`](./entrypoint.sh) 是一個專為使用 [Jib](https://github.com/GoogleContainerTools/jib) 打包的 Image 而設計的進入點，它在執行時會根據前面提到的 [計算演算法](#calculation-algorithm) 來計算出建議的記憶體配置，然後啟動 Java 應用程式。

使用 `entrypoint.sh` 的步驟如下：

1. 將 `entrypoint.sh` 放入 Jib 所使用的 Base Image 中，例如放在 `/tmp` 資料夾下。
2. 在 Jib 的配置中，將 entrypoint 設定為 `/tmp/entrypoint.sh`。

在 Jib 中若 [自定義了 entrypoint](https://github.com/GoogleContainerTools/jib/tree/master/jib-maven-plugin#custom-container-entrypoint)，`<jvmFlags>` 參數將無法被直接引用。因此，`entrypoint.sh` 還整合了另一套 Jib 擴充工具: [jib-jvm-flags-extension-maven](https://github.com/softleader/jib-jvm-flags-extension-maven)。藉由這個 Jib Extension，我們就可以繼續使用 `<jvmFlags>`。

### 支援的參數

`entrypoint.sh` 支援以下作業系統環境變數：

- `MEM_CALC_HOME`：指定可執行檔所在目錄 (預設 `/usr/local/bin`)
- `MEM_CALC_ENABLED`：是否啟用記憶體計算工具，可設定 `true`/`false` (預設 `true`)
- `MEM_CALC_DEBUG`：是否啟用除錯模式，可設定 `true`/`false`，除錯模式中會列印出計算過程中除錯訊息，及啟動 Java 應用程式的完整指令 (預設 `false`)
- `MEM_CALC_ENABLE_PREVIEW`：是否啟用預覽功能，可設定 `true`/`false` (預設 `false`)

> [!NOTE]
> 預覽功能可能包含仍在測試或不穩定的功能, 其 API 也可能隨時調整

## 開發前準備

- Golang: v1.24+
- Jib ContainerTool
- Linux 基礎映像

## Install Script

使用 [`install.sh`](./installsh) 腳本可以在系統（Linux 或 macOS）上輕鬆安裝 `memory-calculator`。該腳本會自動檢測操作的系統和架構，以下載正確的二進位檔案。

執行前請先確認環境要求:

- 此腳本需要將二進位檔案移動到 `/usr/local/bin` 中，請確保執行時擁有讀寫權限
- 若有指定 `--entrypoint` 目錄，請確保執行時擁有讀寫權限
- 此腳本依賴套件：`curl`, `unzip`, `uname`，請確保執行前系統已安裝

### 基本用法

執行以下命令以下載並安裝最新版本：

```sh
curl -sL https://raw.githubusercontent.com/softleader/memory-calculator/main/install.sh | sh
```

### 支援的參數

安裝腳本支援了多個參數以自定義其行為：

| Flag | 說明 | 預設值 |
|---|---|---|
| `--version=<tag>` | 指定要安裝的版本，例如 `1.2.3` | `latest` |
| `--entrypoint=<path>` | 安裝後將 `entrypoint.sh` 複製到指定目錄 | 未設置，即不複製 |
| `--os=<os>` | 覆寫檢測到的操作系統，支援: `linux`, `darwin` | 自動檢測 |
| `--arch=<arch>` | 覆寫檢測到的 CPU 架構，支援: `amd64`, `arm64` | 自動檢測 |

### 範例

**安裝特定版本:**

```sh
curl -sL https://raw.githubusercontent.com/softleader/memory-calculator/main/install.sh | sh -s -- --version=1.2.3
```

**安裝最新版本並複製 `entrypoint.sh`:**

這對於設置容器映像非常有用。

```sh
curl -sL https://raw.githubusercontent.com/softleader/memory-calculator/main/install.sh | sh -s -- --entrypoint=/path/to/my-folder
```

**強制安裝特定平台（進階用法）:**

```sh
curl -sL https://raw.githubusercontent.com/softleader/memory-calculator/main/install.sh | sh -s -- --os=linux --arch=arm64
```

**在 Dockerfile 中安裝:**

參考 [`Dockerfile`](./Dockerfile)

