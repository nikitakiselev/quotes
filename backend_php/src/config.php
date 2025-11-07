<?php

namespace App;

class Config
{
    public string $dbHost;
    public string $dbPort;
    public string $dbUser;
    public string $dbPassword;
    public string $dbName;
    public string $dbSSLMode;
    public string $apiPort;
    public string $corsOrigin;

    public function __construct()
    {
        $this->dbHost = $_ENV['DB_HOST'] ?? 'localhost';
        $this->dbPort = $_ENV['DB_PORT'] ?? '5432';
        $this->dbUser = $_ENV['DB_USER'] ?? 'quotes_user';
        $this->dbPassword = $_ENV['DB_PASSWORD'] ?? 'quotes_password';
        $this->dbName = $_ENV['DB_NAME'] ?? 'quotes_db';
        $this->dbSSLMode = $_ENV['DB_SSLMODE'] ?? 'disable';
        $this->apiPort = $_ENV['API_PORT'] ?? '8082';
        $this->corsOrigin = $_ENV['CORS_ORIGIN'] ?? '*';
    }

    public function getDatabaseDSN(): string
    {
        return sprintf(
            "pgsql:host=%s;port=%s;dbname=%s;user=%s;password=%s",
            $this->dbHost,
            $this->dbPort,
            $this->dbName,
            $this->dbUser,
            $this->dbPassword
        );
    }
}

