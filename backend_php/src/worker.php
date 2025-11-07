<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;
use Nyholm\Psr7\Response;
use App\Config;
use App\Database;
use App\QuoteRepository;
use App\Handlers;
use App\Router;

// Инициализация приложения
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
            $response = new Response(500);
            $response->getBody()->write(json_encode(['error' => 'Internal server error']));
            $response = $response->withHeader('Content-Type', 'application/json')
                ->withHeader('X-Backend', 'php');
            $worker->respond($response);
        }
    }
} catch (\Throwable $e) {
    fwrite(STDERR, "Fatal error: " . $e->getMessage() . "\n");
    fwrite(STDERR, $e->getTraceAsString() . "\n");
    exit(1);
}

