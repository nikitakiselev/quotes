<?php
// Выводим в stderr для отладки
fwrite(STDERR, "Worker starting...\n");

require __DIR__ . '/../vendor/autoload.php';

fwrite(STDERR, "Autoload loaded\n");

use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;
use Nyholm\Psr7\Response;

fwrite(STDERR, "Creating worker...\n");
$worker = PSR7Worker::create(Worker::create());
fwrite(STDERR, "Worker created, waiting for requests...\n");

while ($request = $worker->waitRequest()) {
    fwrite(STDERR, "Request received: " . $request->getUri()->getPath() . "\n");
    $path = $request->getUri()->getPath();
    
    if ($path === '/health') {
        $response = new Response(200);
        $response->getBody()->write(json_encode(['status' => 'ok', 'backend' => 'php']));
        $response = $response->withHeader('Content-Type', 'application/json')
            ->withHeader('X-Backend', 'php');
        $worker->respond($response);
        continue;
    }
    
    if (strpos($path, '/api/quotes') === 0) {
        $response = new Response(200);
        $response->getBody()->write(json_encode(['message' => 'PHP backend is working', 'path' => $path]));
        $response = $response->withHeader('Content-Type', 'application/json')
            ->withHeader('X-Backend', 'php');
        $worker->respond($response);
        continue;
    }
    
    $response = new Response(404);
    $response->getBody()->write(json_encode(['error' => 'Not found']));
    $response = $response->withHeader('Content-Type', 'application/json')
        ->withHeader('X-Backend', 'php');
    $worker->respond($response);
}

