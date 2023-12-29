# memory-calculator

memory-calculator 是基於 [paketo-buildpacks/libjvm](https://github.com/paketo-buildpacks/libjvm/) 所撰寫幫助 JVM 在
Runtime 時計算記憶體設定的工具

## PreRequirement:

- Golang: v1.20+
- Jib ContainerTool
- Linux base image

## Configuration

| Environment Variable          | Description                                                                                                |
|-------------------------------|------------------------------------------------------------------------------------------------------------|
| `$BPL_JVM_HEAD_ROOM`          | Configure the percentage of headroom the memory calculator will allocated.  Defaults to `0`.               |
| `$BPL_JVM_LOADED_CLASS_COUNT` | Configure the number of classes that will be loaded at runtime.  Defaults to 35% of the number of classes. |
| `$BPL_JVM_THREAD_COUNT`       | Configure the number of user threads at runtime.  Defaults to `250`.                                       |
| `$JAVA_HOME`                  | Configure the JRE location                                                                                 |
| `$JAVA_OPTS`                  | Configure the JAVA launch flags                                                                            |
| `$JAVA_TOOL_OPTIONS`          | Configure the JVM launch flags                                                                             |

## Usage

Put `memory-calculator` and `entrypoint.sh` in `/tmp` and custom entrypoint to `/tmp/entrypoint.sh`
Don't use `jvmFlags` but point to `JAVA_OPTS` environment variable