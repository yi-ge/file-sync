<?php

class DeviceAddHandler
{
  function post_xhr($json)
  {
    global $database;
    if (
      !array_key_exists('email', $json) ||
      !array_key_exists('machineId', $json) ||
      !array_key_exists('machineName', $json) ||
      !array_key_exists('verify', $json) ||
      !array_key_exists('publicKey', $json) ||
      !array_key_exists('privateKey', $json)
    ) {
      echo json_encode([
        "status" => -1,
        "msg" => "Missing required parameters",
        "result" => null
      ]);
      return;
    }
    $email = $json['email'];
    $machineId = $json['machineId'];
    $machineName = $json['machineName'];
    $verify = $json['verify'];
    $publicKey = $json['publicKey'];
    $privateKey = $json['privateKey'];

    // Determine if the table exists
    $tableName = 'user';
    $row = $database->query("SHOW TABLES LIKE '" . $tableName . "'")->fetchAll();
    if ('1' != count($row)) { // Table does not exist
      // user table
      $database->query("CREATE TABLE IF NOT EXISTS `" . $tableName . "` (
        `id` BIGINT NOT NULL AUTO_INCREMENT,
        `email` VARCHAR(40) NOT NULL,
        `verify` VARCHAR(40) NOT NULL,
        `publicKey` TEXT NOT NULL,
        `privateKey` TEXT NOT NULL,
        `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`), UNIQUE `email_keys` (`email`)
        ) ENGINE = InnoDB")->fetchAll();

      // device table
      $database->query("CREATE TABLE IF NOT EXISTS `device` (
        `id` BIGINT NOT NULL AUTO_INCREMENT,
        `email` VARCHAR(40) NOT NULL,
        `machineId` VARCHAR(40) NOT NULL,
        `machineName` TEXT NOT NULL,
        `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`), UNIQUE `machineId_keys` (`machineId`)
        ) ENGINE = InnoDB")->fetchAll();

      // config table
      $database->query("CREATE TABLE IF NOT EXISTS `config` (
        `id` BIGINT NOT NULL AUTO_INCREMENT,
        `email` VARCHAR(40) NOT NULL,
        `machineId` VARCHAR(40) NOT NULL,
        `fileId` VARCHAR(40) NOT NULL,
        `path` TEXT NOT NULL,
        `attribute` TEXT NULL,
        `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`)
        ) ENGINE = InnoDB")->fetchAll();

      // file table
      $database->query("CREATE TABLE IF NOT EXISTS `file` (
        `id` BIGINT NOT NULL AUTO_INCREMENT,
        `email` VARCHAR(40) NOT NULL,
        `fileId` VARCHAR(40) NOT NULL,
        `fileName` TEXT NULL,
        `content` LONGTEXT NULL,
        `sha256` VARCHAR(64) NOT NULL,
        `fromMachineId` VARCHAR(40) NOT NULL,
        `updateAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`), UNIQUE `fileId_keys` (`fileId`)
        ) ENGINE = InnoDB")->fetchAll();

      // log table
      $database->query("CREATE TABLE IF NOT EXISTS `log` (
        `id` BIGINT NOT NULL AUTO_INCREMENT,
        `email` VARCHAR(40) NOT NULL,
        `fromMachineId` VARCHAR(40) NOT NULL,
        `action` VARCHAR(40) NOT NULL,
        `content` TEXT NULL,
        `createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`)
        ) ENGINE = InnoDB")->fetchAll();
    }

    if (
      !$database->has("user", [
        "email" => $email
      ])
    ) { // new user
      $database->insert("user", [
        "email" => $email,
        "verify" => $verify,
        "publicKey" => $publicKey,
        "privateKey" => $privateKey
      ]);

      $database->insert("device", [
        "email" => $email,
        "machineId" => $machineId,
        "machineName" => $machineName
      ]);

      $database->insert("log", [
        "email" => $email,
        "machineId" => $machineId,
        "action" => "DeviceAdd",
        "content" => json_encode([
          "newUser" => true,
          "email" => $email,
          "verify" => $verify,
          "publicKey" => $publicKey,
          "privateKey" => $privateKey,
          "machineId" => $machineId,
          "machineName" => $machineName,
        ])
      ]);

      $aes = new Aes();
      // TODO: change key
      $key = "abcdabcdabcdabcdabcdabcdabcdabcd";
      $encrypted = $aes->encrypt($publicKey, $key);

      echo json_encode([
        "status" => 1,
        "msg" => "New user added",
        "result" => [
          "publicKey" => $encrypted,
          "privateKey" => $privateKey
        ]
      ]);
    } else {
      if (
        $database->has("device", [
          "email" => $email,
          "machineId" => $machineId
        ])
      ) {
        echo json_encode([
          "status" => -2,
          "msg" => "Device already exists",
          "result" => null
        ]);
        return;
      }

      $user = $database->get("user", "*", [
        "email" => $email
      ]);

      if (count($user)) {
        $activeVerify = $user['verify'];
        $activePublicKey = $user['publicKey'];
        $activePrivateKey = $user['privateKey'];

        if ($activeVerify != $verify) {
          echo json_encode([
            "status" => -3,
            "msg" => "Verification Rejected",
            "result" => null
          ]);
          return;
        }

        $last_device_id = $database->insert("device", [
          "email" => $email,
          "machineId" => $machineId,
          "machineName" => $machineName
        ]);

        if ($last_device_id > 0) {
          $database->insert("log", [
            "email" => $email,
            "machineId" => $machineId,
            "action" => "DeviceAdd",
            "content" => json_encode([
              "newUser" => false,
              "email" => $email,
              "machineId" => $machineId,
              "machineName" => $machineName,
            ])
          ]);

          echo json_encode([
            "status" => 2,
            "msg" => "Device added",
            "result" => [
              "publicKey" => $activePublicKey,
              "privateKey" => $activePrivateKey
            ]
          ]);
        }
      }

      echo json_encode([
        "status" => -99,
        "msg" => "Unknown error.",
        "result" => null
      ]);
    }
  }
}