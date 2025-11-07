<?php

namespace App;

use PDO;

class QuoteRepository
{
    private PDO $pdo;

    public function __construct(PDO $pdo)
    {
        $this->pdo = $pdo;
    }

    public function getRandom(): ?Quote
    {
        // Используем prepare вместо query для лучшей производительности
        $stmt = $this->pdo->prepare(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             ORDER BY RANDOM() 
             LIMIT 1"
        );
        $stmt->execute();
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ? Quote::fromArray($row) : null;
    }

    public function getAll(int $page, int $pageSize, ?string $search = null): array
    {
        $offset = ($page - 1) * $pageSize;

        // Подсчет общего количества
        if ($search) {
            $countStmt = $this->pdo->prepare(
                "SELECT COUNT(*) FROM quotes WHERE text ILIKE ? OR author ILIKE ?"
            );
            $searchPattern = "%{$search}%";
            $countStmt->execute([$searchPattern, $searchPattern]);
            $total = (int)$countStmt->fetchColumn();
        } else {
            // Используем prepare вместо query
            $countStmt = $this->pdo->prepare("SELECT COUNT(*) FROM quotes");
            $countStmt->execute();
            $total = (int)$countStmt->fetchColumn();
        }

        // Получение цитат
        if ($search) {
            $stmt = $this->pdo->prepare(
                "SELECT id, text, author, likes_count, created_at, updated_at 
                 FROM quotes 
                 WHERE text ILIKE ? OR author ILIKE ? 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?"
            );
            $searchPattern = "%{$search}%";
            $stmt->execute([$searchPattern, $searchPattern, $pageSize, $offset]);
        } else {
            $stmt = $this->pdo->prepare(
                "SELECT id, text, author, likes_count, created_at, updated_at 
                 FROM quotes 
                 ORDER BY created_at DESC 
                 LIMIT ? OFFSET ?"
            );
            $stmt->execute([$pageSize, $offset]);
        }

        $quotes = [];
        while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $quotes[] = Quote::fromArray($row);
        }

