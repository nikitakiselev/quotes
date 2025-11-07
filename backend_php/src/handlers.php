<?php

namespace App;

use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Message\ResponseInterface;
use Nyholm\Psr7\Response;

class Handlers
{
    private QuoteRepository $repository;

    public function __construct(QuoteRepository $repository)
    {
        $this->repository = $repository;
    }

    private function getClientIp(ServerRequestInterface $request): string
    {
        $headers = $request->getHeader('X-Forwarded-For');
        if (!empty($headers)) {
            $ips = explode(',', $headers[0]);
            return trim($ips[0]);
        }

        $headers = $request->getHeader('X-Real-IP');
        if (!empty($headers)) {
            return $headers[0];
        }

        return '127.0.0.1';
    }

    private function jsonResponse(int $statusCode, array $data): ResponseInterface
    {
        $response = new Response($statusCode);
        $response->getBody()->write(json_encode($data, JSON_UNESCAPED_UNICODE));
        return $response->withHeader('Content-Type', 'application/json')
            ->withHeader('X-Backend', 'php');
    }

    public function getRandom(ServerRequestInterface $request): ResponseInterface
    {
        $quote = $this->repository->getRandom();
        if (!$quote) {
            return $this->jsonResponse(404, ['error' => 'no quotes found']);
        }

        $userIp = $this->getClientIp($request);
        $isLiked = $this->repository->isLiked($quote->id, $userIp);

        return $this->jsonResponse(200, $quote->toResponse($isLiked));
    }

    public function getAll(ServerRequestInterface $request): ResponseInterface
    {
        $queryParams = $request->getQueryParams();
        $page = max(1, (int)($queryParams['page'] ?? 1));
        $pageSize = max(1, min(100, (int)($queryParams['page_size'] ?? 10)));
        $search = $queryParams['search'] ?? null;

        [$quotes, $total] = $this->repository->getAll($page, $pageSize, $search);

        // Оптимизация: batch проверка лайков вместо N+1 запросов
        $userIp = $this->getClientIp($request);
        $quoteIds = array_map(fn($quote) => $quote->id, $quotes);
        $likedMap = $this->repository->areLiked($quoteIds, $userIp);
        
        $responses = [];
        foreach ($quotes as $quote) {
            $isLiked = $likedMap[$quote->id] ?? false;
            $responses[] = $quote->toResponse($isLiked);
        }

        $totalPages = QuoteRepository::calculateTotalPages($total, $pageSize);

        return $this->jsonResponse(200, [
            'quotes' => $responses,
            'total' => $total,
            'page' => $page,
            'page_size' => $pageSize,
            'total_pages' => $totalPages,
        ]);
    }

    public function getById(ServerRequestInterface $request, string $id): ResponseInterface
    {
        $quote = $this->repository->getById($id);
        if (!$quote) {
            return $this->jsonResponse(404, ['error' => 'quote not found']);
        }

        $userIp = $this->getClientIp($request);
        $isLiked = $this->repository->isLiked($quote->id, $userIp);

        return $this->jsonResponse(200, $quote->toResponse($isLiked));
    }

    public function create(ServerRequestInterface $request): ResponseInterface
    {
        $body = json_decode($request->getBody()->getContents(), true);
        if (!isset($body['text']) || !isset($body['author'])) {
            return $this->jsonResponse(400, ['error' => 'text and author are required']);
        }

        // Генерируем UUID в формате, совместимом с другими бэкендами
        $uuid = sprintf(
            '%04x%04x-%04x-%04x-%04x-%04x%04x%04x',
            mt_rand(0, 0xffff), mt_rand(0, 0xffff),
            mt_rand(0, 0xffff),
            mt_rand(0, 0x0fff) | 0x4000,
            mt_rand(0, 0x3fff) | 0x8000,
            mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
        );

        $quote = new Quote(
            $uuid,
            $body['text'],
            $body['author'],
            0,
            date('c'),
            date('c')
        );

        $this->repository->create($quote);

        return $this->jsonResponse(201, $quote->toResponse(false));
    }

    public function update(ServerRequestInterface $request, string $id): ResponseInterface
    {
        $quote = $this->repository->getById($id);
        if (!$quote) {
            return $this->jsonResponse(404, ['error' => 'quote not found']);
        }

        $body = json_decode($request->getBody()->getContents(), true);
        if (isset($body['text'])) {
            $quote->text = $body['text'];
        }
        if (isset($body['author'])) {
            $quote->author = $body['author'];
        }
        $quote->updatedAt = date('c');

        if (!$this->repository->update($id, $quote)) {
            return $this->jsonResponse(500, ['error' => 'internal server error']);
        }

        $updatedQuote = $this->repository->getById($id);
        $userIp = $this->getClientIp($request);
        $isLiked = $this->repository->isLiked($updatedQuote->id, $userIp);

        return $this->jsonResponse(200, $updatedQuote->toResponse($isLiked));
    }

    public function delete(ServerRequestInterface $request, string $id): ResponseInterface
    {
        if (!$this->repository->delete($id)) {
            return $this->jsonResponse(404, ['error' => 'quote not found']);
        }

        return new Response(204);
    }

    public function like(ServerRequestInterface $request, string $id): ResponseInterface
    {
        $userIp = $this->getClientIp($request);
        $userAgent = $request->getHeader('User-Agent')[0] ?? null;

        try {
            $this->repository->like($id, $userIp, $userAgent);
        } catch (\RuntimeException $e) {
            if (strpos($e->getMessage(), 'already liked') !== false) {
                return $this->jsonResponse(400, ['error' => 'Вы уже поставили лайк этой цитате']);
            }
            if (strpos($e->getMessage(), 'not found') !== false) {
                return $this->jsonResponse(404, ['error' => 'quote not found']);
            }
            return $this->jsonResponse(500, ['error' => 'internal server error']);
        }

        $quote = $this->repository->getById($id);
        return $this->jsonResponse(200, $quote->toResponse(true));
    }

    public function getTopWeekly(ServerRequestInterface $request): ResponseInterface
    {
        $quote = $this->repository->getTopWeekly();
        if (!$quote) {
            return $this->jsonResponse(404, ['error' => 'no quotes found for the last week']);
        }

        $userIp = $this->getClientIp($request);
        $isLiked = $this->repository->isLiked($quote->id, $userIp);

        return $this->jsonResponse(200, $quote->toResponse($isLiked));
    }

    public function getTopAllTime(ServerRequestInterface $request): ResponseInterface
    {
        $quote = $this->repository->getTopAllTime();
        if (!$quote) {
            return $this->jsonResponse(404, ['error' => 'no quotes found']);
        }

        $userIp = $this->getClientIp($request);
        $isLiked = $this->repository->isLiked($quote->id, $userIp);

        return $this->jsonResponse(200, $quote->toResponse($isLiked));
    }

    public function resetLikes(ServerRequestInterface $request): ResponseInterface
    {
        $this->repository->resetLikes();
        return $this->jsonResponse(200, ['message' => 'Все лайки успешно сброшены']);
    }

    public function health(ServerRequestInterface $request): ResponseInterface
    {
        return $this->jsonResponse(200, ['status' => 'ok']);
    }
}

