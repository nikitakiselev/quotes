<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;
use App\Config;
use App\Database;
use App\QuoteRepository;
use App\Handlers;
use App\Router;

// Загрузка переменных окружения
if (file_exists(__DIR__ . '/../.env')) {
    $lines = file(__DIR__ . '/../.env', FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    foreach ($lines as $line) {
        if (strpos(trim($line), '#') === 0) {
            continue;
        }
        if (strpos($line, '=') === false) {
            continue;
        }
        list($name, $value) = explode('=', $line, 2);
        $_ENV[trim($name)] = trim($value);
    }
}

// Инициализация
try {
    $config = new Config();
    $database = new Database($config);
    $repository = new QuoteRepository($database->getPDO());
    $handlers = new Handlers($repository);
    $router = new Router($handlers);
    
    // Создание RoadRunner worker
    $worker = PSR7Worker::create(Worker::create());
    
    // Обработка запросов
    while ($request = $worker->waitRequest()) {
        try {
            $response = $router->handle($request);
            $worker->respond($response);
        } catch (\Throwable $e) {
            $response = new \Nyholm\Psr7\Response(500);
            $response->getBody()->write(json_encode(['error' => 'Internal server error', 'message' => $e->getMessage()]));
            $worker->respond($response);
        }
    }
} catch (\Throwable $e) {
    // Выводим ошибку в stderr для RoadRunner
    fwrite(STDERR, "Fatal error during initialization: " . $e->getMessage() . "\n");
    fwrite(STDERR, $e->getTraceAsString() . "\n");
    exit(1);
}

