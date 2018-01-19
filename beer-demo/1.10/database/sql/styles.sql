-- MySQL dump 10.11
--
-- Host: internal-db.s33724.gridserver.com    Database: db33724_openbeer_db
-- ------------------------------------------------------
-- Server version	5.0.32-Debian_7etch5~bpo31+1

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
-- Table structure for table `styles`
--

DROP TABLE IF EXISTS `styles`;
CREATE TABLE `styles` (
  `id` int(11) NOT NULL auto_increment,
  `cat_id` int(11) NOT NULL default '0',
  `style_name` varchar(255) character set utf8 collate utf8_unicode_ci NOT NULL default '',
  `last_mod` datetime NOT NULL default '0000-00-00 00:00:00',
  PRIMARY KEY  (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=142 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `styles`
--

LOCK TABLES `styles` WRITE;
/*!40000 ALTER TABLE `styles` DISABLE KEYS */;
INSERT INTO `styles` VALUES (1,1,'Classic English-Style Pale Ale','2010-10-24 13:53:31'),(2,1,'English-Style India Pale Ale','2010-06-15 19:14:38'),(3,1,'Ordinary Bitter','2010-06-15 19:14:54'),(4,1,'Special Bitter or Best Bitter','2010-06-15 19:15:02'),(5,1,'Extra Special Bitter','2010-06-15 19:15:09'),(6,1,'English-Style Summer Ale','2010-06-15 19:15:20'),(7,1,'Scottish-Style Light Ale','2010-06-15 19:15:28'),(8,1,'Scottish-Style Heavy Ale','2010-06-15 19:15:37'),(9,1,'Scottish-Style Export Ale','2010-06-15 19:15:45'),(10,1,'English-Style Pale Mild Ale','2010-06-15 19:16:17'),(11,1,'English-Style Dark Mild Ale','2010-06-15 19:16:27'),(12,1,'English-Style Brown Ale','2010-06-15 19:16:34'),(13,1,'Old Ale','2010-06-15 19:16:40'),(14,1,'Strong Ale','2010-06-15 19:16:46'),(15,1,'Scotch Ale','2010-06-15 19:16:54'),(16,1,'British-Style Imperial Stout','2010-06-15 19:17:07'),(17,1,'British-Style Barley Wine Ale','2010-06-15 19:17:16'),(18,1,'Robust Porter','2010-06-15 19:17:37'),(19,1,'Brown Porter','2010-06-15 19:17:43'),(20,1,'Sweet Stout','2010-06-15 19:17:50'),(21,1,'Oatmeal Stout','2010-06-15 19:17:56'),(22,2,'Irish-Style Red Ale','2010-06-15 19:21:04'),(23,2,'Classic Irish-Style Dry Stout','2010-06-15 19:21:07'),(24,2,'Foreign (Export)-Style Stout','2010-06-15 19:21:14'),(25,2,'Porter','2010-06-15 19:21:39'),(26,3,'American-Style Pale Ale','2010-06-15 19:22:32'),(27,3,'Fresh Hop Ale','2010-06-15 19:22:41'),(28,3,'Pale American-Belgo-Style Ale','2010-06-15 19:22:51'),(29,3,'Dark American-Belgo-Style Ale','2010-06-15 19:23:00'),(30,3,'American-Style Strong Pale Ale','2010-06-15 19:23:07'),(31,3,'American-Style India Pale Ale','2010-06-15 19:23:15'),(32,3,'Imperial or Double India Pale Ale','2010-06-15 19:23:23'),(33,3,'American-Style Amber/Red Ale','2010-06-15 19:23:30'),(34,3,'Imperial or Double Red Ale','2010-06-15 19:23:52'),(35,3,'American-Style Barley Wine Ale','2010-06-15 19:23:59'),(36,3,'American-Style Wheat Wine Ale','2010-06-15 19:24:08'),(37,3,'Golden or Blonde Ale','2010-06-15 19:24:15'),(38,3,'American-Style Brown Ale','2010-06-15 19:24:23'),(39,3,'Smoke Porter','2010-06-15 19:24:30'),(40,3,'American-Style Sour Ale','2010-06-15 19:24:58'),(41,3,'American-Style India Black Ale','2010-06-15 19:25:05'),(42,3,'American-Style Stout','2010-06-15 19:25:11'),(43,3,'American-Style Imperial Stout','2010-06-15 19:25:18'),(44,3,'Specialty Stouts','2010-06-15 19:25:26'),(45,3,'American-Style Imperial Porter','2010-06-15 19:25:33'),(46,3,'Porter','2010-06-15 19:25:53'),(47,4,'German-Style Kolsch','2010-06-15 20:43:02'),(48,4,'Berliner-Style Weisse','2010-06-15 19:29:33'),(49,4,'Leipzig-Style Gose','2010-06-15 19:29:40'),(50,4,'South German-Style Hefeweizen','2010-06-15 19:29:47'),(51,4,'South German-Style Kristal Weizen','2010-06-15 19:29:55'),(52,4,'German-Style Leichtes Weizen','2010-06-15 19:30:01'),(53,4,'South German-Style Bernsteinfarbenes Weizen','2010-06-15 19:30:18'),(54,4,'South German-Style Dunkel Weizen','2010-06-15 19:30:25'),(55,4,'South German-Style Weizenbock','2010-06-15 19:30:33'),(56,4,'Bamberg-Style Weiss Rauchbier','2010-06-15 19:30:47'),(57,4,'German-Style Brown Ale/Altbier','2010-07-11 16:08:56'),(58,4,'Kellerbier - Ale','2010-06-15 19:38:57'),(59,5,'Belgian-Style Flanders/Oud Bruin','2010-06-15 19:32:09'),(60,5,'Belgian-Style Dubbel','2010-06-15 19:32:15'),(61,5,'Belgian-Style Tripel','2010-06-15 19:32:23'),(62,5,'Belgian-Style Quadrupel','2010-06-15 19:32:31'),(63,5,'Belgian-Style Blonde Ale','2010-06-15 19:32:39'),(64,5,'Belgian-Style Pale Ale','2010-06-15 19:32:47'),(65,5,'Belgian-Style Pale Strong Ale','2010-06-15 19:32:55'),(66,5,'Belgian-Style Dark Strong Ale','2010-06-15 19:33:01'),(67,5,'Belgian-Style White','2010-06-15 19:33:16'),(68,5,'Belgian-Style Lambic','2010-06-15 19:33:43'),(69,5,'Belgian-Style Gueuze Lambic','2010-06-15 19:33:52'),(70,5,'Belgian-Style Fruit Lambic','2010-06-15 19:34:01'),(71,5,'Belgian-Style Table Beer','2010-06-15 19:34:08'),(72,5,'Other Belgian-Style Ales','2010-06-15 19:34:15'),(73,5,'French-Style Biere de Garde','2010-07-11 15:55:26'),(74,5,'French & Belgian-Style Saison','2010-06-15 19:34:30'),(75,6,'International-Style Pale Ale','2010-06-15 19:34:51'),(76,6,'Australasian-Style Pale Ale','2010-06-15 19:35:01'),(77,7,'German-Style Pilsener','2010-06-15 19:35:42'),(78,7,'Bohemian-Style Pilsener','2010-06-15 19:35:54'),(79,7,'European Low-Alcohol Lager','2010-06-15 19:36:02'),(80,7,'Munchner-Style Helles','2010-06-15 19:39:29'),(81,7,'Dortmunder/European-Style Export','2010-06-15 19:36:27'),(82,7,'Vienna-Style Lager','2010-06-15 19:36:34'),(83,7,'German-Style Marzen','2010-07-08 17:59:32'),(84,7,'German-Style Oktoberfest','2010-06-15 19:37:01'),(85,7,'European-Style Dark','2010-06-15 19:37:09'),(86,7,'German-Style Schwarzbier','2010-06-15 19:37:36'),(87,7,'Bamberg-Style Marzen','2010-07-08 17:59:39'),(88,7,'Bamberg-Style Helles Rauchbier','2010-06-15 19:37:50'),(89,7,'Bamberg-Style Bock Rauchbier','2010-06-15 19:37:58'),(90,7,'Traditional German-Style Bock','2010-06-15 19:38:06'),(91,7,'German-Style Heller Bock/Maibock','2010-06-15 19:38:15'),(92,7,'German-Style Doppelbock','2010-06-15 19:38:22'),(93,7,'German-Style Eisbock','2010-06-15 19:38:30'),(94,7,'Kellerbier - Lager','2010-06-15 19:38:49'),(95,8,'American-Style Lager','2010-06-15 19:40:43'),(96,8,'American-Style Light Lager','2010-06-15 19:40:46'),(97,8,'American-Style Low-Carb Light Lager','2010-06-15 19:41:05'),(98,8,'American-Style Amber Lager','2010-06-15 19:41:18'),(99,8,'American-Style Premium Lager','2010-06-15 19:41:26'),(100,8,'American-Style Pilsener','2010-06-15 19:41:33'),(101,8,'American-Style Ice Lager','2010-06-15 19:41:43'),(102,8,'American-Style Malt Liquor','2010-06-15 19:41:51'),(103,8,'American-Style Amber Lager','2010-06-15 19:41:58'),(104,8,'American-Style Marzen/Oktoberfest','2010-06-15 19:44:27'),(105,8,'American-Style Dark Lager','2010-06-15 19:42:28'),(106,9,'Baltic-Style Porter','2010-06-15 19:42:40'),(107,9,'Australasian-Style Light Lager','2010-06-15 19:44:09'),(108,9,'Latin American-Style Light Lager','2010-06-15 19:43:29'),(109,9,'Tropical-Style Light Lager','2010-06-15 19:43:37'),(110,10,'International-Style Pilsener','2010-06-15 19:43:49'),(111,10,'Dry Lager','2010-06-15 19:43:57'),(112,11,'Session Beer','2010-06-15 19:45:11'),(113,11,'American-Style Cream Ale or Lager','2010-06-15 19:45:19'),(114,11,'California Common Beer','2010-06-15 19:45:25'),(115,11,'Japanese Sake-Yeast Beer','2010-06-15 19:45:34'),(116,11,'Light American Wheat Ale or Lager','2010-06-15 19:45:42'),(117,11,'Fruit Wheat Ale or Lager','2010-06-15 19:45:58'),(118,11,'Dark American Wheat Ale or Lager','2010-06-15 19:46:06'),(119,11,'American Rye Ale or Lager','2010-06-15 19:46:14'),(120,11,'German-Style Rye Ale','2010-06-15 19:46:24'),(121,11,'Fruit Beer','2010-06-15 19:46:33'),(122,11,'Field Beer','2010-06-15 19:46:41'),(123,11,'Pumpkin Beer','2010-06-15 19:47:18'),(124,11,'Chocolate/Cocoa-Flavored Beer','2010-06-15 19:47:27'),(125,11,'Coffee-Flavored Beer','2010-06-15 19:47:35'),(126,11,'Herb and Spice Beer','2010-06-15 19:47:40'),(127,11,'Specialty Beer','2010-06-15 19:47:48'),(128,11,'Specialty Honey Lager or Ale','2010-06-15 19:47:54'),(129,11,'Gluten-Free Beer','2010-06-15 19:48:01'),(130,11,'Smoke Beer','2010-06-15 19:48:10'),(131,11,'Experimental Beer','2010-06-15 19:48:17'),(132,11,'Out of Category','2010-06-15 19:48:27'),(133,11,'Wood- and Barrel-Aged Beer','2010-06-15 19:48:39'),(134,11,'Wood- and Barrel-Aged Pale to Amber Beer','2010-06-15 19:48:53'),(135,11,'Wood- and Barrel-Aged Dark Beer','2010-06-15 19:49:01'),(136,11,'Wood- and Barrel-Aged Strong Beer','2010-06-15 19:49:08'),(137,11,'Wood- and Barrel-Aged Sour Beer','2010-06-15 19:49:18'),(138,11,'Aged Beer','2010-06-15 19:49:25'),(139,11,'Other Strong Ale or Lager','2010-06-15 19:49:32'),(140,11,'Non-Alcoholic Beer','2010-06-15 19:49:47'),(141,11,'Winter Warmer','2010-07-11 16:43:25');
/*!40000 ALTER TABLE `styles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2011-10-24  1:00:49
