-- MySQL dump 10.13  Distrib 8.0.33, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: snet
-- ------------------------------------------------------
-- Server version	8.0.33-0ubuntu0.23.04.2

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
-- Table structure for table `token`
--

CREATE DATABASE snet;
CREATE USER 'socialnet' IDENTIFIED BY 'socialnet';
GRANT ALL PRIVILEGES ON snet.* TO 'socialnet';
GRANT ALL PRIVILEGES ON snet.* TO 'socialnet'@'%';

use snet;

DROP TABLE IF EXISTS `token`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `token` (
  `hash` char(64) NOT NULL,
  `user_id` char(36) NOT NULL,
  `expires` timestamp NOT NULL,
  KEY `token_users_user_id_fk` (`user_id`),
  KEY `token_hash_index` (`hash`),
  CONSTRAINT `token_users_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `token`
--

/*!40000 ALTER TABLE `token` DISABLE KEYS */;
INSERT INTO `token` VALUES ('cf0dc07992a6378812d2410e674d5cbf1cf0c19aa4ffc06d6f7af3a3e259eb2c','56d721b0-a5cc-a00e-a77d-72b0919cb457','2023-07-20 14:06:00');
INSERT INTO `token` VALUES ('d441d97565c00e3e8abcb65f91fe84b3a11489006f4694c578e443dca52bab97','fb0e9288-e4d2-9561-166a-eb34120bc3c3','2023-07-20 14:06:00');
/*!40000 ALTER TABLE `token` ENABLE KEYS */;

--
-- Table structure for table `user_credentials`
--

DROP TABLE IF EXISTS `user_credentials`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_credentials` (
  `user_id` char(36) NOT NULL,
  `password` char(60) DEFAULT NULL,
  KEY `user_credentials_user_id_index` (`user_id`),
  CONSTRAINT `user_credentials_users_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_credentials`
--

/*!40000 ALTER TABLE `user_credentials` DISABLE KEYS */;
INSERT INTO `user_credentials` VALUES ('56d721b0-a5cc-a00e-a77d-72b0919cb457','$2a$04$E/4SsLEqOan35T7geJaDAeE0vwlN.QjFGhcRBRRr8PR3nDD0luec.');
INSERT INTO `user_credentials` VALUES ('fb0e9288-e4d2-9561-166a-eb34120bc3c3','$2a$04$ur6ft08qwVsQw/DLUJBy5OvKE9rVZ14vFZzfC29dhmV8x/lc219xu');
/*!40000 ALTER TABLE `user_credentials` ENABLE KEYS */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `user_id` char(36) NOT NULL,
  `first_name` varchar(64) NOT NULL,
  `second_name` varchar(64) DEFAULT NULL,
  `sex` enum('male','female') DEFAULT NULL,
  `biography` text,
  `city` varchar(64) DEFAULT NULL,
  `birthdate` timestamp NULL DEFAULT NULL,
  KEY `users_user_id_index` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('56d721b0-a5cc-a00e-a77d-72b0919cb457','Ivan1','Frolov',NULL,'Hokkey','Moskva','2002-02-10 21:00:00');
INSERT INTO `users` VALUES ('fb0e9288-e4d2-9561-166a-eb34120bc3c3','Masha1','Frolova',NULL,'Dance','Piter','2003-02-10 21:00:00');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-07-19 20:07:57
