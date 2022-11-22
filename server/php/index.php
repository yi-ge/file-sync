<?php
require_once 'config.php';
require_once 'libs/Toro.php';

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

class DeviceAddHandler
{
  function post_xhr($json)
  {
    var_dump($json);
    // $email =
    echo json_encode([
      "status" => 1
    ]);
  }
}

class DeviceRemoveHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

class FileConfigsHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

class FileConfigHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

class FileCheckHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

class FileSyncHandler
{
  function post_xhr()
  {
    echo json_encode([
      "status" => 1
    ]);
  }
}

Toro::serve(
  array(
    "/" => "HomeHandler",
    "/device/add" => "DeviceAddHandler",
    "/device/remove" => "DeviceRemoveHandler",
    "/file/configs" => "FileConfigsHandler",
    "/file/config" => "FileConfigHandler",
    "/file/check" => "FileCheckHandler",
    "/file/sync" => "FileSyncHandler",
  ),
  $server_options
);