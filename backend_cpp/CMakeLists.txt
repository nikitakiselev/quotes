cmake_minimum_required(VERSION 3.16)
project(quotes_backend_cpp)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find dependencies
find_package(drogon REQUIRED)
find_package(PostgreSQL REQUIRED)

# Include directories
include_directories(${CMAKE_SOURCE_DIR}/include)

# Source files
set(SOURCES
    src/main.cpp
    src/config.cpp
    src/database.cpp
    src/repository.cpp
    src/handlers.cpp
    src/models.cpp
)

# Find UUID library
find_library(UUID_LIBRARY uuid)
if(UUID_LIBRARY)
    target_link_libraries(${PROJECT_NAME} ${UUID_LIBRARY})
endif()

# Create executable
add_executable(${PROJECT_NAME} ${SOURCES})

# Link libraries
target_link_libraries(${PROJECT_NAME}
    drogon::drogon
    ${PostgreSQL_LIBRARIES}
)

# Compiler flags for optimization
target_compile_options(${PROJECT_NAME} PRIVATE
    -O3
    -march=native
    -flto
)

