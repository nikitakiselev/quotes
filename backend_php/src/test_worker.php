<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;
use Nyholm\Psr7\Response;

$worker = PSR7Worker::create(Worker::create());

while ($request = $worker->waitRequest()) {
    $path = $request->getUri()->getPath();
    
    $response = new Response(200);
    
    if ($path === '/health') {
        $response->getBody()->write(json_encode(['status' => 'ok', 'backend' => 'php']));
    } else {
        $response->getBody()->write(json_encode(['message' => 'PHP backend working', 'path' => $path]));
    }
    
    $response = $response->withHeader('Content-Type', 'application/json')
        ->withHeader('X-Backend', 'php');
    
    $worker->respond($response);
}

