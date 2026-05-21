CREATE DATABASE  IF NOT EXISTS `museo_proyecto` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `museo_proyecto`;
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
-- Table structure for table `artista_genero`
--

DROP TABLE IF EXISTS `artista_genero`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `artista_genero` (
  `id_artista` int NOT NULL,
  `id_genero` int NOT NULL,
  PRIMARY KEY (`id_artista`,`id_genero`),
  KEY `id_genero` (`id_genero`),
  CONSTRAINT `artista_genero_ibfk_1` FOREIGN KEY (`id_artista`) REFERENCES `artista` (`id_artista`) ON DELETE CASCADE,
  CONSTRAINT `artista_genero_ibfk_2` FOREIGN KEY (`id_genero`) REFERENCES `genero` (`id_genero`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cliente`
--

DROP TABLE IF EXISTS `cliente`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cliente` (
  `id_cliente` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `telefono` varchar(20) DEFAULT NULL,
  `login` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `codigo_seguridad` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id_cliente`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `escultura`
--

DROP TABLE IF EXISTS `escultura`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `escultura` (
  `id_obra` int NOT NULL,
  `peso` int DEFAULT NULL,
  `dimensiones` varchar(100) DEFAULT NULL,
  `material` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id_obra`),
  CONSTRAINT `escultura_ibfk_1` FOREIGN KEY (`id_obra`) REFERENCES `obra` (`id_obra`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `genero`
--

DROP TABLE IF EXISTS `genero`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `genero` (
  `id_genero` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(50) NOT NULL,
  PRIMARY KEY (`id_genero`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `membresia`
--

DROP TABLE IF EXISTS `membresia`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `membresia` (
  `id_membresia` int NOT NULL AUTO_INCREMENT,
  `id_cliente` int NOT NULL,
  `id_tarjeta` int NOT NULL,
  `fecha` varchar(20) NOT NULL,
  PRIMARY KEY (`id_membresia`),
  KEY `id_cliente` (`id_cliente`),
  KEY `id_tarjeta` (`id_tarjeta`),
  CONSTRAINT `membresia_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `cliente` (`id_cliente`) ON DELETE CASCADE,
  CONSTRAINT `membresia_ibfk_2` FOREIGN KEY (`id_tarjeta`) REFERENCES `tarjeta_cliente` (`id_tarjeta`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

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
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

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
  `id_trabajador` int DEFAULT NULL,
  `fecha` varchar(45) DEFAULT NULL,
  `status` enum('PENDIENTE','CONCRETADA','CANCELADA') NOT NULL DEFAULT 'PENDIENTE',
  `direccion` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id_orden`),
  KEY `id_obra` (`id_obra`),
  KEY `id_cliente` (`id_cliente`),
  KEY `id_trabajador` (`id_trabajador`),
  CONSTRAINT `orden_ibfk_1` FOREIGN KEY (`id_obra`) REFERENCES `obra` (`id_obra`),
  CONSTRAINT `orden_ibfk_2` FOREIGN KEY (`id_cliente`) REFERENCES `cliente` (`id_cliente`),
  CONSTRAINT `orden_ibfk_3` FOREIGN KEY (`id_trabajador`) REFERENCES `trabajador` (`id_trabajador`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `preguntas`
--

DROP TABLE IF EXISTS `preguntas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `preguntas` (
  `id_pregunta` bigint unsigned NOT NULL AUTO_INCREMENT,
  `id_cliente` int NOT NULL,
  `pregunta` varchar(255) NOT NULL,
  `respuesta` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id_pregunta`),
  UNIQUE KEY `id_pregunta` (`id_pregunta`),
  KEY `fk_cliente` (`id_cliente`),
  CONSTRAINT `fk_cliente` FOREIGN KEY (`id_cliente`) REFERENCES `cliente` (`id_cliente`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tarjeta_cliente`
--

DROP TABLE IF EXISTS `tarjeta_cliente`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tarjeta_cliente` (
  `id_tarjeta` int NOT NULL AUTO_INCREMENT,
  `id_cliente` int NOT NULL,
  `numero_tarjeta` varchar(20) NOT NULL,
  `tipo` varchar(20) NOT NULL,
  PRIMARY KEY (`id_tarjeta`),
  KEY `id_cliente` (`id_cliente`),
  CONSTRAINT `tarjeta_cliente_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `cliente` (`id_cliente`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `trabajador`
--

DROP TABLE IF EXISTS `trabajador`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `trabajador` (
  `id_trabajador` int NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `login` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `admin` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id_trabajador`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-15 11:29:26
