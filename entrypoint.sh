#!/usr/bin/env sh
set -e

DEFAULT_BIN_PATH="/usr/local/bin"
TMP_ENV="/tmp/.env"
JIB_CLASSPATH_FILE="/app/jib-classpath-file"
JIB_MAIN_CLASS_FILE="/app/jib-main-class-file"
JIB_JVM_FLAGS_FILE="/app/jib-jvm-flags-file"

read_jvm_flags() {
  if [ -f "$JIB_JVM_FLAGS_FILE" ]; then
    cat "$JIB_JVM_FLAGS_FILE"
  else
    echo ""
  fi
}

execute_memory_calculator() {
  if [ "${MEM_CALC_ENABLED}" = "false" ]; then
    echo "Memory calculator is disabled"
    return
  fi

  bin_path=$1
  debug=$2
  enablePreview=$3
  "$bin_path/memory-calculator" -o="$TMP_ENV" -v="$debug" --enable-preview="$enablePreview" || {
    echo "Memory calculator failed, version: $("$bin_path/memory-calculator" --version)";
    exit 1;
  }
  if [ -f "$TMP_ENV" ]; then
    . "$TMP_ENV"
  fi
}

DEBUG="${MEM_CALC_DEBUG:-false}"
ENABLE_PREVIEW="${MEM_CALC_ENABLE_PREVIEW:-false}"
BIN="${MEM_CALC_HOME:-$DEFAULT_BIN_PATH}"
JVM_FLAGS=$(read_jvm_flags)

execute_memory_calculator "$BIN" "$DEBUG" "$ENABLE_PREVIEW"

if [ "$DEBUG" = true ]; then set -x; fi
exec java $JVM_FLAGS -cp $(cat "$JIB_CLASSPATH_FILE") $(cat "$JIB_MAIN_CLASS_FILE") "$@"
