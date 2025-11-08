#include "config.h"
#include <cstdlib>
#include <cstring>

Config Config::load() {
    Config cfg;
    
    const char* dbHost = std::getenv("DB_HOST");
    cfg.dbHost = dbHost ? dbHost : "localhost";
    
    const char* dbPort = std::getenv("DB_PORT");
    cfg.dbPort = dbPort ? dbPort : "5432";
    
    const char* dbUser = std::getenv("DB_USER");
    cfg.dbUser = dbUser ? dbUser : "quotes_user";
    
    const char* dbPassword = std::getenv("DB_PASSWORD");
    cfg.dbPassword = dbPassword ? dbPassword : "quotes_password";
    
    const char* dbName = std::getenv("DB_NAME");
    cfg.dbName = dbName ? dbName : "quotes_db";
    
    const char* dbSSLMode = std::getenv("DB_SSLMODE");
    cfg.dbSSLMode = dbSSLMode ? dbSSLMode : "disable";
    
    const char* apiPort = std::getenv("API_PORT");
    cfg.apiPort = apiPort ? std::stoi(apiPort) : 8083;
    
    const char* corsOrigin = std::getenv("CORS_ORIGIN");
    cfg.corsOrigin = corsOrigin ? corsOrigin : "*";
    
    return cfg;
}

std::string Config::getDatabaseDSN() const {
    return "host=" + dbHost + 
           " port=" + dbPort + 
           " dbname=" + dbName + 
           " user=" + dbUser + 
           " password=" + dbPassword + 
           " sslmode=" + dbSSLMode;
}

