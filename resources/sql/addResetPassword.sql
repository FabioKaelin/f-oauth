START TRANSACTION;

USE `oauth`;

CREATE TABLE
    `reset_password` (
        `id` varchar(64) NOT NULL,
        `secret` varchar(64) DEFAULT NULL,
        `useridfs` varchar(512) NOT NULL,
    );


ALTER TABLE `reset_password` ADD PRIMARY KEY (`id`),
ADD KEY `fk_resetpassword_user` (`useridfs`);

ALTER TABLE `reset_password` ADD CONSTRAINT `fk_resetpassword_user` FOREIGN KEY (`useridfs`) REFERENCES `users` (`id`);

COMMIT;