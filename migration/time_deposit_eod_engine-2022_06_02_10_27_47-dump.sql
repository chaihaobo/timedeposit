-- MySQL dump 10.13  Distrib 8.0.29, for macos12.2 (arm64)
--
-- Host: 192.168.31.16    Database: time_deposit
-- ------------------------------------------------------
-- Server version	8.0.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `t_flow_node_logs`
--
create
    database if not exists time_deposit;
use time_deposit;
DROP TABLE IF EXISTS `t_flow_node_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;

CREATE TABLE `t_flow_node_logs`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `account_id`  varchar(50)  DEFAULT NULL,
    `flow_id`     varchar(50)  DEFAULT NULL,
    `flow_name`   varchar(255) DEFAULT NULL,
    `node_name`   varchar(255) DEFAULT NULL,
    `node_result` varchar(255) DEFAULT NULL,
    `node_msg`    varchar(255) DEFAULT NULL,
    `create_time` datetime     DEFAULT NULL,
    `update_time` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1063 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_flow_node_query_logs`
--

DROP TABLE IF EXISTS `t_flow_node_query_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_flow_node_query_logs`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `flow_id`     varchar(50)  NOT NULL COMMENT 'flow id',
    `node_name`   varchar(100) NOT NULL,
    `query_type`  varchar(100) DEFAULT NULL,
    `data`        longtext,
    `create_time` datetime     DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY           `t_flow_node_query_log_flow_id_index` (`flow_id`)
) ENGINE=InnoDB AUTO_INCREMENT=510 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_flow_node_relations`
--

DROP TABLE IF EXISTS `t_flow_node_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_flow_node_relations`
(
    `id`          int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
    `flow_name`   varchar(63) DEFAULT NULL COMMENT 'Flow type name',
    `node_name`   varchar(63) DEFAULT NULL COMMENT 'Flow node name',
    `result_code` varchar(63) DEFAULT NULL COMMENT 'description of node',
    `next_node`   varchar(63) DEFAULT NULL COMMENT 'the flow node name of next step',
    `create_time` datetime    DEFAULT NULL COMMENT 'Create Time',
    `update_time` datetime    DEFAULT NULL COMMENT 'Update Time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb3 COMMENT='flow node info';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_flow_nodes`
--

