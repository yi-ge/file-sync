<?php

class FileCheckHandler
{
  function post_xhr($json)
  {
    global $database;

    if (
      !array_key_exists('email', $json) ||
      !array_key_exists('fileId', $json) ||
      !array_key_exists('sha256', $json)
    ) {
      echo json_encode([
        "status" => -1,
        "msg" => "Missing required parameters",
        "result" => null
      ]);
      return;
    }

    $emailSha1 = $json['email'];
    $fileId = $json['fileId'];
    $sha256 = $json['sha256'];

    $file = $database->select("file", "sha256", [
      "emailSha1" => $emailSha1,
      "fileId" => $fileId,
      "ORDER" => "updateAt DESC",
	    "LIMIT" => 1
    ]);

    if (!$file || $file->rowCount() != 1) {
      echo json_encode([
        "status" => 0,
        "msg" => "File not found",
        "result" => null
      ]);
      return;
    }

    if ($file["sha256"] != $sha256) {
      echo json_encode([
        "status" => 1
      ]);
    } else {
      echo json_encode([
        "status" => 2
      ]);
    }
  }
}