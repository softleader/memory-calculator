#!/usr/bin/env bash

# 1. 先從 $1 參數取值
BIN=${1:-}

# 2. 若無指定, 判斷 MEMORY_CALCULATOR_HOME
if [ -z "$BIN" ]; then
  BIN=$MEMORY_CALCULATOR_HOME
fi

# 3. 如果都没有，使用預設 /usr/local/bin
BIN=${BIN:-/usr/local/bin}

TMP_ENV="/tmp/.env"

$BIN/memory-calculator -o $TMP_ENV

# 檢查 /app/jib-jvm-flags-file 是否存在
if [ -f /app/jib-jvm-flags-file ]; then
    JVM_FLAGS=$(cat /app/jib-jvm-flags-file)
else
    JVM_FLAGS=""
fi

if [ -f "$TMP_ENV" ]; then
  source $TMP_ENV
  exec java $JVM_FLAGS -cp $( cat /app/jib-classpath-file ) $( cat /app/jib-main-class-file )
else
  exec java $JVM_FLAGS $JAVA_OPTS -cp $( cat /app/jib-classpath-file ) $( cat /app/jib-main-class-file )
fi