        return [$quotes, $total];
    }

    public function getById(string $id): ?Quote
    {
        $stmt = $this->pdo->prepare(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             WHERE id = ?"
        );
        $stmt->execute([$id]);
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ? Quote::fromArray($row) : null;
    }

    public function create(Quote $quote): void
    {
        $stmt = $this->pdo->prepare(
            "INSERT INTO quotes (id, text, author, likes_count, created_at, updated_at)
             VALUES (?, ?, ?, ?, ?, ?)"
        );
        $stmt->execute([
            $quote->id,
            $quote->text,
            $quote->author,
            $quote->likesCount,
            $quote->createdAt,
            $quote->updatedAt
        ]);
    }

    public function update(string $id, Quote $quote): bool
    {
        $stmt = $this->pdo->prepare(
            "UPDATE quotes 
             SET text = ?, author = ?, updated_at = ?
             WHERE id = ?"
        );
        $stmt->execute([
            $quote->text,
            $quote->author,
            $quote->updatedAt,
            $id
        ]);
        return $stmt->rowCount() > 0;
    }

    public function delete(string $id): bool
    {
        $stmt = $this->pdo->prepare("DELETE FROM quotes WHERE id = ?");
        $stmt->execute([$id]);
        return $stmt->rowCount() > 0;
    }

    public function like(string $id, string $userIp, ?string $userAgent): void
    {
        $this->pdo->beginTransaction();
        try {
            // Проверяем, не лайкал ли уже
            $checkStmt = $this->pdo->prepare(
                "SELECT id FROM likes WHERE quote_id = ? AND user_ip = ? FOR UPDATE"
            );
            $checkStmt->execute([$id, $userIp]);
            if ($checkStmt->fetch()) {
                throw new \RuntimeException("you have already liked this quote");
            }

            // Увеличиваем счетчик
            $updateStmt = $this->pdo->prepare(
                "UPDATE quotes 
                 SET likes_count = likes_count + 1, updated_at = NOW()
                 WHERE id = ?"
            );
            $updateStmt->execute([$id]);
            if ($updateStmt->rowCount() === 0) {
                throw new \RuntimeException("quote not found");
            }

            // Сохраняем лайк (генерируем UUID)
            $likeId = sprintf(
                '%04x%04x-%04x-%04x-%04x-%04x%04x%04x',
                mt_rand(0, 0xffff), mt_rand(0, 0xffff),
                mt_rand(0, 0xffff),
                mt_rand(0, 0x0fff) | 0x4000,
                mt_rand(0, 0x3fff) | 0x8000,
                mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff)
            );
            $insertStmt = $this->pdo->prepare(
                "INSERT INTO likes (id, quote_id, user_ip, user_agent, created_at)
                 VALUES (?, ?, ?, ?, NOW())
                 ON CONFLICT (quote_id, user_ip) DO NOTHING"
            );
            $insertStmt->execute([$likeId, $id, $userIp, $userAgent]);

            $this->pdo->commit();
        } catch (\Exception $e) {
            $this->pdo->rollBack();
            throw $e;
        }
    }

    public function isLiked(string $quoteId, string $userIp): bool
    {
        $stmt = $this->pdo->prepare(
            "SELECT COUNT(*) FROM likes WHERE quote_id = ? AND user_ip = ?"
        );
        $stmt->execute([$quoteId, $userIp]);
        return (int)$stmt->fetchColumn() > 0;
    }

    /**
     * Batch проверка лайков для оптимизации (устранение N+1 проблемы)
     * @param string[] $quoteIds
     * @param string $userIp
     * @return array<string, bool> Map quote_id => is_liked
     */
    public function areLiked(array $quoteIds, string $userIp): array
    {
        if (empty($quoteIds)) {
            return [];
        }

        // Используем ANY для эффективного batch запроса
        $placeholders = implode(',', array_fill(0, count($quoteIds), '?'));
        $stmt = $this->pdo->prepare(
            "SELECT quote_id FROM likes WHERE quote_id IN ($placeholders) AND user_ip = ?"
        );
        
        $params = array_merge($quoteIds, [$userIp]);
        $stmt->execute($params);

        $likedIds = [];
        while ($row = $stmt->fetch(PDO::FETCH_ASSOC)) {
            $likedIds[$row['quote_id']] = true;
        }

        // Создаем полный массив с false для нелайкнутых
        $result = [];
        foreach ($quoteIds as $quoteId) {
            $result[$quoteId] = isset($likedIds[$quoteId]);
        }

        return $result;
    }

    public function getTopWeekly(): ?Quote
    {
        // Используем prepare вместо query
        $stmt = $this->pdo->prepare(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             WHERE created_at >= NOW() - INTERVAL '7 days'
             ORDER BY likes_count DESC, created_at DESC
             LIMIT 1"
        );
        $stmt->execute();
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ? Quote::fromArray($row) : null;
    }

    public function getTopAllTime(): ?Quote
    {
        // Используем prepare вместо query
        $stmt = $this->pdo->prepare(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             ORDER BY likes_count DESC, created_at DESC
             LIMIT 1"
        );
        $stmt->execute();
        $row = $stmt->fetch(PDO::FETCH_ASSOC);
        return $row ? Quote::fromArray($row) : null;
    }

    public function resetLikes(): void
    {
        $this->pdo->beginTransaction();
        try {
            $this->pdo->exec("UPDATE quotes SET likes_count = 0, updated_at = NOW()");
            $this->pdo->exec("DELETE FROM likes");
            $this->pdo->commit();
        } catch (\Exception $e) {
            $this->pdo->rollBack();
            throw $e;
        }
    }

    public static function calculateTotalPages(int $total, int $pageSize): int
    {
        return (int)ceil($total / $pageSize);
    }
}

