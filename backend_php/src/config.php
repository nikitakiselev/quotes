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
        // Используем getenv() для получения переменных окружения из Docker
        $this->dbHost = getenv('DB_HOST') ?: ($_ENV['DB_HOST'] ?? 'localhost');
        $this->dbPort = getenv('DB_PORT') ?: ($_ENV['DB_PORT'] ?? '5432');
        $this->dbUser = getenv('DB_USER') ?: ($_ENV['DB_USER'] ?? 'quotes_user');
        $this->dbPassword = getenv('DB_PASSWORD') ?: ($_ENV['DB_PASSWORD'] ?? 'quotes_password');
        $this->dbName = getenv('DB_NAME') ?: ($_ENV['DB_NAME'] ?? 'quotes_db');
        $this->dbSSLMode = getenv('DB_SSLMODE') ?: ($_ENV['DB_SSLMODE'] ?? 'disable');
        $this->apiPort = getenv('API_PORT') ?: ($_ENV['API_PORT'] ?? '8082');
        $this->corsOrigin = getenv('CORS_ORIGIN') ?: ($_ENV['CORS_ORIGIN'] ?? '*');
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

