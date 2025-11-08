#include "repository.h"
#include <libpq-fe.h>
#include <stdexcept>
#include <sstream>
#include <algorithm>

Quote QuoteRepository::rowToQuote(PGresult* res, int row) {
    Quote quote;
    quote.id = PQgetvalue(res, row, 0);
    quote.text = PQgetvalue(res, row, 1);
    quote.author = PQgetvalue(res, row, 2);
    quote.likesCount = std::stoi(PQgetvalue(res, row, 3));
    quote.createdAt = PQgetvalue(res, row, 4);
    quote.updatedAt = PQgetvalue(res, row, 5);
    return quote;
}

Quote QuoteRepository::getRandom() {
    const char* query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes ORDER BY RANDOM() LIMIT 1";
    PGresult* res = PQexec(db_->getConnection(), query);
    
    if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) == 0) {
        PQclear(res);
        throw std::runtime_error("no quotes found");
    }
    
    Quote quote = rowToQuote(res, 0);
    PQclear(res);
    return quote;
}

std::pair<std::vector<Quote>, int64_t> QuoteRepository::getAll(int page, int pageSize, const std::string& search) {
    int offset = (page - 1) * pageSize;
    std::vector<Quote> quotes;
    int64_t total = 0;
    
    PGconn* conn = db_->getConnection();
    
    // Count total
    if (!search.empty()) {
        std::string searchPattern = "%" + search + "%";
        const char* countQuery = "SELECT COUNT(*) FROM quotes WHERE text ILIKE $1 OR author ILIKE $1";
        const char* paramValues[1] = {searchPattern.c_str()};
        int paramLengths[1] = {static_cast<int>(search.length())};
        int paramFormats[1] = {0};
        
        int paramLengths[1] = {static_cast<int>(searchPattern.length())};
        PGresult* countRes = PQexecParams(conn, countQuery, 1, nullptr, paramValues, paramLengths, paramFormats, 0);
        if (PQresultStatus(countRes) == PGRES_TUPLES_OK) {
            total = std::stoll(PQgetvalue(countRes, 0, 0));
        }
        PQclear(countRes);
        
        // Get quotes with search
        const char* query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes WHERE text ILIKE $1 OR author ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3";
        const char* values[3] = {searchPattern.c_str(), std::to_string(pageSize).c_str(), std::to_string(offset).c_str()};
        int lengths[3] = {static_cast<int>(searchPattern.length()), -1, -1};
        int formats[3] = {0, 0, 0};
        
        PGresult* res = PQexecParams(conn, query, 3, nullptr, values, lengths, formats, 0);
        if (PQresultStatus(res) == PGRES_TUPLES_OK) {
            int ntuples = PQntuples(res);
            quotes.reserve(ntuples);
            for (int i = 0; i < ntuples; ++i) {
                quotes.push_back(rowToQuote(res, i));
            }
        }
        PQclear(res);
    } else {
        const char* countQuery = "SELECT COUNT(*) FROM quotes";
        PGresult* countRes = PQexec(conn, countQuery);
        if (PQresultStatus(countRes) == PGRES_TUPLES_OK) {
            total = std::stoll(PQgetvalue(countRes, 0, 0));
        }
        PQclear(countRes);
        
        // Get quotes
        std::string query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes ORDER BY created_at DESC LIMIT " + 
                           std::to_string(pageSize) + " OFFSET " + std::to_string(offset);
        PGresult* res = PQexec(conn, query.c_str());
        if (PQresultStatus(res) == PGRES_TUPLES_OK) {
            int ntuples = PQntuples(res);
            quotes.reserve(ntuples);
            for (int i = 0; i < ntuples; ++i) {
                quotes.push_back(rowToQuote(res, i));
            }
        }
        PQclear(res);
    }
    
    return {quotes, total};
}

Quote QuoteRepository::getById(const std::string& id) {
    const char* query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes WHERE id = $1";
    const char* values[1] = {id.c_str()};
    int lengths[1] = {static_cast<int>(id.length())};
    int formats[1] = {0};
    
    PGresult* res = PQexecParams(db_->getConnection(), query, 1, nullptr, values, lengths, formats, 0);
    
    if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) == 0) {
        PQclear(res);
        throw std::runtime_error("quote not found");
    }
    
    Quote quote = rowToQuote(res, 0);
    PQclear(res);
    return quote;
}

