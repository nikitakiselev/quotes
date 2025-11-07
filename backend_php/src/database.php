<?php

namespace App;

use PDO;
use PDOException;

class Database
{
    private PDO $pdo;

    public function __construct(Config $config)
    {
        $dsn = sprintf(
            "pgsql:host=%s;port=%s;dbname=%s",
            $config->dbHost,
            $config->dbPort,
            $config->dbName
        );

        $options = [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
            PDO::ATTR_EMULATE_PREPARES => false,  // Используем нативные prepared statements
            PDO::ATTR_PERSISTENT => true,  // Persistent connections для переиспользования
            PDO::ATTR_TIMEOUT => 5,
            // Оптимизация для PostgreSQL
            PDO::ATTR_STRINGIFY_FETCHES => false,
        ];

        try {
            $this->pdo = new PDO(
                $dsn,
                $config->dbUser,
                $config->dbPassword,
                $options
            );
            // Проверяем соединение
            $this->pdo->exec("SELECT 1");
        } catch (PDOException $e) {
            throw new \RuntimeException("Failed to connect to database: " . $e->getMessage());
        }
    }

    public function getPDO(): PDO
    {
        return $this->pdo;
    }
}

