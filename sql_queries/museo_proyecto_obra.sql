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
  `estatus` enum('Disponible','Reservada','Vendida') DEFAULT 'Disponible',
  `foto` varchar(255) DEFAULT NULL,
  `material` varchar(50) DEFAULT NULL,
  `peso` int DEFAULT NULL,
  `dimensiones` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_obra`),
  KEY `id_artista` (`id_artista`),
  KEY `id_genero` (`id_genero`),
  CONSTRAINT `obra_ibfk_1` FOREIGN KEY (`id_artista`) REFERENCES `artista` (`id_artista`),
  CONSTRAINT `obra_ibfk_2` FOREIGN KEY (`id_genero`) REFERENCES `genero` (`id_genero`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `obra`
--

LOCK TABLES `obra` WRITE;
/*!40000 ALTER TABLE `obra` DISABLE KEYS */;
INSERT INTO `obra` VALUES (1,'La muerte de CSP',2,2,2000,'1889-05-31','Disponible','https://media.discordapp.net/attachments/936121149203574814/1474599528547090492/image.png?ex=69a30114&is=69a1af94&hm=ec493a7604c3949f011c1bacaf3eb8769ed2c7f64dd8c301731f4c5460e4183d&=&format=webp&quality=lossless&width=694&height=694','oreo sobre oreo',3,'74 x 92 cm'),(2,'El indie mas chud',2,2,2000,'1889-05-31','Disponible','https://media.discordapp.net/attachments/936121149203574814/1399975429556932728/image.png?ex=69a27238&is=69a120b8&hm=4fd7700f1184721f8e1608011b1e792a2598c9546b6060b4717a41beaf65315b&=&format=webp&quality=lossless&width=554&height=694','oreo sobre oreo',3,'74 x 92 cm'),(5,'Will never see the light of day',2,2,2000,'1889-05-31','Disponible','https://media.discordapp.net/attachments/936121149203574814/943326128024719381/23.png?ex=69a2f52b&is=69a1a3ab&hm=d05428598279ea48b2f766d9f832328442778e8f131d9acbdf0f51e77a9b9066&=&format=webp&quality=lossless&width=822&height=608','oreo sobre oreo',3,'74 x 92 cm'),(6,'Dreams of a frog',2,2,2000,'1889-05-31','Disponible','https://media.discordapp.net/attachments/936121149203574814/944428862866067497/des_frog.png?ex=69a303ab&is=69a1b22b&hm=5efc2eb4fd297f8110b7a59d9987384487b38d4dba4392fdcd79a6da2aa3b423&=&format=webp&quality=lossless&width=694&height=694','oreo sobre oreo',3,'74 x 92 cm');
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

-- Dump completed on 2026-02-28 15:52:37
