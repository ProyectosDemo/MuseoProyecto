-- MySQL dump 10.13  Distrib 8.0.43, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: museo_proyecto
-- ------------------------------------------------------
-- Server version	9.4.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `obra`
--

DROP TABLE IF EXISTS `obra`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `obra` (
  `id_obra` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `id_artista` int DEFAULT NULL,
  `id_genero` int DEFAULT NULL,
  `precio_obra` int DEFAULT NULL,
  `fecha_creacion` date DEFAULT NULL,
  `status` enum('DISPONIBLE','RESERVADA','VENDIDA') NOT NULL DEFAULT 'DISPONIBLE',
  `foto` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id_obra`),
  KEY `id_artista` (`id_artista`),
  KEY `id_genero` (`id_genero`),
  CONSTRAINT `obra_ibfk_1` FOREIGN KEY (`id_artista`) REFERENCES `artista` (`id_artista`),
  CONSTRAINT `obra_ibfk_2` FOREIGN KEY (`id_genero`) REFERENCES `genero` (`id_genero`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `obra`
--

LOCK TABLES `obra` WRITE;
/*!40000 ALTER TABLE `obra` DISABLE KEYS */;
INSERT INTO `obra` VALUES (1,'La muerte de CSP',2,2,2000,'1889-05-31','DISPONIBLE','images/obras/la_muerte_de_csp.png'),(2,'El Indie mas chud',2,2,2000,'1889-06-01','DISPONIBLE','images\\obras\\el_indie_mas_chud.png'),(5,'Will never see the light of day',2,2,2000,'1889-05-31','DISPONIBLE','images\\obras\\will_never_see_the_light_of_day.png'),(6,'Dreams of a frog',2,2,2000,'1889-05-31','DISPONIBLE','images\\obras\\dreams_of_a_frog.png'),(15,'Escultura Alfa',4,6,3500,'1890-01-15','DISPONIBLE','images\\obras\\escultura1.jpg'),(16,'Escultura Beta',6,18,4200,'1891-02-20','DISPONIBLE','images\\obras\\escultura2.jpg'),(17,'Escultura Gamma',5,3,2800,'1892-03-10','DISPONIBLE','images\\obras\\escultura3.jpg'),(18,'Escultura Delta',9,9,5000,'1893-04-05','DISPONIBLE','images\\obras\\escultura4.jpg'),(19,'Escultura Épsilon',10,7,3300,'1894-05-12','DISPONIBLE','images\\obras\\escultura5.jpg'),(20,'Escultura Zeta',7,17,3900,'1895-06-18','DISPONIBLE','images\\obras\\escultura6.jpg'),(21,'Escultura Eta',8,15,4500,'1896-07-20','DISPONIBLE','images\\obras\\escultura7.jpg'),(22,'Escultura Theta',11,14,4700,'1897-08-22','DISPONIBLE','images\\obras\\escultura8.jpg');
/*!40000 ALTER TABLE `obra` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-06 21:02:01
