<?php
// Минимальный worker для тестирования
error_reporting(E_ALL);
ini_set('display_errors', 0);
ini_set('log_errors', 1);

fwrite(STDERR, "=== Worker starting ===\n");

try {
    fwrite(STDERR, "Loading autoload...\n");
    require __DIR__ . '/../vendor/autoload.php';
    fwrite(STDERR, "Autoload loaded\n");
    
    fwrite(STDERR, "Loading classes...\n");
    use Spiral\RoadRunner\Http\PSR7Worker;
    use Spiral\RoadRunner\Worker;
    use Nyholm\Psr7\Response;
    fwrite(STDERR, "Classes loaded\n");
    
    fwrite(STDERR, "Creating Worker...\n");
    $rrWorker = Worker::create();
    fwrite(STDERR, "Worker created\n");
    
    fwrite(STDERR, "Creating PSR7Worker...\n");
    $worker = PSR7Worker::create($rrWorker);
    fwrite(STDERR, "PSR7Worker created, entering loop...\n");
    
    while (true) {
        fwrite(STDERR, "Waiting for request...\n");
        $request = $worker->waitRequest();
        fwrite(STDERR, "Request received: " . $request->getUri()->getPath() . "\n");
        
        $path = $request->getUri()->getPath();
        $response = new Response(200);
        
        if ($path === '/health') {
            $response->getBody()->write(json_encode(['status' => 'ok', 'backend' => 'php']));
        } else {
            $response->getBody()->write(json_encode(['message' => 'PHP backend working', 'path' => $path]));
        }
        
        $response = $response->withHeader('Content-Type', 'application/json')
            ->withHeader('X-Backend', 'php');
        
        fwrite(STDERR, "Sending response...\n");
        $worker->respond($response);
        fwrite(STDERR, "Response sent\n");
    }
} catch (\Throwable $e) {
    fwrite(STDERR, "ERROR: " . $e->getMessage() . "\n");
    fwrite(STDERR, $e->getTraceAsString() . "\n");
    exit(1);
}

