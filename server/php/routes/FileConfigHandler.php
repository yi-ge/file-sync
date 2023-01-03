<?php

class FileConfigHandler
{
  function post_xhr($json)
  {
    global $database;

    if (
      !array_key_exists('email', $json) ||
      !array_key_exists('machineId', $json) ||
      !array_key_exists('timestamp', $json) ||
      !array_key_exists('fileId', $json) ||
      !array_key_exists('action', $json) ||
      !array_key_exists('actionMachineId', $json) ||
      !array_key_exists('attribute', $json) ||
      !array_key_exists('token', $json)
    ) {
      echo json_encode([
        "status" => -1,
        "msg" => "Missing required parameters",
        "result" => null
      ]);
      return;
    }

    $emailSha1 = $json['email'];
    $token = $json['token'];

    $user = $database->get("user", "*", [
      "emailSha1" => $emailSha1
    ]);

    if (!$user) {
      echo json_encode([
        "status" => -2,
        "msg" => "Invalid email",
        "result" => null
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
        "result" => null
      ]);
      return;
    }

    if ($json['action'] == 'add') { // Add
      if (
        !array_key_exists('path', $json) ||
        !array_key_exists('fileName', $json)
      ) {
        echo json_encode([
          "status" => -4,
          "Missing required parameters",
          "result" => null
        ]);
        return;
      }

      $last_config_id = $database->insert("config", [
        "email" => $user['email'],
        "machineId" => $json['actionMachineId'],
        "fileId" => $json['fileId'],
        "fileName" => $json['fileName'],
        "path" => $json['path'],
        "attribute" => $json['attribute']
      ]);

      echo json_encode([
        "status" => 1,
        "msg" => "OK",
        "result" => [
          "lastId" => $last_config_id,
          "fileId" => $json['fileId']
        ]
      ]);
      return;
    } else if ($json['action'] == 'remove') { // Remove
      $data = $database->update("config", [
        "deletedAt" => time()
      ], [
        "email" => $user['email'],
        "fileId" => $json['fileId'],
        "machineId" => $json['actionMachineId'],
        "deletedAt" => null
      ]);

      if ($data->rowCount() == 1) {
        echo json_encode([
          "status" => 1,
          "msg" => "OK",
          "result" => null
        ]);
        return;
      }

      echo json_encode([
        "status" =>-4,
        "msg" => "Delete fail.",
        "result" => null
      ]);
      return;
    }

    echo json_encode([
      "status" => -2,
      "msg" => "Invalid action",
      "result" => null
    ]);
  }
}