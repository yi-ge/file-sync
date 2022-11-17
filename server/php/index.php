<?php
require_once 'config.php';
require_once 'libs/Toro.php';

class HomeHandler
{
  function get()
  {
    echo "Hello, world";
  }

  function get_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

Toro::serve(
  array(
    "/" => "HomeHandler",
    "/device/add" => "HomeHandler",
    "/device/remove" => "HomeHandler",
    "/file/config" => "HomeHandler",
    "/file/check" => "HomeHandler",
    "/file/sync" => "HomeHandler",
  ),
  [
    'cors' => true
  ]
);