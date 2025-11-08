#pragma once

#include <libpq-fe.h>
#include <memory>
#include <string>

class Database {
public:
    explicit Database(const std::string& dsn);
    ~Database();
    
    Database(const Database&) = delete;
    Database& operator=(const Database&) = delete;
    
    PGconn* getConnection() const { return conn_; }
    bool isConnected() const;
    
private:
    PGconn* conn_;
};

