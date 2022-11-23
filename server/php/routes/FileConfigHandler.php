<?php

class FileConfigHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}