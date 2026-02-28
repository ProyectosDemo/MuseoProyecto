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
-- Table structure for table `orden`
--

DROP TABLE IF EXISTS `orden`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orden` (
  `id_orden` int NOT NULL AUTO_INCREMENT,
  `id_obra` int NOT NULL,
  `id_cliente` int NOT NULL,
  `id_trabajador` int NOT NULL,
  `precio_obra` decimal(10,2) NOT NULL,
  `iva` decimal(10,2) GENERATED ALWAYS AS ((`precio_obra` * 0.16)) STORED,
  `ganancia_museo` decimal(10,2) NOT NULL,
  `total` decimal(10,2) GENERATED ALWAYS AS ((`precio_obra` + `iva`)) STORED,
  `fecha_orden` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `estatus` enum('Pendiente','Concretada','Cancelada') DEFAULT 'Pendiente',
  PRIMARY KEY (`id_orden`),
  KEY `id_obra` (`id_obra`),
  KEY `id_cliente` (`id_cliente`),
  KEY `id_trabajador` (`id_trabajador`),
  CONSTRAINT `orden_ibfk_1` FOREIGN KEY (`id_obra`) REFERENCES `obra` (`id_obra`),
  CONSTRAINT `orden_ibfk_2` FOREIGN KEY (`id_cliente`) REFERENCES `cliente` (`id_cliente`),
  CONSTRAINT `orden_ibfk_3` FOREIGN KEY (`id_trabajador`) REFERENCES `trabajador` (`id_trabajador`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orden`
--

LOCK TABLES `orden` WRITE;
/*!40000 ALTER TABLE `orden` DISABLE KEYS */;
INSERT INTO `orden` (`id_orden`, `id_obra`, `id_cliente`, `id_trabajador`, `precio_obra`, `ganancia_museo`, `fecha_orden`, `estatus`) VALUES (1,1,3,1,2000.00,400.50,'2026-02-21 20:00:00','Pendiente'),(2,1,3,1,2000.00,400.00,'2026-02-21 20:00:00','Pendiente'),(4,1,3,1,2000.00,400.00,'2026-02-21 20:00:00','Pendiente'),(5,1,3,1,2000.00,400.00,'2026-02-21 20:00:00','Pendiente'),(6,2,3,1,1500.50,300.10,'2026-02-21 20:00:00','Pendiente'),(7,2,3,1,1500.50,300.10,'2026-02-21 20:00:00','Concretada');
/*!40000 ALTER TABLE `orden` ENABLE KEYS */;
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
