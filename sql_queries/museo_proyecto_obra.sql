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
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `obra`
--

LOCK TABLES `obra` WRITE;
/*!40000 ALTER TABLE `obra` DISABLE KEYS */;
INSERT INTO `obra` VALUES (1,'La muerte de CSP',2,8,2000,'1889-05-31','DISPONIBLE','images/obras/la_muerte_de_csp.png'),(2,'El Indie mas chud',2,8,2000,'1889-06-01','DISPONIBLE','images\\obras\\el_indie_mas_chud.png'),(5,'Will never see the light of day',2,8,2000,'1889-05-31','DISPONIBLE','images\\obras\\will_never_see_the_light_of_day.png'),(6,'Dreams of a frog',2,8,2000,'1889-05-31','DISPONIBLE','images\\obras\\dreams_of_a_frog.png'),(15,'Hekapoo',4,5,3500,'1890-01-15','DISPONIBLE','images\\obras\\yesenia1.jpeg'),(16,'Escultura Beta',6,4,4200,'1891-02-20','DISPONIBLE','images\\obras\\escultura2.jpg'),(17,'Escultura Gamma',5,4,2800,'1892-03-10','DISPONIBLE','images\\obras\\escultura3.jpg'),(18,'Escultura Delta',9,4,5000,'1893-04-05','DISPONIBLE','images\\obras\\escultura4.jpg'),(19,'Escultura Épsilon',10,4,3300,'1894-05-12','DISPONIBLE','images\\obras\\escultura5.jpg'),(20,'Escultura Zeta',7,4,3900,'1895-06-18','DISPONIBLE','images\\obras\\escultura6.jpg'),(21,'Escultura Eta',8,4,4500,'1896-07-20','VENDIDA','images\\obras\\escultura7.jpg'),(22,'Escultura Theta',11,4,4700,'1897-08-22','DISPONIBLE','images\\obras\\escultura8.jpg'),(23,'Tracer',4,5,2301,'2015-05-10','DISPONIBLE','images\\obras\\yesenia2.jpeg'),(24,'Undertale',4,5,4000,'2025-05-15','DISPONIBLE','images\\obras\\yesenia3.jpeg'),(25,'Pensamiento',5,11,1000,'2020-12-01','DISPONIBLE','images\\obras\\conceptual1.png'),(26,'Ave',5,7,900,'2015-01-30','VENDIDA','images\\obras\\fotografia1.jpg'),(27,'Manzana',6,16,3000,'2008-07-14','DISPONIBLE','images\\obras\\realismo2.png'),(28,'Arbol Arcoiris',8,3,4500,'2010-10-05','DISPONIBLE','images\\obras\\pintura1.jpg'),(29,'Control',8,15,560,'2020-12-01','DISPONIBLE','images\\obras\\surrealismo1.png'),(30,'Mal dia',13,12,5000,'2002-04-04','DISPONIBLE','images\\obras\\abstracto1.png'),(31,'Reflejo',13,14,3000,'2006-08-09','DISPONIBLE','images\\obras\\cubismo1.png'),(32,'Amor de Madre',7,16,3999,'1996-10-23','DISPONIBLE','images\\obras\\realismo1.png'),(33,'Zebra Arcoiris',7,3,2500,'2002-06-19','VENDIDA','images\\obras\\pintura2.png'),(34,'Sendero',9,12,5000,'2017-09-28','DISPONIBLE','images\\obras\\abstracto2.png'),(35,'Cuerpo',9,11,1000,'2001-01-30','DISPONIBLE','images\\obras\\conceptual2.png'),(36,'Observación',10,14,4300,'2019-11-04','DISPONIBLE','images\\obras\\cubismo2.png'),(37,'Paisaje',10,7,2850,'2001-07-24','VENDIDA','images\\obras\\fotografia2.png'),(38,'Barba',10,15,700,'2012-12-12','DISPONIBLE','images\\obras\\surrealismo2.png'),(39,'Liwiwo',2,8,9000,'2025-11-11','VENDIDA','images\\obras\\liwiwo.png');
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

-- Dump completed on 2026-03-07 13:08:05
