<?php
if (function_exists('shmop_open')) {
  require_once "../libs/SimpleBlock.php";
  $memory = new SimpleBlock(66);
  $memory->delete();
  echo "ok";
  return;
}

echo "fail";
