<?php
require_once 'config.php';
require_once 'libs/Toro.php';

class HelloHandler
{
  function get()
  {
    echo "Hello, world";
  }

  function get_xhr()
  {
    echo json_encode([
      "status" => 200
    ]);
  }
}

Toro::serve(
  array(
    "/" => "HelloHandler",
  ),
  [
    'cors' => true
  ]
);