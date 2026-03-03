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
  `peso` int DEFAULT NULL,
  `dimensiones` varchar(45) DEFAULT NULL,
  `material` varchar(45) DEFAULT NULL,
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
INSERT INTO `obra` VALUES (1,'La muerte de CSP',2,2,2000,'1889-05-31','DISPONIBLE','https://media.discordapp.net/attachments/936121149203574814/1474599528547090492/image.png?ex=69a30114&is=69a1af94&hm=ec493a7604c3949f011c1bacaf3eb8769ed2c7f64dd8c301731f4c5460e4183d&=&format=webp&quality=lossless&width=694&height=694',12,'70x50x5','Óleo sobre lienzo'),(2,'El Indie mas chud',2,2,2000,'1889-06-01','DISPONIBLE','https://media.discordapp.net/attachments/936121149203574814/1399975429556932728/image.png?ex=69a7b838&is=69a666b8&hm=24a4df4e9ce892adcb050a9f8ce1cae92d4d5cfd4454644b009f446f28a8a2be&=&format=webp&quality=lossless&width=554&height=694',10,'65x50x5','Acrílico sobre lienzo'),(5,'Will never see the light of day',2,2,2000,'1889-05-31','DISPONIBLE','https://media.discordapp.net/attachments/936121149203574814/943326128024719381/23.png?ex=69a2f52b&is=69a1a3ab&hm=d05428598279ea48b2f766d9f832328442778e8f131d9acbdf0f51e77a9b9066&=&format=webp&quality=lossless&width=822&height=608',15,'80x60x5','Óleo sobre lienzo'),(6,'Dreams of a frog',2,2,2000,'1889-05-31','DISPONIBLE','https://media.discordapp.net/attachments/936121149203574814/944428862866067497/des_frog.png?ex=69a303ab&is=69a1b22b&hm=5efc2eb4fd297f8110b7a59d9987384487b38d4dba4392fdcd79a6da2aa3b423&=&format=webp&quality=lossless&width=694&height=694',8,'50x50x5','Acrílico sobre lienzo'),(15,'Escultura Alfa',4,6,3500,'1890-01-15','DISPONIBLE','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT1RxRBZejkWV6Dbj8iHBDYUnum0wRiJP4JHg&s',3500,'200x80x70','Mármol'),(16,'Escultura Beta',6,18,4200,'1891-02-20','RESERVADA','https://www.oldskull.net/wp-content/uploads/2023/05/Jonathan-Hateley-esculturas-de-bronce-yoga.jpg',4200,'220x90x75','Bronce'),(17,'Escultura Gamma',5,3,2800,'1892-03-10','DISPONIBLE','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTE9NcTZNa0LND6CyAD44WJFoBcge880IDUaQ&s',2800,'180x70x65','Mármol'),(18,'Escultura Delta',9,9,5000,'1893-04-05','RESERVADA','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTOMUKBEsNbMEYGIntM5Lq6zOifOgJjp7R8Tg&s',5000,'250x100x90','Bronce'),(19,'Escultura Épsilon',10,7,3300,'1894-05-12','RESERVADA','https://www.cinconoticias.com/wp-content/uploads/esculturas-italianas.jpg',3300,'190x80x70','Mármol'),(20,'Escultura Zeta',7,17,3900,'1895-06-18','DISPONIBLE','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRI55AnJ1-NsEBER99Nvu6BZGgxDLz4cf4cOQ&s',3900,'210x85x80','Bronce'),(21,'Escultura Eta',8,15,4500,'1896-07-20','RESERVADA','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQxTz2G2BCk9wn_rPOEXmlDFT0FoCIw2j9gEw&s',4500,'230x95x85','Mármol'),(22,'Escultura Theta',11,14,4700,'1897-08-22','DISPONIBLE','https://www.oldskull.net/wp-content/uploads/2015/09/Animal-Sculptures-beth-cavener-1.jpg',4700,'240x100x90','Bronce');
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

-- Dump completed on 2026-03-03 14:15:56
