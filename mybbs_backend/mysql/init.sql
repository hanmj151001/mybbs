USE bbsgo_db;
SET NAMES utf8mb4;
-- 初始化用户表
CREATE TABLE `t_user` (
                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
                          `username` varchar(32) DEFAULT NULL,
                          `email` varchar(128) DEFAULT NULL,

                          `nickname` varchar(16) DEFAULT NULL,
                          `uuid`     varchar(32) DEFAULT NULL,
                          `background_image` text,
                          `password` varchar(512) DEFAULT NULL,
                          `home_page` varchar(1024) DEFAULT NULL,
                          `description` text,
                          `score` bigint(20) NOT NULL,
                          `status` bigint(20) NOT NULL,
                          `topic_count` bigint(20) NOT NULL,
                          `comment_count` bigint(20) NOT NULL,
                          `roles` text,
                          `forbidden_end_time` bigint(20) NOT NULL DEFAULT '0',
                          `create_time` bigint(20) DEFAULT NULL,
                          `update_time` bigint(20) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `username` (`username`),
                          UNIQUE KEY `email` (`email`),
                          KEY `idx_user_score` (`score`),
                          KEY `idx_user_status` (`status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;