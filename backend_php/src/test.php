<?php

require __DIR__ . '/../vendor/autoload.php';

use Spiral\RoadRunner\Http\PSR7Worker;
use Spiral\RoadRunner\Worker;
use Nyholm\Psr7\Response;

$worker = PSR7Worker::create(Worker::create());

while ($request = $worker->waitRequest()) {
    $response = new Response(200);
    $response->getBody()->write(json_encode(['status' => 'ok', 'message' => 'Test response']));
    $response = $response->withHeader('Content-Type', 'application/json');
    $worker->respond($response);
}

