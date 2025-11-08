#include <drogon/drogon.h>
#include "config.h"
#include "database.h"
#include "repository.h"
#include "handlers.h"
#include <iostream>
#include <memory>

int main() {
    try {
        // Load configuration
        Config cfg = Config::load();
        
        // Connect to database
        Database db(cfg.getDatabaseDSN());
        if (!db.isConnected()) {
            std::cerr << "Failed to connect to database" << std::endl;
            return 1;
        }
        
        // Initialize repository and handlers
        QuoteRepository repo(&db);
        QuoteHandlers handlers(&repo);
        
        // Setup CORS and X-Backend header
        drogon::app().registerPostHandlingAdvice([](const drogon::HttpRequestPtr& req, const drogon::HttpResponsePtr& resp) {
            std::string origin = req->getHeader("Origin");
            if (origin.empty()) {
                origin = "*";
            }
            resp->addHeader("Access-Control-Allow-Origin", origin);
            resp->addHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH");
            resp->addHeader("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With");
            resp->addHeader("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Backend");
            resp->addHeader("Access-Control-Allow-Credentials", "true");
            resp->addHeader("X-Backend", "cpp");
        });
        
        // Handle OPTIONS for CORS
        drogon::app().registerHandler("/api/*",
            [](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                if (req->getMethod() == drogon::Options) {
                    auto resp = drogon::HttpResponse::newHttpResponse();
                    resp->setStatusCode(drogon::k204NoContent);
                    resp->addHeader("X-Backend", "cpp");
                    callback(resp);
                }
            },
            {drogon::Options});
        
        // Health check
        drogon::app().registerHandler("/health",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                handlers.health(req, std::move(callback));
            },
            {drogon::Get});
        
        // API routes
        drogon::app().registerHandler("/api/quotes/random",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                handlers.getRandom(req, std::move(callback));
            },
            {drogon::Get});
        
        drogon::app().registerHandler("/api/quotes/top/weekly",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                handlers.getTopWeekly(req, std::move(callback));
            },
            {drogon::Get});
        
        drogon::app().registerHandler("/api/quotes/top/alltime",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                handlers.getTopAllTime(req, std::move(callback));
            },
            {drogon::Get});
        
        drogon::app().registerHandler("/api/quotes/likes/reset",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                handlers.resetLikes(req, std::move(callback));
            },
            {drogon::Delete});
        
        drogon::app().registerHandler("/api/quotes",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                if (req->getMethod() == drogon::Get) {
                    handlers.getAll(req, std::move(callback));
                } else if (req->getMethod() == drogon::Post) {
                    handlers.create(req, std::move(callback));
                }
            },
            {drogon::Get, drogon::Post});
        
        drogon::app().registerHandler("/api/quotes/{id}",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback, const std::string& id) {
                if (req->getMethod() == drogon::Get) {
                    handlers.getById(req, std::move(callback), id);
                } else if (req->getMethod() == drogon::Put) {
                    handlers.update(req, std::move(callback), id);
                } else if (req->getMethod() == drogon::Delete) {
                    handlers.deleteQuote(req, std::move(callback), id);
                }
            },
            {drogon::Get, drogon::Put, drogon::Delete});
        
        drogon::app().registerHandler("/api/quotes/{id}/like",
            [&handlers](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback, const std::string& id) {
                handlers.like(req, std::move(callback), id);
            },
            {drogon::Put});
        
        // 404 handler
        drogon::app().registerHandler(".*",
            [](const drogon::HttpRequestPtr& req, std::function<void(const drogon::HttpResponsePtr&)>&& callback) {
                Json::Value error;
                error["error"] = "Not found";
                auto resp = drogon::HttpResponse::newHttpJsonResponse(error);
                resp->setStatusCode(drogon::k404NotFound);
                resp->addHeader("X-Backend", "cpp");
                callback(resp);
            },
            {drogon::Get, drogon::Post, drogon::Put, drogon::Delete});
        
        // Start server
        std::string addr = "0.0.0.0";
        std::cout << "Starting C++ backend on " << addr << ":" << cfg.apiPort << std::endl;
        
        drogon::app()
            .setLogLevel(trantor::Logger::kWarn)
            .addListener(addr, static_cast<uint16_t>(cfg.apiPort))
            .setThreadNum(0)  // Auto-detect optimal thread count (CPU cores)
            .run();
            
    } catch (const std::exception& e) {
        std::cerr << "Fatal error: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
}

