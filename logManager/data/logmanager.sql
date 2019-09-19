-- create database log_manager --
CREATE DATABASE log_manager DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

-- ----------------------------
-- Table structure for log_collects
-- ----------------------------
DROP TABLE IF EXISTS `log_collects`;
CREATE TABLE `log_collects` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host` varchar(255) NOT NULL,
  `topic` varchar(255) NOT NULL,
  `logPath` varchar(1024) NOT NULL,
  `pid` int(11) NOT NULL,
  `applyPath` varchar(1024) NOT NULL,
  `createTime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for log_project
-- ----------------------------
DROP TABLE IF EXISTS `log_project`;
CREATE TABLE `log_project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pname` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `applyPath` varchar(10240) NOT NULL,
  `createTime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pname` (`pname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
