-- Migration: Add expires_at, used, created_at columns to reset_password table
-- Run this against the 'oauth' database after deploying the password reset feature.

USE `oauth`;

-- Increase secret column to hold 64-char hex token (was 64, stays fine)
-- Add expiry timestamp
ALTER TABLE `reset_password`
    ADD COLUMN `expires_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN `used` TINYINT(1) NOT NULL DEFAULT 0,
    ADD COLUMN `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
