CREATE SCHEMA IF NOT EXISTS `development` DEFAULT CHARACTER SET utf8mb4 ;
CREATE SCHEMA IF NOT EXISTS `test` DEFAULT CHARACTER SET utf8mb4 ;

GRANT ALL PRIVILEGES ON `test`.* TO 'user'@'%';
FLUSH PRIVILEGES;
SET CHARSET utf8mb4;
--
-- -- -----------------------------------------------------
-- -- Table `clean_architecture`.`users`
-- -- -----------------------------------------------------
-- CREATE TABLE IF NOT EXISTS `clean_architecture`.`user` (
--     `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
--     `name` VARCHAR(64) NOT NULL COMMENT 'ユーザ名',
--     PRIMARY KEY (`id`));
--
-- USE clean_architecture;
--
-- SET NAMES utf8mb4;
--
-- INSERT INTO `user` (`id`,`name`) VALUES ("1","name1");
-- INSERT INTO `user` (`id`,`name`) VALUES ("2","name2");
-- INSERT INTO `user` (`id`,`name`) VALUES ("3","name3");
-- INSERT INTO `user` (`id`,`name`) VALUES ("4","name4");
-- INSERT INTO `user` (`id`,`name`) VALUES ("5","name5");
-- INSERT INTO `user` (`id`,`name`) VALUES ("6","name6");
-- INSERT INTO `user` (`id`,`name`) VALUES ("7","name7");
-- INSERT INTO `user` (`id`,`name`) VALUES ("8","name8");
-- INSERT INTO `user` (`id`,`name`) VALUES ("9","name9");
-- INSERT INTO `user` (`id`,`name`) VALUES ("10","name10");