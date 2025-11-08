#include "database.h"
#include <stdexcept>

Database::Database(const std::string& dsn) {
    conn_ = PQconnectdb(dsn.c_str());
    
    if (PQstatus(conn_) != CONNECTION_OK) {
        std::string error = PQerrorMessage(conn_);
        PQfinish(conn_);
        throw std::runtime_error("Failed to connect to database: " + error);
    }
    
    // Set connection to use prepared statements
    PQexec(conn_, "SET statement_timeout = 5000");
}

Database::~Database() {
    if (conn_) {
        PQfinish(conn_);
    }
}

bool Database::isConnected() const {
    return conn_ && PQstatus(conn_) == CONNECTION_OK;
}

