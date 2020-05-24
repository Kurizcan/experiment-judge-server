/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 80013
Source Host           : localhost:3306
Source Database       : db_experiment_judge

Target Server Type    : MYSQL
Target Server Version : 80013
File Encoding         : 65001

Date: 2020-05-01 16:30:44
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `problem`
-- ----------------------------
DROP TABLE IF EXISTS `problem`;
CREATE TABLE `problem` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `problemId` int(11) NOT NULL,
  `dataSource` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `solution` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `output` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `dataBase` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `problemId` (`problemId`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8;