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

# Build Drogon from source
WORKDIR /tmp
RUN git clone --depth 1 --branch v1.9.0 https://github.com/drogonframework/drogon.git && \
    cd drogon && \
    git submodule update --init --recursive && \
    mkdir build && cd build && \
    cmake .. \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_INSTALL_PREFIX=/usr/local && \
    cmake --build . -j$(nproc) && \
    cmake --install . && \
    # Verify installation
    ls -la /usr/local/lib/cmake/drogon/ || echo "Drogon cmake files not found" && \
    find /usr/local -name "*drogon*" -type f | head -10

# Copy source code
WORKDIR /app
COPY backend_cpp /app

# Build application
RUN mkdir -p build && cd build && \
    cmake .. \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_CXX_FLAGS="-O3 -march=native -flto" \
        -DCMAKE_PREFIX_PATH=/usr/local \
        -Ddrogon_DIR=/usr/local/lib/cmake/drogon && \
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

# Copy Drogon libraries
COPY --from=builder /usr/local/lib/libdrogon.so* /usr/local/lib/
COPY --from=builder /usr/local/lib/libtrantor.so* /usr/local/lib/
RUN ldconfig

# Copy binary
COPY --from=builder /app/install/bin/quotes_backend_cpp /app/main

WORKDIR /app

# Expose port
EXPOSE 8083

# Run
CMD ["/app/main"]
