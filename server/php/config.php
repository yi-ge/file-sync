<?php
// DB config
// Hint: getenv() is not thread safe.
$DB_CONFIG = [
  'type' => 'mysql',
  'host' => getenv('MYSQL_HOST') ?: 'localhost',
  'database' => getenv('MYSQL_NAME') ?: 'file_sync',
  'port' => getenv('MYSQL_PORT') ?: 3306,
  'username' => getenv('MYSQL_USER') ?: 'root',
  'password' => getenv('MYSQL_PASS') ?: '',
  'charset' => 'utf8',
];

// error_reporting(0);
require_once 'libs/Medoo.php';

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