DROP TABLE IF EXISTS `t_flow_nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_flow_nodes`
(
    `id`          int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
    `flow_name`   varchar(63)  DEFAULT NULL COMMENT 'Flow type name',
    `node_name`   varchar(63)  DEFAULT NULL COMMENT 'Flow node name',
    `node_path`   varchar(255) DEFAULT NULL,
    `node_detail` varchar(255) DEFAULT NULL COMMENT 'description of node',
    `create_time` datetime     DEFAULT NULL COMMENT 'Create Time',
    `update_time` datetime     DEFAULT NULL COMMENT 'Update Time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `flow_name_node_name` (`flow_name`,`node_name`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb3 COMMENT='flow node info';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_flow_task_infos`
--

DROP TABLE IF EXISTS `t_flow_task_infos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_flow_task_infos`
(
    `id`            int unsigned NOT NULL AUTO_INCREMENT,
    `flow_id`       varchar(50)  DEFAULT NULL,
    `account_id`    varchar(50)  DEFAULT NULL,
    `flow_name`     varchar(255) NOT NULL,
    `flow_status`   varchar(255) DEFAULT NULL,
    `cur_node_name` varchar(255) DEFAULT NULL,
    `cur_status`    varchar(255) DEFAULT NULL,
    `end_status`    varchar(255) DEFAULT NULL,
    `start_time`    datetime     DEFAULT NULL,
    `end_time`      datetime     DEFAULT NULL,
    `enable`        tinyint(1) default 1 null,
    `create_time`   datetime     DEFAULT NULL,
    `update_time`   datetime     DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=185 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_flow_transactions`
--

DROP TABLE IF EXISTS `t_flow_transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_flow_transactions`
(
    `id`                   int unsigned NOT NULL AUTO_INCREMENT,
    `trans_id`             varchar(255)  DEFAULT NULL,
    `flow_id`              varchar(100)  DEFAULT NULL,
    `mambu_trans_id`       varchar(100)  DEFAULT NULL,
    `terminal_rrn`         varchar(255)  DEFAULT NULL,
    `source_account_no`    varchar(50)   DEFAULT NULL,
    `source_account_name`  varchar(255)  DEFAULT NULL,
    `benefit_account_no`   varchar(50)   DEFAULT NULL,
    `benefit_account_name` varchar(255)  DEFAULT NULL,
    `amount`               double        DEFAULT NULL,
    `channel`              varchar(255)  DEFAULT NULL,
    `transaction_type`     varchar(50)   DEFAULT NULL,
    `result`               int           DEFAULT '-1' COMMENT '1:succeed, 0:failed',
    `encoded_key`          varchar(255)  DEFAULT NULL,
    `create_time`          datetime      DEFAULT NULL,
    `update_time`          datetime      DEFAULT NULL,
    `error_msg`            varchar(1000) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=82 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `t_mambu_request_logs`
--

DROP TABLE IF EXISTS `t_mambu_request_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_mambu_request_logs`
(
    `id`            bigint NOT NULL AUTO_INCREMENT,
    `flow_id`       varchar(100) DEFAULT NULL,
    `account_id`    varchar(100) DEFAULT NULL,
    `node_name`     varchar(100) DEFAULT NULL,
    `type`          varchar(50)  DEFAULT NULL,
    `request_url`   varchar(255) DEFAULT NULL,
    `request_body`  text,
    `response_code` int          DEFAULT NULL,
    `response_body` text,
    `error`         text,
    `create_time`   datetime     DEFAULT CURRENT_TIMESTAMP,
    `update_time`   datetime     DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY             `t_mambu_request_log_flow_id_index` (`flow_id`)
) ENGINE=InnoDB AUTO_INCREMENT=717 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (1, 'eod_flow', 'start_node', 'StartNode', 'start_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (2, 'eod_flow', 'undo_maturity_node', 'UndoMaturityNode', 'undo_maturity_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (3, 'eod_flow', 'start_new_maturity_node', 'StartNewMaturityNode', 'start_new_maturity_node',
        '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (4, 'eod_flow', 'apply_profit_node', 'ApplyProfitNode', 'apply_profit_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (5, 'eod_flow', 'withdraw_netprofit_node', 'WithdrawNetprofitNode', 'withdraw_netprofit_node',
        '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (6, 'eod_flow', 'deposit_netprofit_node', 'DepositNetprofitNode', 'deposit_netprofit_node',
        '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (7, 'eod_flow', 'withdraw_balance_node', 'WithdrawBalanceNode', 'withdraw_balance_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (8, 'eod_flow', 'deposit_balance_node', 'DepositBalanceNode', 'deposit_balance_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (9, 'eod_flow', 'search_last_profit_applied_node', 'SearchLastProfitAppliedNode',
        'search_last_profit_applied_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (10, 'eod_flow', 'withdraw_additional_profit_node', 'WithdrawAdditionalProfitNode',
        'withdraw_additional_profit_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (11, 'eod_flow', 'deposit_additional_profit_node', 'DepositAdditionalProfitNode',
        'deposit_additional_profit_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (12, 'eod_flow', 'patch_account_node', 'PatchAccountNode', 'patch_account_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (13, 'eod_flow', 'close_account_node', 'CloseAccountNode', 'close_account_node', '2022-05-30 10:43:44',
        '2022-05-30 10:43:47');
INSERT INTO time_deposit.t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time,
                                       update_time)
VALUES (14, 'eod_flow', 'end_node', 'EndNode', 'end_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');

INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (1, 'eod_flow', 'start_node', 'success', 'undo_maturity_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (2, 'eod_flow', 'undo_maturity_node', 'success', 'start_new_maturity_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (3, 'eod_flow', 'start_new_maturity_node', 'success', 'apply_profit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (4, 'eod_flow', 'apply_profit_node', 'success', 'withdraw_netprofit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (5, 'eod_flow', 'withdraw_netprofit_node', 'success', 'deposit_netprofit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (6, 'eod_flow', 'deposit_netprofit_node', 'success', 'withdraw_balance_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (7, 'eod_flow', 'withdraw_balance_node', 'success', 'deposit_balance_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (8, 'eod_flow', 'deposit_balance_node', 'success', 'deposit_additional_profit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (10, 'eod_flow', 'withdraw_additional_profit_node', 'success', 'patch_account_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (11, 'eod_flow', 'deposit_additional_profit_node', 'success', 'withdraw_additional_profit_node',
        '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (12, 'eod_flow', 'patch_account_node', 'success', 'close_account_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (13, 'eod_flow', 'close_account_node', 'success', 'end_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (14, 'eod_flow', 'undo_maturity_node', 'skip', 'start_new_maturity_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (15, 'eod_flow', 'start_new_maturity_node', 'skip', 'apply_profit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (16, 'eod_flow', 'apply_profit_node', 'skip', 'withdraw_netprofit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (17, 'eod_flow', 'withdraw_netprofit_node', 'skip', 'deposit_netprofit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (18, 'eod_flow', 'deposit_netprofit_node', 'skip', 'withdraw_balance_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (19, 'eod_flow', 'withdraw_balance_node', 'skip', 'deposit_balance_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (20, 'eod_flow', 'deposit_balance_node', 'skip', 'deposit_additional_profit_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (21, 'eod_flow', 'withdraw_additional_profit_node', 'skip', 'patch_account_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (22, 'eod_flow', 'deposit_additional_profit_node', 'skip', 'withdraw_additional_profit_node',
        '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (23, 'eod_flow', 'patch_account_node', 'skip', 'close_account_node', '2022-05-30 10:50:11',
        '2022-05-30 10:50:12');
INSERT INTO time_deposit.t_flow_node_relations (id, flow_name, node_name, result_code, next_node,
                                                create_time, update_time)
VALUES (24, 'eod_flow', 'close_account_node', 'skip', 'end_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');


/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-02 10:27:49
