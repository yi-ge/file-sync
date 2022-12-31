<?php

class FileConfigsHandler
{
  function post_xhr($json)
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

    $fileList = $database->select("file", "*", [
      "email" => $user['email']
    ]);

    echo json_encode([
      "status" => 1,
      "msg" => "OK",
      "result" => $fileList
    ]);
  }
}