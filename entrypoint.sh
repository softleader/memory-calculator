#!/usr/bin/env bash
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
    return
  fi

  local bin_path=$1
  local debug=$2
  $bin_path/memory-calculator -o="$TMP_ENV" -v="$debug" || { echo "Memory calculator failed"; exit 1; }
  if [ -f "$TMP_ENV" ]; then
    source "$TMP_ENV"
  fi
}

execute_java_app() {
  local jvm_flags=$1
  local args=$2
  [ "$DEBUG" = true ] && set -x
  exec java $jvm_flags -cp $(cat "$JIB_CLASSPATH_FILE") $(cat "$JIB_MAIN_CLASS_FILE") $args
  [ "$DEBUG" = true ] && set +x
}

DEBUG="${DEBUG:-false}"
BIN="${MEM_CALC_HOME:-$DEFAULT_BIN_PATH}"
JVM_FLAGS=$(read_jvm_flags)
ARGS="$@"

execute_memory_calculator "$BIN" "$DEBUG"
execute_java_app "$JVM_FLAGS" "$ARGS"
