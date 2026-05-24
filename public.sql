/*
 Navicat Premium Dump SQL

 Source Server         : local-mysql
 Source Server Type    : MySQL
 Source Server Version : 80036 (8.0.36)
 Source Host           : 127.0.0.1:3306
 Source Schema         : psychological_test

 Target Server Type    : MySQL
 Target Server Version : 80036 (8.0.36)
 File Encoding         : 65001

 Date: 12/11/2025 20:11:24
*/

-- ----------------------------
-- Database structure for psychological_test
-- ----------------------------
CREATE DATABASE IF NOT EXISTS `psychological_test` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `psychological_test`;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识，自增主键',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录名，唯一',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户密码，BCrypt加密存储',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户昵称，用于显示',
  `avatar_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户头像URL地址',
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '微信开放平台唯一标识',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号码',
  `user_type` int DEFAULT '0' COMMENT '用户类型(0-普通用户,1-管理员)',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `gender` int DEFAULT '0' COMMENT '性别(0-未知,1-男,2-女)',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `deleted_time` datetime DEFAULT NULL COMMENT '软删除标记',
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_username_idx` (`username`),
  UNIQUE KEY `users_mobile_idx` (`mobile`),
  UNIQUE KEY `users_open_id_idx` (`open_id`),
  KEY `users_password_idx` (`password`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='存储系统用户基本信息';

-- ----------------------------
-- Table structure for user_accounts (用户余额账户表)
-- ----------------------------
DROP TABLE IF EXISTS `user_accounts`;
CREATE TABLE `user_accounts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '用户ID',
    `balance` DECIMAL(12,2) NOT NULL DEFAULT '0.00' COMMENT '账户余额（元）',
    `frozen_amount` DECIMAL(12,2) NOT NULL DEFAULT '0.00' COMMENT '冻结金额',
    `total_recharge` DECIMAL(12,2) NOT NULL DEFAULT '0.00' COMMENT '累计充值金额',
    `total_consume` DECIMAL(12,2) NOT NULL DEFAULT '0.00' COMMENT '累计消费金额',
    `version` INT NOT NULL DEFAULT '0' COMMENT '乐观锁版本号',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户余额账户表';

-- ----------------------------
-- Table structure for test_categories
-- ----------------------------
DROP TABLE IF EXISTS `test_categories`;
CREATE TABLE `test_categories` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `sort` int DEFAULT '0' COMMENT '排序优先级',
  `status` int DEFAULT '1' COMMENT '状态(1-启用, 0-禁用)',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试分类表';

-- ----------------------------
-- Records of test_categories
-- ----------------------------
BEGIN;
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('健康测试', 10);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('性格测试', 20);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('情感测试', 30);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('职业测试', 40);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('趣味测试', 50);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('智商测试', 60);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('能力测试', 70);
INSERT INTO `test_categories` (`name`, `sort`) VALUES ('亲子测试', 80);
COMMIT;

-- ----------------------------
-- Table structure for test_types
-- ----------------------------
DROP TABLE IF EXISTS `test_types`;
CREATE TABLE `test_types` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '类型ID，自增主键',
  `category_id` bigint DEFAULT NULL COMMENT '所属分类ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '测试名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '测试描述',
  `test_time` int DEFAULT NULL COMMENT '预计测试时间(分钟)',
  `status` int DEFAULT '1' COMMENT '状态(0-禁用,1-启用)',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `img` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '测试图标URL',
  `price` DECIMAL(10, 2) DEFAULT 0.00 COMMENT '查看报告所需支付金额（元）',
  PRIMARY KEY (`id`),
  KEY `fk_test_types_category` (`category_id`),
  CONSTRAINT `fk_test_types_category` FOREIGN KEY (`category_id`) REFERENCES `test_categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试类型定义';

-- ----------------------------
-- Table structure for questions
-- ----------------------------
DROP TABLE IF EXISTS `questions`;
CREATE TABLE `questions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '问题ID，自增主键',
  `test_type_id` bigint DEFAULT NULL COMMENT '关联的测试类型ID',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '问题内容文本',
  `sequence` int DEFAULT NULL COMMENT '问题顺序号',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `questions_test_type_id_idx` (`test_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试题目';

-- ----------------------------
-- Table structure for options
-- ----------------------------
DROP TABLE IF EXISTS `options`;
CREATE TABLE `options` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '选项ID，自增主键',
  `question_id` bigint DEFAULT NULL COMMENT '关联的问题ID',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '选项内容文本',
  `score` int DEFAULT '0' COMMENT '选项对应的分数',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `options_question_id_idx` (`question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试题目选项';

-- ----------------------------
-- Table structure for test_records
-- ----------------------------
DROP TABLE IF EXISTS `test_records`;
CREATE TABLE `test_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID，自增主键',
  `user_id` bigint DEFAULT NULL COMMENT '关联用户ID',
  `test_type_id` int DEFAULT NULL COMMENT '测试类型ID',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '进行中' COMMENT '测试状态(进行中/已完成)',
  `start_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '测试开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '测试结束时间',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '软删除标记(0-未删除,1-已删除)',
  PRIMARY KEY (`id`),
  KEY `test_records_user_id_idx` (`user_id`),
  KEY `test_records_test_type_id_idx` (`test_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户测试记录';

-- ----------------------------
-- Table structure for user_answers
-- ----------------------------
DROP TABLE IF EXISTS `user_answers`;
CREATE TABLE `user_answers` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '答案ID，自增主键',
  `user_id` bigint DEFAULT NULL COMMENT '关联用户ID',
  `test_type_id` bigint DEFAULT NULL COMMENT '测试类型ID',
  `question_id` bigint DEFAULT NULL COMMENT '问题ID',
  `option_id` bigint DEFAULT NULL COMMENT '选择的选项ID',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `answer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户作答选择',
  `deleted_time` datetime DEFAULT NULL COMMENT '软删除标记',
  PRIMARY KEY (`id`),
  KEY `user_answers_user_id_idx` (`user_id`),
  KEY `user_answers_test_type_id_idx` (`test_type_id`),
  KEY `user_answers_question_id_idx` (`question_id`),
  KEY `user_answers_option_id_idx` (`option_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户测试答案记录';

-- ----------------------------
-- Table structure for test_results
-- ----------------------------
DROP TABLE IF EXISTS `test_results`;
CREATE TABLE `test_results` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '结果ID，自增主键',
  `test_type_id` bigint DEFAULT NULL COMMENT '关联的测试类型ID',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '结果详细描述',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `test_results_test_type_id_idx` (`test_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预定义的测试结果解释';

-- ----------------------------
-- Table structure for favorites
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '收藏ID，自增主键',
  `user_id` bigint NOT NULL COMMENT '关联用户ID',
  `test_type_id` bigint NOT NULL COMMENT '测试类型ID',
  `result_id` bigint DEFAULT NULL COMMENT '测试结果ID（可选）',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '收藏创建时间',
  `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_time` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_test` (`user_id`, `test_type_id`),
  KEY `idx_favorites_user_deleted` (`user_id`, `deleted_time`),
  KEY `idx_favorites_test_type` (`test_type_id`),
  KEY `idx_favorites_result` (`result_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收藏记录表';

-- 创建测试结果分析建议表
CREATE TABLE IF NOT EXISTS user_test_results (
    id BIGINT NOT NULL AUTO_INCREMENT COMMENT '结果ID，自增主键',
    test_type_id BIGINT NOT NULL COMMENT '测试类型ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    result_data JSON NOT NULL COMMENT '结果数据(包含总分、分析结果、建议等)',
    test_date DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '测试时间',
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_test_type (test_type_id),
    KEY idx_user_id (user_id),
    KEY idx_test_date (test_date),
    FOREIGN KEY fk_result_test_type (test_type_id) REFERENCES test_types(id),
    FOREIGN KEY fk_result_user (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户测试结果记录';

-- ----------------------------
-- Foreign Key Constraints
-- ----------------------------
ALTER TABLE `questions` ADD CONSTRAINT `fk_questions_test_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`);
ALTER TABLE `options` ADD CONSTRAINT `fk_options_question` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);
ALTER TABLE `test_records` ADD CONSTRAINT `fk_test_records_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `user_answers` ADD CONSTRAINT `fk_user_answers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `user_answers` ADD CONSTRAINT `fk_user_answers_question` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`);
ALTER TABLE `user_answers` ADD CONSTRAINT `fk_user_answers_option` FOREIGN KEY (`option_id`) REFERENCES `options` (`id`);

ALTER TABLE `test_results` ADD CONSTRAINT `fk_test_results_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`);
ALTER TABLE `favorites` ADD CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `favorites` ADD CONSTRAINT `fk_favorites_test_type` FOREIGN KEY (`test_type_id`) REFERENCES `test_types` (`id`);
ALTER TABLE `favorites` ADD CONSTRAINT `fk_favorites_result` FOREIGN KEY (`result_id`) REFERENCES `user_test_results` (`id`) ON DELETE SET NULL;

-- ----------------------------
-- Table structure for verification_codes
-- ----------------------------
DROP TABLE IF EXISTS `verification_codes`;
CREATE TABLE `verification_codes` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '验证码ID',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号码',
  `code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '验证码',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '业务类型(LOGIN/REGISTER/...)',
  `status` tinyint DEFAULT '0' COMMENT '状态(0-未使用, 1-已使用, 2-已过期)',
  `expire_at` datetime NOT NULL COMMENT '过期时间',
  `client_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端IP',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_mobile_type_status` (`mobile`, `type`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信验证码记录表';

-- ----------------------------
-- Table structure for test_orders (业务订单表)
-- ----------------------------
DROP TABLE IF EXISTS `test_orders`;
CREATE TABLE `test_orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    `order_no` VARCHAR(32) NOT NULL UNIQUE COMMENT '订单编号（业务唯一）',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `test_id` INT UNSIGNED NOT NULL COMMENT '测试类型ID',
    `test_name` VARCHAR(100) NOT NULL COMMENT '测试名称（冗余存储）',
    `original_price` DECIMAL(10,2) NOT NULL COMMENT '原价',
    `discount_amount` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
    `payable_amount` DECIMAL(10,2) NOT NULL COMMENT '应付金额（实际支付）',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭',
    `paid_at` DATETIME DEFAULT NULL COMMENT '支付完成时间',
    `paid_channel` VARCHAR(20) DEFAULT NULL COMMENT '支付渠道：alipay/wechat',
    `paid_transaction_id` VARCHAR(64) DEFAULT NULL COMMENT '第三方支付流水号',
    `result_report_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '生成的报告ID',
    `expire_at` DATETIME NOT NULL COMMENT '订单过期时间',
    `client_ip` VARCHAR(40) DEFAULT NULL COMMENT '用户下单IP',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    `is_deleted` TINYINT NOT NULL DEFAULT '0',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_user_status` (`user_id`, `status`),
    KEY `idx_status_created` (`status`, `created_at`),
    KEY `idx_expire` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='性格测试订单表';

-- ----------------------------
-- Table structure for payment_transactions (支付流水表)
-- ----------------------------
DROP TABLE IF EXISTS `payment_transactions`;
CREATE TABLE `payment_transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `transaction_no` VARCHAR(32) NOT NULL UNIQUE COMMENT '内部支付流水号（唯一）',
    `order_no` VARCHAR(32) NOT NULL COMMENT '业务订单号',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `amount` DECIMAL(10,2) NOT NULL COMMENT '支付金额',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '币种',
    `channel` VARCHAR(20) NOT NULL COMMENT '支付渠道：alipay/wechat',
    `channel_app_id` VARCHAR(64) DEFAULT NULL COMMENT '渠道应用ID',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '状态：0-待支付 1-支付中 2-支付成功 3-支付失败 4-已关闭 5-已退款',
    `channel_extra` JSON DEFAULT NULL COMMENT '渠道扩展参数',
    `alipay_trade_no` VARCHAR(64) DEFAULT NULL COMMENT '支付宝交易号',
    `wechat_prepay_id` VARCHAR(64) DEFAULT NULL COMMENT '微信预支付ID',
    `wechat_transaction_id` VARCHAR(64) DEFAULT NULL COMMENT '微信订单号',
    `prepay_at` DATETIME DEFAULT NULL COMMENT '预下单时间',
    `paid_at` DATETIME DEFAULT NULL COMMENT '支付成功时间',
    `closed_at` DATETIME DEFAULT NULL COMMENT '关闭时间',
    `expire_at` DATETIME NOT NULL COMMENT '过期时间',
    `notify_raw_data` TEXT DEFAULT NULL COMMENT '原始回调数据',
    `notified_at` DATETIME DEFAULT NULL COMMENT '收到回调时间',
    `idempotency_key` VARCHAR(64) UNIQUE COMMENT '幂等键',
    `client_ip` VARCHAR(40) DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_order_no` (`order_no`),
    KEY `idx_status_created` (`status`, `created_at`),
    KEY `idx_paid_at` (`paid_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付流水表';

-- ----------------------------
-- Table structure for payment_qrcodes (支付二维码表)
-- ----------------------------
DROP TABLE IF EXISTS `payment_qrcodes`;
CREATE TABLE `payment_qrcodes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `transaction_no` VARCHAR(32) NOT NULL UNIQUE COMMENT '关联支付流水号',
    `qr_type` TINYINT NOT NULL COMMENT '类型：1-支付宝 2-微信',
    `qr_content` TEXT NOT NULL COMMENT '二维码内容',
    `qr_image_url` VARCHAR(500) DEFAULT NULL COMMENT '二维码图片URL',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '0-有效 1-已扫描 2-已支付 3-已过期',
    `scanned_at` DATETIME DEFAULT NULL COMMENT '被扫描时间',
    `scanned_by` VARCHAR(50) DEFAULT NULL COMMENT '扫描者标识',
    `expire_at` DATETIME NOT NULL COMMENT '过期时间',
    `valid_seconds` INT UNSIGNED NOT NULL DEFAULT '300',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `idx_expire_status` (`expire_at`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付二维码表';

-- ----------------------------
-- Table structure for payment_notifications (支付回调通知表)
-- ----------------------------
DROP TABLE IF EXISTS `payment_notifications`;
CREATE TABLE `payment_notifications` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `channel` VARCHAR(20) NOT NULL COMMENT 'alipay/wechat',
    `notify_type` VARCHAR(30) NOT NULL COMMENT '类型：pay/refund/dispute',
    `raw_headers` TEXT COMMENT 'HTTP头',
    `raw_body` TEXT NOT NULL COMMENT '原始请求体',
    `parsed_data` JSON COMMENT '解析后的数据',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '0-待处理 1-处理中 2-成功 3-失败',
    `retry_count` INT UNSIGNED DEFAULT '0',
    `error_msg` VARCHAR(500) DEFAULT NULL,
    `transaction_no` VARCHAR(32) DEFAULT NULL,
    `order_no` VARCHAR(32) DEFAULT NULL,
    `received_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `processed_at` DATETIME DEFAULT NULL,
    KEY `idx_status_received` (`status`, `received_at`),
    KEY `idx_transaction` (`transaction_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付回调通知表';

-- ----------------------------
-- Table structure for idempotency_tokens (幂等令牌表)
-- ----------------------------
DROP TABLE IF EXISTS `idempotency_tokens`;
CREATE TABLE `idempotency_tokens` (
    `token` VARCHAR(64) NOT NULL PRIMARY KEY,
    `business_type` VARCHAR(20) NOT NULL COMMENT '业务类型',
    `business_id` VARCHAR(32) NOT NULL COMMENT '业务单号',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '0-处理中 1-成功 2-失败',
    `expire_at` DATETIME NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_business` (`business_type`, `business_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='幂等令牌表';

-- ----------------------------
-- Table structure for refund_records (退款记录表)
-- ----------------------------
DROP TABLE IF EXISTS `refund_records`;
CREATE TABLE `refund_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `refund_no` VARCHAR(32) NOT NULL UNIQUE COMMENT '内部退款单号',
    `transaction_no` VARCHAR(32) NOT NULL COMMENT '原支付流水号',
    `order_no` VARCHAR(32) NOT NULL COMMENT '原订单号',
    `refund_amount` DECIMAL(10,2) NOT NULL COMMENT '退款金额',
    `reason` VARCHAR(200) DEFAULT NULL COMMENT '退款原因',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '0-申请中 1-处理中 2-成功 3-失败',
    `channel_refund_id` VARCHAR(64) DEFAULT NULL COMMENT '渠道退款单号',
    `channel_response` TEXT DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `completed_at` DATETIME DEFAULT NULL,
    UNIQUE KEY `uk_refund_order` (`order_no`),
    KEY `idx_transaction` (`transaction_no`),
    KEY `idx_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款记录表';

-- ----------------------------
-- Table structure for reconciliation_records (对账记录表)
-- ----------------------------
DROP TABLE IF EXISTS `reconciliation_records`;
CREATE TABLE `reconciliation_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `reconcile_date` DATE NOT NULL COMMENT '对账日期',
    `channel` VARCHAR(20) NOT NULL COMMENT '渠道',
    `total_transactions` INT UNSIGNED NOT NULL DEFAULT '0',
    `total_amount` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
    `matched_count` INT UNSIGNED NOT NULL DEFAULT '0',
    `unmatched_count` INT UNSIGNED NOT NULL DEFAULT '0',
    `status` TINYINT NOT NULL DEFAULT '0' COMMENT '0-待对账 1-对账中 2-完成 3-异常',
    `detail_file_url` VARCHAR(500) DEFAULT NULL,
    `exception_detail` TEXT DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `completed_at` DATETIME DEFAULT NULL,
    UNIQUE KEY `uk_date_channel` (`reconcile_date`, `channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对账记录表';

-- ----------------------------
-- Table structure for reconciliation_exceptions (对账异常明细表)
-- ----------------------------
DROP TABLE IF EXISTS `reconciliation_exceptions`;
CREATE TABLE `reconciliation_exceptions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `reconcile_date` DATE NOT NULL,
    `channel` VARCHAR(20) NOT NULL,
    `exception_type` VARCHAR(30) NOT NULL COMMENT 'LOCAL_MISS/AMOUNT_MISMATCH等',
    `local_transaction_no` VARCHAR(32) DEFAULT NULL,
    `channel_trade_no` VARCHAR(64) DEFAULT NULL,
    `order_no` VARCHAR(32) DEFAULT NULL,
    `local_amount` DECIMAL(10,2) DEFAULT NULL,
    `channel_amount` DECIMAL(10,2) DEFAULT NULL,
    `local_status` TINYINT DEFAULT NULL,
    `channel_status` VARCHAR(20) DEFAULT NULL,
    `handle_status` TINYINT NOT NULL DEFAULT '0',
    `handle_by` VARCHAR(50) DEFAULT NULL,
    `handled_at` DATETIME DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_date_type` (`reconcile_date`, `exception_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对账异常表';

-- ----------------------------
-- Table structure for reconciliation_channel_data (渠道对账单数据临时表)
-- ----------------------------
DROP TABLE IF EXISTS `reconciliation_channel_data`;
CREATE TABLE `reconciliation_channel_data` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `reconcile_date` DATE NOT NULL,
    `channel` VARCHAR(20) NOT NULL,
    `channel_trade_no` VARCHAR(64) NOT NULL,
    `channel_order_no` VARCHAR(64) DEFAULT NULL,
    `trade_time` DATETIME NOT NULL,
    `amount` DECIMAL(10,2) NOT NULL,
    `trade_type` VARCHAR(10) NOT NULL,
    `trade_status` VARCHAR(20) NOT NULL,
    `matched` TINYINT NOT NULL DEFAULT '0',
    `local_transaction_no` VARCHAR(32) DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_channel_trade` (`reconcile_date`, `channel`, `channel_trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='渠道对账单数据';

-- ----------------------------
-- Table structure for payment_logs (支付操作日志表)
-- ----------------------------
DROP TABLE IF EXISTS `payment_logs`;
CREATE TABLE `payment_logs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `transaction_no` VARCHAR(32) DEFAULT NULL,
    `order_no` VARCHAR(32) DEFAULT NULL,
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型：CREATE/PREPAY/NOTIFY等',
    `request_data` JSON DEFAULT NULL,
    `response_data` JSON DEFAULT NULL,
    `result` TINYINT DEFAULT NULL COMMENT '0-失败 1-成功',
    `error_msg` VARCHAR(500) DEFAULT NULL,
    `ip_address` VARCHAR(40) DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_transaction_action` (`transaction_no`, `action`),
    KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付操作日志表';
