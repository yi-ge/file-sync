<?php
error_reporting(0);
require_once 'config.php';
require_once 'Medoo.php';

use Medoo\Medoo;

$database = new Medoo([
  'database_type' => $DB_CONFIG['type'],
  'database_name' => $DB_CONFIG['database'],
  'server' => $DB_CONFIG['host'],
  'username' => $DB_CONFIG['username'],
  'password' => $DB_CONFIG['password'],
  'charset' => $DB_CONFIG['charset'],
  'port' => $DB_CONFIG['port'],
]);

header('Access-Control-Allow-Headers: Authorization, DNT, User-Agent, Keep-Alive, Origin, X-Requested-With, Content-Type, Accept, x-clientid');
header('Access-Control-Allow-Methods: PUT, POST, GET, DELETE, OPTIONS');
header('Access-Control-Allow-Origin: *');
header('Content-Type: application/json; charset=utf-8');

echo json_encode(['err' => '非法访问']);