#include "models.h"

Quote::Response Quote::toResponse(bool isLiked) const {
    return {
        id,
        text,
        author,
        likesCount,
        isLiked,
        createdAt,
        updatedAt
    };
}

