<?php
// Disable error reporting
// error_reporting(0);

// Report runtime errors
// error_reporting(E_ERROR | E_WARNING | E_PARSE);

require_once 'libs/Medoo.php';

use Medoo\Medoo;

// Database configuration
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

// Server configuration
$server_options = [
  'cors' => true
];

// If you are using PHP version 7.5 or higher, it is recommended to switch to a higher version of Medoo
$database = new Medoo([
  'database_type' => $DB_CONFIG['type'],
  'database_name' => $DB_CONFIG['database'],
  'server' => $DB_CONFIG['host'],
  'username' => $DB_CONFIG['username'],
  'password' => $DB_CONFIG['password'],
  'charset' => $DB_CONFIG['charset'],
  'port' => $DB_CONFIG['port'],
]);