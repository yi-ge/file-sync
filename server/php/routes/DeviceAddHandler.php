<?php

class DeviceAddHandler
{
  function post_xhr($json)
  {
    var_dump($json);
    $email = $json['email'];
    $machineId = $json['machine_id'];
    $machineName = $json['machineName'];
    $verify = $json['verify'];
    $publicKey = $json['publicKey'];
    $privateKey = $json['privateKey'];

    if (!$email || !$machineId || !$machineName || !$verify || !$publicKey || !$privateKey) {
      echo json_encode([
        "status" => -1,
        "msg" => "Missing required parameters"
      ]);
    } else {
      echo json_encode([
        "status" => 1
      ]);
    }
  }
}