<?php

namespace App;

class Quote
{
    public string $id;
    public string $text;
    public string $author;
    public int $likesCount;
    public string $createdAt;
    public string $updatedAt;

    public function __construct(
        string $id,
        string $text,
        string $author,
        int $likesCount,
        string $createdAt,
        string $updatedAt
    ) {
        $this->id = $id;
        $this->text = $text;
        $this->author = $author;
        $this->likesCount = $likesCount;
        $this->createdAt = $createdAt;
        $this->updatedAt = $updatedAt;
    }

    public function toResponse(bool $isLiked): array
    {
        return [
            'id' => $this->id,
            'text' => $this->text,
            'author' => $this->author,
            'likes_count' => $this->likesCount,
            'is_liked' => $isLiked,
            'created_at' => $this->createdAt,
            'updated_at' => $this->updatedAt,
        ];
    }

    public static function fromArray(array $data): self
    {
        return new self(
            $data['id'],
            $data['text'],
            $data['author'],
            (int)$data['likes_count'],
            $data['created_at'],
            $data['updated_at']
        );
    }
}

