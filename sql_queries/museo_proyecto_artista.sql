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
-- Table structure for table `artista`
--

DROP TABLE IF EXISTS `artista`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `artista` (
  `id_artista` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `fecha_nacimiento` varchar(45) DEFAULT NULL,
  `nacionalidad` varchar(50) DEFAULT NULL,
  `biografia` text,
  `foto` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id_artista`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `artista`
--

LOCK TABLES `artista` WRITE;
/*!40000 ALTER TABLE `artista` DISABLE KEYS */;
INSERT INTO `artista` VALUES (2,'Pr Quoute','1990-05-19','Chile','75% Serotonin 20% Dinner Time 5% Pure Hate','https://images-ext-1.discordapp.net/external/pI-IR8vEBlxHGKg6MBllpemZqXuy21ZiJVbQT1euEyY/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/201505280545128448/536a938685ec618f2129cfd4ff287fca.png?format=webp&quality=lossless'),(4,'Yesenia','2002-02-01','Venezuela','3cm de alto y enojada','https://media.discordapp.net/attachments/731352148947894363/1477514662089064489/Screenshot_2026-02-26-23-56-47-312_com.discord-edit.jpg?ex=69a5b302&is=69a46182&hm=3523044e3162f0062fe719799dbdb3954ab475e912cf44a900836a8e6b8576ae&=&format=webp'),(5,'Yisus','2002-02-01','Venezuela','Giga chad the kawai','https://images-ext-1.discordapp.net/external/KtJ4_uP4a9EmluKTEagLb2vTD4rjDh8aro3b7ZaUup8/https/imgur.com/uE6Gyz5.png?format=webp&quality=lossless&width=517&height=693'),(6,'Claude Monet','1840-11-14','Francesa','Pintor francés, fundador del impresionismo, conocido por sus paisajes y el uso de la luz.','https://example.com/monet.jpg'),(7,'Pablo Picasso','1881-10-25','Española','Pintor y escultor español, cofundador del cubismo y uno de los artistas más influyentes del siglo XX.','https://example.com/picasso.jpg'),(8,'Frida Kahlo','1907-07-06','Mexicana','Pintora mexicana conocida por sus autorretratos y la expresión de su dolor físico y emocional.','https://example.com/kahlo.jpg'),(9,'Auguste Rodin','1840-11-12','Francesa','Escultor francés, considerado el padre de la escultura moderna, famoso por \"El Pensador\".','https://example.com/rodin.jpg'),(10,'Andy Warhol','1928-08-06','Estadounidense','Artista estadounidense, líder del pop art, conocido por sus retratos de la cultura popular y medios de comunicación.','https://example.com/warhol.jpg'),(11,'Leonardo da Vinci','1452-04-15','Italiana','Polímata italiano del Renacimiento, famoso por pinturas como \"La Gioconda\" y \"La Última Cena\".','https://example.com/davinci.jpg'),(12,'Salvador Dalí','1904-05-11','Española','Pintor surrealista español, conocido por sus imágenes oníricas y excéntricas.','https://example.com/dali.jpg'),(13,'Georgia O’Keeffe','1887-11-15','Estadounidense','Pintora estadounidense, conocida como la “madre del modernismo americano”, famosa por sus flores y paisajes del desierto.','https://example.com/okeeffe.jpg');
/*!40000 ALTER TABLE `artista` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-03 13:46:08
