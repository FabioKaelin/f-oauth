-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: mariadb
-- Generation Time: Aug 10, 2023 at 06:46 PM
-- Server version: 11.0.2-MariaDB-1:11.0.2+maria~ubu2204
-- PHP Version: 8.1.20
SET
    SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

START TRANSACTION;

SET
    time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;

/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;

/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;

/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `oauth`
--
-- --------------------------------------------------------
--
-- Table structure for table `users`
--

CREATE DATABASE `oauth`;

USE `oauth`;

CREATE TABLE
    `users` (
        `id` varchar(64) NOT NULL DEFAULT uuid (),
        `name` varchar(64) DEFAULT NULL,
        `password` varchar(512) NOT NULL,
        `email` varchar(64) DEFAULT NULL,
        `role` varchar(64) DEFAULT NULL,
        `provider` varchar(64) DEFAULT NULL,
        `photo` varchar(512) DEFAULT NULL,
        `verified` varchar(64) NOT NULL,
        `created_at` date NOT NULL DEFAULT current_timestamp(),
        `updated_at` date NOT NULL DEFAULT current_timestamp()
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

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;

/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;