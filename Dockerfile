FROM eclipse-temurin:17-jre-alpine

# 記憶體計算工具
RUN apk add --no-cache curl unzip && \
    curl -sL https://raw.githubusercontent.com/softleader/memory-calculator/refs/heads/main/install.sh | sh -s -- --entrypoint=/

ENTRYPOINT ["/entrypoint.sh"]
