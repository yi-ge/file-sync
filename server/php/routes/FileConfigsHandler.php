<?php

class FileConfigsHandler
{
    public function post_xhr($json)
    {
        global $database;

        if (
            !array_key_exists('email', $json) ||
            !array_key_exists('machineId', $json) ||
            !array_key_exists('timestamp', $json) ||
            !array_key_exists('token', $json)
        ) {
            echo json_encode([
                "status" => -1,
                "msg" => "Missing required parameters",
                "result" => null,
            ]);
            return;
        }

        $emailSha1 = $json['email'];
        $token = $json['token'];

        $user = $database->get("user", "*", [
            "emailSha1" => $emailSha1,
        ]);

        if (!$user) {
            echo json_encode([
                "status" => -2,
                "msg" => "Invalid email",
                "result" => null,
            ]);
            return;
        }

        unset($json["token"]);
        ksort($json);
        $sign = '';
        foreach ($json as $k => $v) {
            $sign .= $k . "=" . $v . "&";
        }
        $sign = $sign . $user["verify"];

        $token = Aes::safetyBase64Decode($token);

        $publicKey = $user["publicKey"];
        $publicKeyId = openssl_pkey_get_public($publicKey);

        if (openssl_verify($sign, $token, $publicKeyId, OPENSSL_ALGO_SHA1) != 1) {
            echo json_encode([
                "status" => -3,
                "msg" => "Invalid token",
                "result" => null,
            ]);
            return;
        }

        // $configList = $database->select("config", "*", [
        //   "email" => $user['email'],
        //   "deletedAt" => null,
        // ]);

        $configList = $database->query("
      SELECT
        config.id as id,
        config.fileName as fileName,
        config.fileId as fileId,
        file.updateAt as updateAt,
        config.machineId as machineId,
        device.machineName as machineName,
        config.path as path,
        config.attribute as attribute,
        config.createdAt as createdAt
      FROM
        config
      LEFT JOIN file on file.id = (
        SELECT f.id FROM file AS f
        WHERE config.fileId = f.fileId
        ORDER BY f.updateAt DESC
        LIMIT 1
      )
      LEFT JOIN device on config.machineId = device.machineId
      WHERE
        <config.email> = :email AND
        deletedAt is NULL
      ORDER BY
        config.fileId,
        config.createdAt,
        config.id", [
            ":email" => $user['email'],
        ])->fetchAll(PDO::FETCH_ASSOC);

        echo json_encode([
            "status" => 1,
            "msg" => "OK",
            "result" => $configList,
        ]);
    }
}