void QuoteRepository::create(const Quote& quote) {
    const char* query = "INSERT INTO quotes (id, text, author, likes_count, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)";
    const char* values[6] = {
        quote.id.c_str(),
        quote.text.c_str(),
        quote.author.c_str(),
        std::to_string(quote.likesCount).c_str(),
        quote.createdAt.c_str(),
        quote.updatedAt.c_str()
    };
    int lengths[6] = {
        static_cast<int>(quote.id.length()),
        static_cast<int>(quote.text.length()),
        static_cast<int>(quote.author.length()),
        -1, -1, -1
    };
    int formats[6] = {0, 0, 0, 0, 0, 0};
    
    PGresult* res = PQexecParams(db_->getConnection(), query, 6, nullptr, values, lengths, formats, 0);
    if (PQresultStatus(res) != PGRES_COMMAND_OK) {
        std::string error = PQerrorMessage(db_->getConnection());
        PQclear(res);
        throw std::runtime_error("failed to create quote: " + error);
    }
    PQclear(res);
}

void QuoteRepository::update(const std::string& id, const Quote& quote) {
    const char* query = "UPDATE quotes SET text = $1, author = $2, updated_at = NOW() WHERE id = $3";
    const char* values[3] = {quote.text.c_str(), quote.author.c_str(), id.c_str()};
    int lengths[3] = {
        static_cast<int>(quote.text.length()),
        static_cast<int>(quote.author.length()),
        static_cast<int>(id.length())
    };
    int formats[3] = {0, 0, 0};
    
    PGresult* res = PQexecParams(db_->getConnection(), query, 3, nullptr, values, lengths, formats, 0);
    if (PQresultStatus(res) != PGRES_COMMAND_OK || PQcmdTuples(res)[0] == '0') {
        PQclear(res);
        throw std::runtime_error("quote not found");
    }
    PQclear(res);
}

void QuoteRepository::deleteQuote(const std::string& id) {
    const char* query = "DELETE FROM quotes WHERE id = $1";
    const char* values[1] = {id.c_str()};
    int lengths[1] = {static_cast<int>(id.length())};
    int formats[1] = {0};
    
    PGresult* res = PQexecParams(db_->getConnection(), query, 1, nullptr, values, lengths, formats, 0);
    if (PQresultStatus(res) != PGRES_COMMAND_OK || PQcmdTuples(res)[0] == '0') {
        PQclear(res);
        throw std::runtime_error("quote not found");
    }
    PQclear(res);
}

void QuoteRepository::like(const std::string& id, const std::string& userIP, const std::string& userAgent) {
    PGconn* conn = db_->getConnection();
    
    // Begin transaction
    PGresult* beginRes = PQexec(conn, "BEGIN");
    PQclear(beginRes);
    
    try {
        // Check if already liked
        const char* checkQuery = "SELECT id FROM likes WHERE quote_id = $1 AND user_ip = $2 FOR UPDATE";
        const char* checkValues[2] = {id.c_str(), userIP.c_str()};
        int checkLengths[2] = {static_cast<int>(id.length()), static_cast<int>(userIP.length())};
        int checkFormats[2] = {0, 0};
        
        PGresult* checkRes = PQexecParams(conn, checkQuery, 2, nullptr, checkValues, checkLengths, checkFormats, 0);
        if (PQntuples(checkRes) > 0) {
            PQclear(checkRes);
            PQexec(conn, "ROLLBACK");
            throw std::runtime_error("you have already liked this quote");
        }
        PQclear(checkRes);
        
        // Update likes count
        const char* updateQuery = "UPDATE quotes SET likes_count = likes_count + 1, updated_at = NOW() WHERE id = $1";
        const char* updateValues[1] = {id.c_str()};
        int updateLengths[1] = {static_cast<int>(id.length())};
        int updateFormats[1] = {0};
        
        PGresult* updateRes = PQexecParams(conn, updateQuery, 1, nullptr, updateValues, updateLengths, updateFormats, 0);
        if (PQresultStatus(updateRes) != PGRES_COMMAND_OK || PQcmdTuples(updateRes)[0] == '0') {
            PQclear(updateRes);
            PQexec(conn, "ROLLBACK");
            throw std::runtime_error("quote not found");
        }
        PQclear(updateRes);
        
        // Insert like (using gen_random_uuid() for UUID generation)
        const char* insertQuery = "INSERT INTO likes (id, quote_id, user_ip, user_agent, created_at) VALUES (gen_random_uuid()::text, $1, $2, $3, NOW()) ON CONFLICT (quote_id, user_ip) DO NOTHING";
        const char* insertValues[3] = {id.c_str(), userIP.c_str(), userAgent.c_str()};
        int insertLengths[3] = {
            static_cast<int>(id.length()),
            static_cast<int>(userIP.length()),
            static_cast<int>(userAgent.length())
        };
        int insertFormats[3] = {0, 0, 0};
        
        PGresult* insertRes = PQexecParams(conn, insertQuery, 3, nullptr, insertValues, insertLengths, insertFormats, 0);
        PQclear(insertRes);
        
        // Commit
        PGresult* commitRes = PQexec(conn, "COMMIT");
        PQclear(commitRes);
    } catch (...) {
        PQexec(conn, "ROLLBACK");
        throw;
    }
}

