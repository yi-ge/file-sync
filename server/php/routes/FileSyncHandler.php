<?php

class FileSyncHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}