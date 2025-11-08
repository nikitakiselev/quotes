#pragma once

#include <string>

class Config {
public:
    std::string dbHost;
    std::string dbPort;
    std::string dbUser;
    std::string dbPassword;
    std::string dbName;
    std::string dbSSLMode;
    int apiPort;
    std::string corsOrigin;

    static Config load();
    std::string getDatabaseDSN() const;
};

