#pragma once

#include <string>
#include <vector>
#include <ctime>

struct Quote {
    std::string id;
    std::string text;
    std::string author;
    int likesCount;
    std::string createdAt;
    std::string updatedAt;
    
    struct Response {
        std::string id;
        std::string text;
        std::string author;
        int likesCount;
        bool isLiked;
        std::string createdAt;
        std::string updatedAt;
    };
    
    Response toResponse(bool isLiked) const;
};

struct PaginatedQuotesResponse {
    std::vector<Quote::Response> quotes;
    int64_t total;
    int page;
    int pageSize;
    int totalPages;
};

inline int calculateTotalPages(int64_t total, int pageSize) {
    return (total + pageSize - 1) / pageSize;
}

