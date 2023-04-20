USE DB;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+01:00";


--
-- Base de données :  `DB`
--

-- --------------------------------------------------------

--
-- Structure de la table `MESSAGES_REQUEST`
--

CREATE TABLE IF NOT EXISTS `MESSAGES_REQUEST` (
  `ID` varchar(36) CHARACTER SET utf8mb3 DEFAULT NULL,
  `MSG` longtext CHARACTER SET utf8mb3 DEFAULT NULL,
  `HOST` varchar(64) CHARACTER SET utf8mb3 DEFAULT NULL,
  `TYPE` varchar(64) CHARACTER SET utf8mb3 DEFAULT NULL,
  `DT_CREATION` datetime DEFAULT NULL,
  `APP_NAME` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SERVICE_ADDRESS` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CHECKSUM` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Index pour les tables exportées
--

--
-- Index pour la table `MESSAGES_REQUEST`
--
ALTER TABLE `MESSAGES_REQUEST`
  ADD KEY `IDX_TYPE_DTC_ID` (`TYPE`,`DT_CREATION`,`ID`),
  ADD KEY `IDX_ID` (`ID`),
  ADD KEY `INDEX_DT_CREATION` (`DT_CREATION`),
  ADD KEY `INDEX_CHECKSUM` (`CHECKSUM`),
  ADD KEY `MESSAGES_REQUEST_TYPE` (`TYPE`);


  
--
-- Structure de la table `MESSAGES_RESPONSE`
--

CREATE TABLE IF NOT EXISTS `MESSAGES_RESPONSE` (
  `ID` varchar(36) CHARACTER SET utf8mb3 DEFAULT NULL,
  `MSG` longtext CHARACTER SET utf8mb3 DEFAULT NULL,
  `STATUS` varchar(64) CHARACTER SET utf8mb3 DEFAULT NULL,
  `COR_ID` varchar(36) CHARACTER SET utf8mb3 DEFAULT NULL,
  `DT_CREATION` datetime DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Index pour les tables exportées
--

--
-- Index pour la table `MESSAGES_RESPONSE`
--
ALTER TABLE `MESSAGES_RESPONSE`
  ADD KEY `IDX_COR_ID` (`COR_ID`),
  ADD KEY `IDX_ID` (`ID`),
  ADD FULLTEXT KEY `MSG` (`MSG`);

