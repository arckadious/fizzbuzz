USE DB;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+01:00";

CREATE TABLE IF NOT EXISTS `MESSAGES_REQUEST` (
  `ID` MEDIUMINT NOT NULL AUTO_INCREMENT,
  `COR_ID` varchar(36) CHARACTER SET utf8mb3 DEFAULT NULL UNIQUE,
  `SERVICE_ADDRESS` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `HOST` varchar(64) CHARACTER SET utf8mb3 DEFAULT NULL,
  `MSG` longtext CHARACTER SET utf8mb3 DEFAULT NULL,
  `DT_CREATION` datetime DEFAULT CURRENT_TIMESTAMP,
  `APP_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CHECKSUM` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`)

) DEFAULT CHARSET=latin1;

ALTER TABLE `MESSAGES_REQUEST`
  ADD KEY `INDEX_DT_CREATION` (`DT_CREATION`),
  ADD KEY `INDEX_CHECKSUM` (`CHECKSUM`),
  ADD KEY `IDX_COR_ID` (`COR_ID`);


-- INSERT INTO `MESSAGES_REQUEST` (`MSG`, `HOST`, `APP_NAME`, `SERVICE_ADDRESS`, `CHECKSUM`) VALUES ('testrequest', '123.456.678.910', 'fizzbuzz', 'http://api.localhost/v1/fizzbuzz', 'euzhghzoig');

CREATE TABLE IF NOT EXISTS `MESSAGES_RESPONSE` (
  `ID` MEDIUMINT NOT NULL AUTO_INCREMENT,
  `COR_ID` varchar(36) CHARACTER SET utf8mb3 DEFAULT NULL UNIQUE,
  `STATUS` varchar(64) CHARACTER SET utf8mb3 DEFAULT NULL,
  `MSG` longtext CHARACTER SET utf8mb3 DEFAULT NULL,
  `DT_CREATION` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`)

) DEFAULT CHARSET=latin1;

ALTER TABLE `MESSAGES_RESPONSE`
  ADD KEY `IDX_COR_ID` (`COR_ID`),
  ADD KEY `INDEX_DT_CREATION` (`DT_CREATION`);