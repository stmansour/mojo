-- MySQL dump 10.13  Distrib 5.7.22, for osx10.12 (x86_64)
--
-- Host: localhost    Database: mojo
-- ------------------------------------------------------
-- Server version	5.7.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `DataUpdate`
--

DROP TABLE IF EXISTS `DataUpdate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DataUpdate` (
  `DUID` bigint(20) NOT NULL AUTO_INCREMENT,
  `GID` bigint(20) NOT NULL DEFAULT '0',
  `DtStart` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `DtStop` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`DUID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DataUpdate`
--

LOCK TABLES `DataUpdate` WRITE;
/*!40000 ALTER TABLE `DataUpdate` DISABLE KEYS */;
INSERT INTO `DataUpdate` VALUES (1,1,'2017-05-02 00:14:52','2017-05-02 00:15:02','2017-05-02 00:15:02',0);
/*!40000 ALTER TABLE `DataUpdate` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `EGroup`
--

DROP TABLE IF EXISTS `EGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `EGroup` (
  `GID` bigint(20) NOT NULL AUTO_INCREMENT,
  `GroupName` varchar(50) NOT NULL DEFAULT '',
  `GroupDescription` varchar(1000) NOT NULL DEFAULT '',
  `DtStart` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `DtStop` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`GID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `EGroup`
--

LOCK TABLES `EGroup` WRITE;
/*!40000 ALTER TABLE `EGroup` DISABLE KEYS */;
INSERT INTO `EGroup` VALUES (1,'FAA','Employees of the Federal Aviation Administration','2017-05-02 00:14:52','2017-05-02 00:15:02','2017-05-02 00:15:02',0),(2,'MojoTest','Steve-only test group','2020-08-16 07:52:42','2020-08-16 07:52:42','2020-08-16 07:52:42',0),(3,'AmazonTest','Steve + Amazon test accounts','2020-08-16 07:52:42','2020-08-16 07:52:42','2020-08-16 07:52:42',0),(4,'AccordTest','Steve + Amazon test + Accord accounts','2020-08-16 07:52:42','2020-08-16 07:52:42','2020-08-16 07:52:42',0),(5,'SteveJoe','Steve + Joe accounts','2020-08-16 07:52:42','2020-08-16 07:52:42','2020-08-16 07:52:42',0);
/*!40000 ALTER TABLE `EGroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `PGroup`
--

DROP TABLE IF EXISTS `PGroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PGroup` (
  `PID` bigint(20) NOT NULL DEFAULT '0',
  `GID` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `PGroup`
--

LOCK TABLES `PGroup` WRITE;
/*!40000 ALTER TABLE `PGroup` DISABLE KEYS */;
INSERT INTO `PGroup` VALUES (1,1,'2017-05-02 00:14:55',0),(2,1,'2017-05-02 00:14:56',0),(3,1,'2017-05-02 00:14:56',0),(4,1,'2017-05-02 00:14:56',0),(5,1,'2017-05-02 00:14:57',0),(6,1,'2017-05-02 00:14:57',0),(7,1,'2017-05-02 00:14:58',0),(8,1,'2017-05-02 00:14:58',0),(9,1,'2017-05-02 00:15:00',0),(10,1,'2017-05-02 00:15:02',0),(11,2,'2020-08-16 07:52:42',0),(12,2,'2020-08-16 07:52:42',0),(11,3,'2020-08-16 07:52:42',0),(12,3,'2020-08-16 07:52:42',0),(13,3,'2020-08-16 07:52:42',0),(14,3,'2020-08-16 07:52:42',0),(15,3,'2020-08-16 07:52:42',0),(11,4,'2020-08-16 07:52:42',0),(12,4,'2020-08-16 07:52:42',0),(13,4,'2020-08-16 07:52:42',0),(14,4,'2020-08-16 07:52:42',0),(15,4,'2020-08-16 07:52:42',0),(16,4,'2020-08-16 07:52:42',0),(17,4,'2020-08-16 07:52:42',0),(11,5,'2020-08-16 07:52:42',0),(12,5,'2020-08-16 07:52:42',0),(16,5,'2020-08-16 07:52:42',0);
/*!40000 ALTER TABLE `PGroup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `People`
--

DROP TABLE IF EXISTS `People`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `People` (
  `PID` bigint(20) NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(100) DEFAULT '',
  `MiddleName` varchar(100) DEFAULT '',
  `LastName` varchar(100) DEFAULT '',
  `PreferredName` varchar(100) DEFAULT '',
  `JobTitle` varchar(100) DEFAULT '',
  `OfficePhone` varchar(100) DEFAULT '',
  `OfficeFax` varchar(100) DEFAULT '',
  `Email1` varchar(50) DEFAULT '',
  `Email2` varchar(5) NOT NULL DEFAULT '',
  `MailAddress` varchar(50) DEFAULT '',
  `MailAddress2` varchar(50) DEFAULT '',
  `MailCity` varchar(100) DEFAULT '',
  `MailState` varchar(50) DEFAULT '',
  `MailPostalCode` varchar(50) DEFAULT '',
  `MailCountry` varchar(50) DEFAULT '',
  `RoomNumber` varchar(50) DEFAULT '',
  `MailStop` varchar(100) DEFAULT '',
  `Status` smallint(6) DEFAULT '0',
  `OptOutDate` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`PID`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `People`
--

LOCK TABLES `People` WRITE;
/*!40000 ALTER TABLE `People` DISABLE KEYS */;
INSERT INTO `People` VALUES (1,'Dave','C','Aakre','','ATSS','701-451-6805','','Dave.C.Aakre@faa.gov','','1801 23rd Ave N','','Fargo','ND','58102','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:55',0),(2,'Natalie','J','Aanerud','','ATCS','303-651-4241','','Natalie.J.Aanerud@faa.gov','','2211 17th Ave','','Longmont','CO','80501','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:56',0),(3,'John','','Aarhus','','Operations Supervisor - P80 Tracon','503-493-7580','','John.Aarhus@faa.gov','','7108 NE Airport Way','','Portland','OR','97218','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:56',0),(4,'Naga','CTR','Aarimanda','','N/A','N/A','','Naga.CTR.Aarimanda@faa.gov','','','','','','','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:56',0),(5,'Erik','','Aarness','','ATCS','651-463-5583','','Erik.Aarness@faa.gov','','512 Division St','','Farmington','MN','55024','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:57',0),(6,'Ryan','CTR','Aaron','','Helpdesk Specialist III','405-954-8747','','Ryan.CTR.Aaron@faa.gov','','6500 S MacArthur Blvd','','Oklahoma City','OK','73169','','BS06','',0,'0000-00-00 00:00:00','2017-05-02 00:14:57',0),(7,'Jeremy','','Aaronson','','Program Manager','202-267-7171','','Jeremy.Aaronson@faa.gov','','800 Independence Ave SW','','Washington','DC','20591','','','',0,'0000-00-00 00:00:00','2017-05-02 00:14:58',0),(8,'Lindsay','','Aaronson','','IdeaHub Operations Lead, Social Collaboration & Engagement Division','202-267-4016','','Lindsay.Aaronson@faa.gov','','800 Independence Ave SW','','Washington','DC','20591','','409W','',0,'0000-00-00 00:00:00','2017-05-02 00:14:58',0),(9,'Willie','','Aaron','','N/A','N/A','','Willie.Aaron@faa.gov','','1850 S Sigsbee St','','Indianapolis','IN','46241','','','',0,'0000-00-00 00:00:00','2017-05-02 00:15:00',0),(10,'John','','Aartman','','Front Line Manager','661-277-3843','','John.Aartman@faa.gov','','100 E Sparks Dr','','Edwards AFB','CA','93524','','','',0,'0000-00-00 00:00:00','2017-05-02 00:15:02',0),(11,'Steven','F','Mansour','','CTO, Accord Interests','323-512-0111 X305','','sman@accordinterests.com','','11719 Bee Cave Road','Suite 301','Austin','TX','78738','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(12,'Steve','','Mansour','','Recording Musician, Engineer, Producer','323-512-0111 X305','','sman@stevemansour.com','','2215 Wellington Drive','','Milpitas','CA','95035','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(13,'Bouncie','','McBounce','','Vagabond','123-456-7890','','bounce@simulator.amazonses.com','','123 Elm St','','Anytown','CA','90210','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(14,'Wendy','','Whiner','','Complainer','123-321-7890','','complaint@simulator.amazonses.com','','321 Elm St','','Anytown','CA','90210','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(15,'Stealthy','','McStealth','','Bad Guy','816-321-0123','','suppressionlist@simulator.amazonses.com','','700 Elm St','','Anytown','CA','90210','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(16,'Joe','G','Mansour','','Principal, Accord Interests','323-512-0111 X303','','jgm@accordinterests.com','','11719 Bee Cave Road','Suite 301','Austin','TX','78738','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0),(17,'Melissa','','Wheeler','','General Manager, Isola Bella','405.721.2194 x205','','mwheeler@accordinterests.com.com','','8309 NW 140th St','','Oklahoma City','OK','73142','USA','','',0,'0000-00-00 00:00:00','2020-08-16 07:52:42',0);
/*!40000 ALTER TABLE `People` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Query`
--

DROP TABLE IF EXISTS `Query`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Query` (
  `QID` bigint(20) NOT NULL AUTO_INCREMENT,
  `QueryName` varchar(50) DEFAULT '',
  `QueryDescr` varchar(1000) DEFAULT '',
  `QueryJSON` varchar(3000) DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`QID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Query`
--

LOCK TABLES `Query` WRITE;
/*!40000 ALTER TABLE `Query` DISABLE KEYS */;
INSERT INTO `Query` VALUES (1,'MojoTest','Steve-only test group','SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=2 WHERE People.Status=0','2020-08-16 07:52:42',0),(2,'AmazonTest','Steve + Amazon test accounts','SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=3 WHERE People.Status=0','2020-08-16 07:52:42',0),(3,'AccordTest','Steve + Amazon test + Accord accounts','SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=4 WHERE People.Status=0','2020-08-16 07:52:42',0),(4,'SteveJoe','Steve + Joe accounts','SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=5 WHERE People.Status=0','2020-08-16 07:52:42',0),(5,'FAA','The first 50 people in the FAA','SELECT People.* FROM People INNER JOIN PGroup ON PGroup.PID=People.PID AND PGroup.GID=1 WHERE People.Status=0','2020-08-16 07:52:42',0);
/*!40000 ALTER TABLE `Query` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-08-16  0:52:42
