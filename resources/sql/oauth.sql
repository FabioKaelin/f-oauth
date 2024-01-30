-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: mariadb
-- Generation Time: Aug 10, 2023 at 06:46 PM
-- Server version: 11.0.2-MariaDB-1:11.0.2+maria~ubu2204
-- PHP Version: 8.1.20
-- SET
--     SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

START TRANSACTION;

-- SET
--     time_zone = "+00:00";


--
-- Database: `oauth`
--
-- --------------------------------------------------------
--
-- Table structure for table `users`
--
CREATE DATABASE IF NOT EXISTS `oauth`;


USE `oauth`;

CREATE TABLE
    `users` (
        `id` varchar(64) NOT NULL,
        `name` varchar(64) DEFAULT NULL,
        `password` varchar(512) NOT NULL,
        `email` varchar(64) DEFAULT NULL,
        `role` varchar(64) DEFAULT NULL,
        `provider` varchar(64) DEFAULT NULL,
        `photo` varchar(512) DEFAULT NULL,
        `verified` tinyint (1) DEFAULT 0,
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

--
-- Dumping data for table `users`
--
--
-- Indexes for dumped tables
--
--
-- Indexes for table `users`
--
ALTER TABLE `users` ADD PRIMARY KEY (`id`);

COMMIT;
