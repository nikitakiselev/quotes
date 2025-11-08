#pragma once

#include "models.h"
#include "database.h"
#include <memory>
#include <vector>
#include <string>
#include <unordered_map>

class QuoteRepository {
public:
    explicit QuoteRepository(Database* db) : db_(db) {}
    
    Quote getRandom();
    std::pair<std::vector<Quote>, int64_t> getAll(int page, int pageSize, const std::string& search = "");
    Quote getById(const std::string& id);
    void create(const Quote& quote);
    void update(const std::string& id, const Quote& quote);
    void deleteQuote(const std::string& id);
    void like(const std::string& id, const std::string& userIP, const std::string& userAgent);
    bool isLiked(const std::string& quoteId, const std::string& userIP);
    std::unordered_map<std::string, bool> areLiked(const std::vector<std::string>& quoteIds, const std::string& userIP);
    Quote getTopWeekly();
    Quote getTopAllTime();
    void resetLikes();
    
private:
    Database* db_;
    Quote rowToQuote(PGresult* res, int row);
};

