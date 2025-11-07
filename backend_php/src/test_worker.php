<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Worker;
use Spiral\RoadRunner\Http\PSR7Worker;
use Nyholm\Psr7\Factory\Psr17Factory;

$worker = Worker::create();
$psr17Factory = new Psr17Factory();
$psr7Worker = new PSR7Worker($worker, $psr17Factory, $psr17Factory, $psr17Factory);

while ($request = $psr7Worker->waitRequest()) {
    try {
        $path = $request->getUri()->getPath();
        
        if ($path === '/health') {
            $response = $psr17Factory->createResponse(200)
                ->withHeader('Content-Type', 'application/json')
                ->withHeader('X-Backend', 'php')
                ->withBody($psr17Factory->createStream(json_encode(['status' => 'ok', 'backend' => 'php'])));
        } else {
            $response = $psr17Factory->createResponse(200)
                ->withHeader('Content-Type', 'application/json')
                ->withHeader('X-Backend', 'php')
                ->withBody($psr17Factory->createStream(json_encode(['message' => 'PHP backend working', 'path' => $path])));
        }
        
        $psr7Worker->respond($response);
    } catch (\Throwable $e) {
        $psr7Worker->getWorker()->error((string)$e);
        $response = $psr17Factory->createResponse(500)
            ->withHeader('Content-Type', 'application/json')
            ->withBody($psr17Factory->createStream(json_encode(['error' => 'Internal server error'])));
        $psr7Worker->respond($response);
    }
}
