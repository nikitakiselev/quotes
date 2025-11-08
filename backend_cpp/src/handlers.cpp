#include "handlers.h"
#include <drogon/HttpResponse.h>
#include <drogon/HttpTypes.h>
#include <drogon/utils/Utilities.h>
#include <json/json.h>
#include <uuid/uuid.h>
#include <ctime>
#include <iomanip>
#include <sstream>
#include <random>

std::string QuoteHandlers::getClientIP(const drogon::HttpRequestPtr& req) {
    auto forwarded = req->getHeader("X-Forwarded-For");
    if (!forwarded.empty()) {
        size_t pos = forwarded.find(',');
        return pos != std::string::npos ? forwarded.substr(0, pos) : forwarded;
    }
    
    auto realIP = req->getHeader("X-Real-IP");
    if (!realIP.empty()) {
        return realIP;
    }
    
    return "127.0.0.1";
}

void QuoteHandlers::getRandom(const drogon::HttpRequestPtr& req,
                              std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        Quote quote = repo_->getRandom();
        std::string userIP = getClientIP(req);
        bool isLiked = repo_->isLiked(quote.id, userIP);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = isLiked;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::getAll(const drogon::HttpRequestPtr& req,
                           std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        auto& params = req->getParameters();
        int page = 1;
        int pageSize = 10;
        std::string search;
        
        auto pageIt = params.find("page");
        if (pageIt != params.end()) {
            page = std::max(1, std::stoi(pageIt->second));
        }
        auto pageSizeIt = params.find("page_size");
        if (pageSizeIt != params.end()) {
            pageSize = std::max(1, std::min(100, std::stoi(pageSizeIt->second)));
        }
        auto searchIt = params.find("search");
        if (searchIt != params.end()) {
            search = searchIt->second;
        }
        
        auto [quotes, total] = repo_->getAll(page, pageSize, search);
        std::string userIP = getClientIP(req);
        
        // Batch check likes
        std::vector<std::string> quoteIds;
        quoteIds.reserve(quotes.size());
        for (const auto& quote : quotes) {
            quoteIds.push_back(quote.id);
        }
        auto likedMap = repo_->areLiked(quoteIds, userIP);
        
        Json::Value json;
        Json::Value quotesArray(Json::arrayValue);
        for (const auto& quote : quotes) {
            Json::Value quoteJson;
            quoteJson["id"] = quote.id;
            quoteJson["text"] = quote.text;
            quoteJson["author"] = quote.author;
            quoteJson["likes_count"] = quote.likesCount;
            quoteJson["is_liked"] = likedMap[quote.id];
            quoteJson["created_at"] = quote.createdAt;
            quoteJson["updated_at"] = quote.updatedAt;
            quotesArray.append(quoteJson);
        }
        
        json["quotes"] = quotesArray;
        json["total"] = static_cast<Json::Int64>(total);
        json["page"] = page;
        json["page_size"] = pageSize;
        json["total_pages"] = static_cast<int>(calculateTotalPages(total, pageSize));
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k500InternalServerError);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::getById(const drogon::HttpRequestPtr& req,
                            std::function<void(const drogon::HttpResponsePtr&)>&& callback,
                            const std::string& id) {
    try {
        Quote quote = repo_->getById(id);
        std::string userIP = getClientIP(req);
        bool isLiked = repo_->isLiked(quote.id, userIP);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = isLiked;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::create(const drogon::HttpRequestPtr& req,
                           std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        auto jsonPtr = req->getJsonObject();
        if (!jsonPtr || !jsonPtr->isMember("text") || !jsonPtr->isMember("author")) {
            Json::Value error;
            error["error"] = "text and author are required";
            auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
            resp->setStatusCode(drogon::k400BadRequest);
            resp->addHeader("X-Backend", "cpp");
            callback(resp);
            return;
        }
        
        // Generate UUID v4
        std::random_device rd;
        std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(0, 15);
        std::uniform_int_distribution<> dis2(8, 11);
        
        std::stringstream uuidStream;
        for (int i = 0; i < 32; ++i) {
            if (i == 8 || i == 12 || i == 16 || i == 20) {
                uuidStream << "-";
            }
            int v = (i == 12) ? dis2(gen) : dis(gen);
            uuidStream << std::hex << v;
        }
        std::string uuidStr = uuidStream.str();
        
        // Get current time
        auto now = std::time(nullptr);
        std::stringstream ss;
        ss << std::put_time(std::gmtime(&now), "%Y-%m-%d %H:%M:%S");
        std::string timestamp = ss.str();
        
        Quote quote;
        quote.id = uuidStr;
        quote.text = (*jsonPtr)["text"].asString();
        quote.author = (*jsonPtr)["author"].asString();
        quote.likesCount = 0;
        quote.createdAt = timestamp;
        quote.updatedAt = timestamp;
        
        repo_->create(quote);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = false;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->setStatusCode(drogon::k201Created);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k500InternalServerError);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::update(const drogon::HttpRequestPtr& req,
                           std::function<void(const drogon::HttpResponsePtr&)>&& callback,
                           const std::string& id) {
    try {
        Quote quote = repo_->getById(id);
        auto jsonPtr = req->getJsonObject();
        
        if (jsonPtr && jsonPtr->isMember("text")) {
            quote.text = (*jsonPtr)["text"].asString();
        }
        if (jsonPtr && jsonPtr->isMember("author")) {
            quote.author = (*jsonPtr)["author"].asString();
        }
        
        auto now = std::time(nullptr);
        std::stringstream ss;
        ss << std::put_time(std::gmtime(&now), "%Y-%m-%d %H:%M:%S");
        quote.updatedAt = ss.str();
        
        repo_->update(id, quote);
        
        Quote updatedQuote = repo_->getById(id);
        std::string userIP = getClientIP(req);
        bool isLiked = repo_->isLiked(updatedQuote.id, userIP);
        
        Json::Value json;
        json["id"] = updatedQuote.id;
        json["text"] = updatedQuote.text;
        json["author"] = updatedQuote.author;
        json["likes_count"] = updatedQuote.likesCount;
        json["is_liked"] = isLiked;
        json["created_at"] = updatedQuote.createdAt;
        json["updated_at"] = updatedQuote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::deleteQuote(const drogon::HttpRequestPtr& req,
                                std::function<void(const drogon::HttpResponsePtr&)>&& callback,
                                const std::string& id) {
    try {
        repo_->deleteQuote(id);
        auto resp = drogon::HttpResponse::newHttpResponse();
        resp->setStatusCode(drogon::k204NoContent);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::like(const drogon::HttpRequestPtr& req,
                         std::function<void(const drogon::HttpResponsePtr&)>&& callback,
                         const std::string& id) {
    try {
        std::string userIP = getClientIP(req);
        std::string userAgent = req->getHeader("User-Agent");
        
        repo_->like(id, userIP, userAgent);
        
        Quote quote = repo_->getById(id);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = true;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        std::string errorMsg = e.what();
        Json::Value error;
        
        if (errorMsg.find("already liked") != std::string::npos) {
            error["error"] = "Вы уже поставили лайк этой цитате";
            auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
            resp->setStatusCode(drogon::k400BadRequest);
            resp->addHeader("X-Backend", "cpp");
            callback(resp);
        } else if (errorMsg.find("not found") != std::string::npos) {
            error["error"] = "quote not found";
            auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
            resp->setStatusCode(drogon::k404NotFound);
            resp->addHeader("X-Backend", "cpp");
            callback(resp);
        } else {
            error["error"] = "internal server error";
            auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
            resp->setStatusCode(drogon::k500InternalServerError);
            resp->addHeader("X-Backend", "cpp");
            callback(resp);
        }
    }
}

void QuoteHandlers::getTopWeekly(const drogon::HttpRequestPtr& req,
                                 std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        Quote quote = repo_->getTopWeekly();
        std::string userIP = getClientIP(req);
        bool isLiked = repo_->isLiked(quote.id, userIP);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = isLiked;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::getTopAllTime(const drogon::HttpRequestPtr& req,
                                  std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        Quote quote = repo_->getTopAllTime();
        std::string userIP = getClientIP(req);
        bool isLiked = repo_->isLiked(quote.id, userIP);
        
        Json::Value json;
        json["id"] = quote.id;
        json["text"] = quote.text;
        json["author"] = quote.author;
        json["likes_count"] = quote.likesCount;
        json["is_liked"] = isLiked;
        json["created_at"] = quote.createdAt;
        json["updated_at"] = quote.updatedAt;
        
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k404NotFound);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::resetLikes(const drogon::HttpRequestPtr& req,
                               std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    try {
        repo_->resetLikes();
        Json::Value json;
        json["message"] = "Все лайки успешно сброшены";
        auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    } catch (const std::exception& e) {
        Json::Value error;
        error["error"] = e.what();
        auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
        resp->setStatusCode(drogon::k500InternalServerError);
        resp->addHeader("X-Backend", "cpp");
        callback(resp);
    }
}

void QuoteHandlers::health(const drogon::HttpRequestPtr& req,
                           std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
    Json::Value json;
    json["status"] = "ok";
    auto resp = drogon::HttpResponse::newHttpJsonResponse(json);
    resp->addHeader("X-Backend", "cpp");
    callback(resp);
}

