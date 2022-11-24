<?php

class DeviceListHandler
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

    $email = $json['email'];
    // $machineId = $json['machineId'];
    // $timestamp = $json['timestamp'];
    $token = $json['token'];

    $user = $database->get("user", "publicKey", [
      "email" => $email
    ]);

    if (count($user) == 0) {
      echo json_encode([
        "status" => -2,
        "msg" => "Invalid email",
        "result" => null
      ]);
      return;
    }

    $publicKey = $user["publicKey"];
    $publicKeyId = openssl_pkey_get_public($publicKey);
    unset($json["token"]);
    sort($json);
    $toSign = json_encode($json);

    if (openssl_verify($toSign, base64_decode($token), $publicKeyId) != 1) {
      echo json_encode([
        "status" => -3,
        "msg" => "Invalid token",
        "result" => null
      ]);
      return;
    }

    $machineList = $database->select("machine", "*", [
      "email" => $email
    ]);

    echo json_encode([
      "status" => 1,
      "msg" => "OK",
      "result" => $machineList
    ]);
  }
}