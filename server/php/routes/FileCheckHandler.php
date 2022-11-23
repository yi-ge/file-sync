<?php

class FileCheckHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}