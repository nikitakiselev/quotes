<?php

namespace App;

use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Message\ResponseInterface;
use Nyholm\Psr7\Response;

class Router
{
    private Handlers $handlers;

    public function __construct(Handlers $handlers)
    {
        $this->handlers = $handlers;
        $this->setupCors();
    }

    private function setupCors(): void
    {
        // CORS будет обрабатываться через middleware в RoadRunner
    }

    private function addCorsHeaders(ResponseInterface $response, string $origin): ResponseInterface
    {
        return $response
            ->withHeader('Access-Control-Allow-Origin', $origin)
            ->withHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS, PATCH')
            ->withHeader('Access-Control-Allow-Headers', 'Origin, Content-Type, Accept, Authorization, X-Requested-With')
            ->withHeader('Access-Control-Expose-Headers', 'Content-Length, Content-Type, X-Backend')
            ->withHeader('Access-Control-Allow-Credentials', 'true');
    }

    public function handle(ServerRequestInterface $request): ResponseInterface
    {
        $method = $request->getMethod();
        $path = $request->getUri()->getPath();

        // Обработка OPTIONS для CORS
        if ($method === 'OPTIONS') {
            $origin = $request->getHeader('Origin')[0] ?? '*';
            return $this->addCorsHeaders(new Response(204), $origin);
        }

        // Health check
        if ($path === '/health' && $method === 'GET') {
            $response = $this->handlers->health($request);
            return $this->addCorsHeaders($response, '*');
        }

        // API routes
        if (strpos($path, '/api/quotes') === 0) {
            $response = $this->handleQuotes($request, $method, $path);
            $origin = $request->getHeader('Origin')[0] ?? '*';
            return $this->addCorsHeaders($response, $origin);
        }

        // 404
        return $this->jsonResponse(404, ['error' => 'Not found']);
    }

    private function handleQuotes(ServerRequestInterface $request, string $method, string $path): ResponseInterface
    {
        // Извлекаем ID из пути, если есть
        $pathParts = explode('/', trim($path, '/'));
        $id = null;
        $isLikeRoute = false;
        
        // Проверяем паттерн /api/quotes/:id/like
        if (count($pathParts) === 5 && $pathParts[2] === 'quotes' && $pathParts[4] === 'like') {
            $id = $pathParts[3];
            $isLikeRoute = true;
        } elseif (count($pathParts) >= 4 && $pathParts[2] === 'quotes') {
            // Проверяем, что это не специальные роуты
            if ($pathParts[3] !== 'random' && $pathParts[3] !== 'top' && $pathParts[3] !== 'likes') {
                $id = $pathParts[3];
            }
        }

        // Специфичные роуты (должны быть раньше параметризованных)
        if ($path === '/api/quotes/random' && $method === 'GET') {
            return $this->handlers->getRandom($request);
        }

        if ($path === '/api/quotes/top/weekly' && $method === 'GET') {
            return $this->handlers->getTopWeekly($request);
        }

        if ($path === '/api/quotes/top/alltime' && $method === 'GET') {
            return $this->handlers->getTopAllTime($request);
        }

        if ($path === '/api/quotes/likes/reset' && $method === 'DELETE') {
            return $this->handlers->resetLikes($request);
        }

        // Роуты с ID
        if ($id) {
            if ($isLikeRoute && $method === 'PUT') {
                return $this->handlers->like($request, $id);
            }

            if (!$isLikeRoute && $method === 'GET') {
                return $this->handlers->getById($request, $id);
            }

            if (!$isLikeRoute && $method === 'PUT') {
                return $this->handlers->update($request, $id);
            }

            if (!$isLikeRoute && $method === 'DELETE') {
                return $this->handlers->delete($request, $id);
            }
        }

        // Общие роуты
        if ($path === '/api/quotes' && $method === 'GET') {
            return $this->handlers->getAll($request);
        }

        if ($path === '/api/quotes' && $method === 'POST') {
            return $this->handlers->create($request);
        }

        return $this->jsonResponse(404, ['error' => 'Not found']);
    }

    private function jsonResponse(int $statusCode, array $data): ResponseInterface
    {
        $response = new Response($statusCode);
        $response->getBody()->write(json_encode($data, JSON_UNESCAPED_UNICODE));
        return $response->withHeader('Content-Type', 'application/json')
            ->withHeader('X-Backend', 'php');
    }
}

