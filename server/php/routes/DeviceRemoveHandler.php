<?php

class DeviceRemoveHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => -99,
      "msg" => "Unknown error.",
      "result" => null
    ]);
  }
}