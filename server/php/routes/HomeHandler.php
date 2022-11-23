<?php

class HomeHandler
{
  function get()
  {
    header('HTTP/1.1 403 Unauthorized');
    echo "403 Unauthorized";
  }

  function get_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}