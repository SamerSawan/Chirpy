-- +goose Up
ALTER TABLE users ADD COLUMN password TEXT NOT NULL DEFAULT 'unset';