cmake_minimum_required(VERSION 3.16)
project(quotes_backend_cpp)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Include FetchContent for downloading Drogon
include(FetchContent)

# Download and build Drogon
FetchContent_Declare(
    drogon
    GIT_REPOSITORY https://github.com/drogonframework/drogon.git
    GIT_TAG v1.9.0
)

# Set Drogon options before fetching
set(BUILD_TESTING OFF CACHE BOOL "" FORCE)
set(BUILD_EXAMPLES OFF CACHE BOOL "" FORCE)

FetchContent_MakeAvailable(drogon)

# Find PostgreSQL
find_package(PostgreSQL REQUIRED)

# Include directories
include_directories(${CMAKE_SOURCE_DIR}/include)
include_directories(${PostgreSQL_INCLUDE_DIRS})

# Source files
set(SOURCES
    src/main.cpp
    src/config.cpp
    src/database.cpp
    src/repository.cpp
    src/handlers.cpp
    src/models.cpp
)

# Create executable
add_executable(${PROJECT_NAME} ${SOURCES})

# Link libraries
target_link_libraries(${PROJECT_NAME}
    PRIVATE
    drogon
    ${PostgreSQL_LIBRARIES}
)

# Compiler flags for optimization
target_compile_options(${PROJECT_NAME} PRIVATE
    -O3
    -march=native
    -flto
)
