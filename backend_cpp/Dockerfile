# Build stage
FROM ubuntu:22.04 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libpq-dev \
    libdrogon-dev \
    libjsoncpp-dev \
    uuid-dev \
    pkg-config \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy source code
WORKDIR /app
COPY backend_cpp /app

# Build
RUN mkdir -p build && cd build && \
    cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_CXX_FLAGS="-O3 -march=native -flto" && \
    cmake --build . -j$(nproc) && \
    cmake --install . --prefix /app/install

# Runtime stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libpq5 \
    libdrogon1.8.0 \
    libjsoncpp25 \
    libuuid1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary
COPY --from=builder /app/install/bin/quotes_backend_cpp /app/main

WORKDIR /app

# Expose port
EXPOSE 8083

# Run
CMD ["/app/main"]