bool QuoteRepository::isLiked(const std::string& quoteId, const std::string& userIP) {
    const char* query = "SELECT COUNT(*) FROM likes WHERE quote_id = $1 AND user_ip = $2";
    const char* values[2] = {quoteId.c_str(), userIP.c_str()};
    int lengths[2] = {
        static_cast<int>(quoteId.length()),
        static_cast<int>(userIP.length())
    };
    int formats[2] = {0, 0};
    
    PGresult* res = PQexecParams(db_->getConnection(), query, 2, nullptr, values, lengths, formats, 0);
    bool result = false;
    if (PQresultStatus(res) == PGRES_TUPLES_OK) {
        result = std::stoi(PQgetvalue(res, 0, 0)) > 0;
    }
    PQclear(res);
    return result;
}

std::unordered_map<std::string, bool> QuoteRepository::areLiked(const std::vector<std::string>& quoteIds, const std::string& userIP) {
    std::unordered_map<std::string, bool> result;
    
    if (quoteIds.empty()) {
        return result;
    }
    
    // Build query with IN clause (simpler and works reliably)
    std::string query = "SELECT quote_id FROM likes WHERE quote_id IN (";
    std::vector<const char*> values(quoteIds.size() + 1);
    std::vector<int> lengths(quoteIds.size() + 1);
    std::vector<int> formats(quoteIds.size() + 1, 0);
    
    for (size_t i = 0; i < quoteIds.size(); ++i) {
        if (i > 0) query += ", ";
        query += "$" + std::to_string(i + 1);
        values[i] = quoteIds[i].c_str();
        lengths[i] = static_cast<int>(quoteIds[i].length());
        result[quoteIds[i]] = false; // Initialize to false
    }
    query += ") AND user_ip = $" + std::to_string(quoteIds.size() + 1);
    values[quoteIds.size()] = userIP.c_str();
    lengths[quoteIds.size()] = static_cast<int>(userIP.length());
    
    PGresult* res = PQexecParams(db_->getConnection(), query.c_str(), 
                                 static_cast<int>(values.size()), nullptr, 
                                 values.data(), lengths.data(), formats.data(), 0);
    
    if (PQresultStatus(res) == PGRES_TUPLES_OK) {
        int ntuples = PQntuples(res);
        for (int i = 0; i < ntuples; ++i) {
            std::string quoteId = PQgetvalue(res, i, 0);
            result[quoteId] = true;
        }
    }
    PQclear(res);
    
    return result;
}

Quote QuoteRepository::getTopWeekly() {
    const char* query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes WHERE created_at >= NOW() - INTERVAL '7 days' ORDER BY likes_count DESC, created_at DESC LIMIT 1";
    PGresult* res = PQexec(db_->getConnection(), query);
    
    if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) == 0) {
        PQclear(res);
        throw std::runtime_error("no quotes found for the last week");
    }
    
    Quote quote = rowToQuote(res, 0);
    PQclear(res);
    return quote;
}

Quote QuoteRepository::getTopAllTime() {
    const char* query = "SELECT id, text, author, likes_count, created_at, updated_at FROM quotes ORDER BY likes_count DESC, created_at DESC LIMIT 1";
    PGresult* res = PQexec(db_->getConnection(), query);
    
    if (PQresultStatus(res) != PGRES_TUPLES_OK || PQntuples(res) == 0) {
        PQclear(res);
        throw std::runtime_error("no quotes found");
    }
    
    Quote quote = rowToQuote(res, 0);
    PQclear(res);
    return quote;
}

void QuoteRepository::resetLikes() {
    PGconn* conn = db_->getConnection();
    
    PGresult* beginRes = PQexec(conn, "BEGIN");
    PQclear(beginRes);
    
    try {
        PGresult* updateRes = PQexec(conn, "UPDATE quotes SET likes_count = 0, updated_at = NOW()");
        if (PQresultStatus(updateRes) != PGRES_COMMAND_OK) {
            PQclear(updateRes);
            PQexec(conn, "ROLLBACK");
            throw std::runtime_error("failed to reset likes");
        }
        PQclear(updateRes);
        
        PGresult* deleteRes = PQexec(conn, "DELETE FROM likes");
        PQclear(deleteRes);
        
        PGresult* commitRes = PQexec(conn, "COMMIT");
        PQclear(commitRes);
    } catch (...) {
        PQexec(conn, "ROLLBACK");
        throw;
    }
}

