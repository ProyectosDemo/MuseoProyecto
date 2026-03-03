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
INSERT INTO `artista` VALUES (2,'Pr Quoute','1990-05-19','Chile','75% Serotonin 20% Dinner Time 5% Pure Hate','https://images-ext-1.discordapp.net/external/pI-IR8vEBlxHGKg6MBllpemZqXuy21ZiJVbQT1euEyY/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/201505280545128448/536a938685ec618f2129cfd4ff287fca.png?format=webp&quality=lossless'),(4,'Yesenia','2002-02-01','Venezuela','3cm de alto y enojada','https://media.discordapp.net/attachments/731352148947894363/1477514662089064489/Screenshot_2026-02-26-23-56-47-312_com.discord-edit.jpg?ex=69a5b302&is=69a46182&hm=3523044e3162f0062fe719799dbdb3954ab475e912cf44a900836a8e6b8576ae&=&format=webp'),(5,'Yisus','2002-02-01','Venezuela','Giga chad the kawai','https://images-ext-1.discordapp.net/external/KtJ4_uP4a9EmluKTEagLb2vTD4rjDh8aro3b7ZaUup8/https/imgur.com/uE6Gyz5.png?format=webp&quality=lossless&width=517&height=693'),(6,'Claude Monet','1840-11-14','Francesa','Pintor francés, fundador del impresionismo, conocido por sus paisajes y el uso de la luz.','https://grovegallery.com/cdn/shop/articles/most-popular-pop-art-artist_d21ed650-7306-4411-aa74-d3d49f98c148.jpg?v=1749715530'),(7,'Pablo Picasso','1881-10-25','Española','Pintor y escultor español, cofundador del cubismo y uno de los artistas más influyentes del siglo XX.','https://upload.wikimedia.org/wikipedia/commons/9/98/Pablo_picasso_1.jpg'),(8,'Hatsune Miku','1907-07-06','Japonesa','La mas grande vocaloid','https://images-ext-1.discordapp.net/external/vFdfws7kZtUNVeTZd4pXHAzacTjVmLQ0lwjWStH1ch4/https/mudae.net/uploads/6428828/qjf2Pcw~tlsSkZh.png?format=webp&quality=lossless'),(9,'Auguste Rodin','1840-11-12','Francesa','Escultor francés, considerado el padre de la escultura moderna, famoso por \"El Pensador\".','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS6m1-xlIPRA-wzQPQXtDWy3GaGsnelmpNGiA&s'),(10,'Andy Warhol','1928-08-06','Estadounidense','Artista estadounidense, líder del pop art, conocido por sus retratos de la cultura popular y medios de comunicación.','https://static-assets.artlogic.net/w_1200,c_limit,f_auto,fl_lossy,q_auto/ws-benbrownfinearts/usr/images/artists/group_images_override/items/39/3926a038ab63445bb308aec7ea82d2a4/warhol-headshot.jpg'),(11,'Niño Interrogación','1452-04-15','Italiana','Polímata italiano del Renacimiento, famoso por pinturas como \"La Gioconda\" y \"La Última Cena\".','https://media.discordapp.net/attachments/821932100507795457/1478454992036827186/1772561533141.jpg?ex=69a87602&is=69a72482&hm=7b19e55fd2a51b7d481f75feeb40010711cad86188f2fae57c5513a1eed60fb4&=&format=webp'),(12,'Superman','1904-05-11','Krypton','Pintor surrealista español, conocido por sus imágenes oníricas y excéntricas.','https://media.discordapp.net/attachments/821932100507795457/1478454394927448104/1772561386126.jpg?ex=69a87574&is=69a723f4&hm=95ca499a820a86eeee998a7797ad3a90f0a04146b0d3535dd62a069a56709245&=&format=webp&width=707&height=694'),(13,'Georgia O’Keeffe','1887-11-15','Estadounidense','Pintora estadounidense, conocida como la “madre del modernismo americano”, famosa por sus flores y paisajes del desierto.','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCRwROoIoWTntWZShTIkeK59m6b53e6b2cLw&s');
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

-- Dump completed on 2026-03-03 14:15:56
