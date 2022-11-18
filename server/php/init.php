$sql = "CREATE TABLE `file_sync`.`device` (`id` BIGINT NOT NULL AUTO_INCREMENT , `email` VARCHAR(40) NOT NULL ,
`machineId`
VARCHAR(40) NOT NULL , `machineName` TEXT NOT NULL , `verify` VARCHAR(40) NOT NULL , `publicKey` TEXT NOT NULL ,
`privateKey` TEXT NOT NULL , `createdAt` DATETIME NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8 COLLATE
utf8_bin;";