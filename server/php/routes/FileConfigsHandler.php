<?php

class FileConfigsHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}