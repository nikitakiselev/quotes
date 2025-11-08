# Build stage
FROM ubuntu:22.04 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libpq-dev \
    libjsoncpp-dev \
    uuid-dev \
    pkg-config \
    git \
    libssl-dev \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy source code
WORKDIR /app
COPY backend_cpp /app

# Build application (Drogon will be downloaded and built automatically)
RUN mkdir -p build && cd build && \
    cmake .. \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_CXX_FLAGS="-O3 -march=native -flto" && \
    cmake --build . -j$(nproc) && \
    cmake --install . --prefix /app/install

# Runtime stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libpq5 \
    libjsoncpp25 \
    libuuid1 \
    libssl3 \
    zlib1g \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary (includes all dependencies statically or we copy shared libs)
COPY --from=builder /app/install/bin/quotes_backend_cpp /app/main

# Copy Drogon shared libraries if needed
COPY --from=builder /app/build/_deps/drogon-build/lib/libdrogon.so* /usr/local/lib/ 2>/dev/null || true
COPY --from=builder /app/build/_deps/drogon-build/lib/libtrantor.so* /usr/local/lib/ 2>/dev/null || true
RUN ldconfig

WORKDIR /app

# Expose port
EXPOSE 8083

# Run
CMD ["/app/main"]
