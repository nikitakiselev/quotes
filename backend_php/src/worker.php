<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Worker;
use Spiral\RoadRunner\Http\PSR7Worker;
use Nyholm\Psr7\Factory\Psr17Factory;
use App\Config;
use App\Database;
use App\QuoteRepository;
use App\Handlers;
use App\Router;

// Инициализация приложения
try {
    fwrite(STDERR, "Initializing application...\n");
    $config = new Config();
    fwrite(STDERR, "Config loaded\n");
    
    $database = new Database($config);
    fwrite(STDERR, "Database connected\n");
    
    $repository = new QuoteRepository($database->getPDO());
    $handlers = new Handlers($repository);
    $router = new Router($handlers);
    fwrite(STDERR, "Router initialized\n");
    
    // Создание RoadRunner worker
    $worker = Worker::create();
    $psr17Factory = new Psr17Factory();
    $psr7Worker = new PSR7Worker($worker, $psr17Factory, $psr17Factory, $psr17Factory);
    fwrite(STDERR, "Worker created, entering loop...\n");
    
    // Обработка запросов
    while ($request = $psr7Worker->waitRequest()) {
        try {
            $response = $router->handle($request);
            $psr7Worker->respond($response);
        } catch (\Throwable $e) {
            fwrite(STDERR, "Error handling request: " . $e->getMessage() . "\n");
            $psr7Worker->getWorker()->error((string)$e);
            $response = $psr17Factory->createResponse(500)
                ->withHeader('Content-Type', 'application/json')
                ->withHeader('X-Backend', 'php')
                ->withBody($psr17Factory->createStream(json_encode(['error' => 'Internal server error'])));
            $psr7Worker->respond($response);
        }
    }
} catch (\Throwable $e) {
    fwrite(STDERR, "Fatal error: " . $e->getMessage() . "\n");
    fwrite(STDERR, $e->getTraceAsString() . "\n");
    exit(1);
}
