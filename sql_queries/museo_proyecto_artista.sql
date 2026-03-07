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
INSERT INTO `artista` VALUES (2,'Pr Quote','1990-05-19','Chile','75% Serotonin 20% Dinner Time 5% Pure Hate','images\\artistas\\pr_quote.png'),(4,'Yesenia','2002-02-01','Venezuela','3cm de alto y enojada','images\\artistas\\yesenia.png'),(5,'Yisus','2002-02-01','Venezuela','Giga chad the kawai','images\\artistas\\yisus.png'),(6,'Claude Monet','1840-11-14','Francesa','Pintor francés, fundador del impresionismo, conocido por sus paisajes y el uso de la luz.','images\\artistas\\claude.jpg'),(7,'Pablo Picasso','1881-10-25','Española','Pintor y escultor español, cofundador del cubismo y uno de los artistas más influyentes del siglo XX.','images\\artistas\\picaso.jpg'),(8,'Hatsune Miku','1907-07-06','Japonesa','La mas grande vocaloid','images\\artistas\\miku.png'),(9,'Auguste Rodin','1840-11-12','Francesa','Escultor francés, considerado el padre de la escultura moderna, famoso por \"El Pensador\".','images\\artistas\\auguste.jpg'),(10,'Andy Warhol','1928-08-06','Estadounidense','Artista estadounidense, líder del pop art, conocido por sus retratos de la cultura popular y medios de comunicación.','images\\artistas\\andy.png'),(11,'Niño Interrogación','1452-04-15','Italiana','Polímata italiano del Renacimiento, famoso por pinturas como \"La Gioconda\" y \"La Última Cena\".','images\\artistas\\niño.jpg'),(13,'Georgia O’Keeffe','1887-11-15','Estadounidense','Pintora estadounidense, conocida como la “madre del modernismo americano”, famosa por sus flores y paisajes del desierto.','images\\artistas\\georgia.jpg');
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

-- Dump completed on 2026-03-07 14:15:02
