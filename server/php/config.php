<?php

$DB_CONFIG = [
  'type' => 'mysql',
  'host' => $_ENV['MYSQL_HOST'] ?: 'localhost',
  'database' => $_ENV['MYSQL_NAME'] ?: 'root',
  'port' => $_ENV['MYSQL_PORT'] ?: '3306',
  'username' => $_ENV['MYSQL_USER'] ?: 'root',
  'password' => $_ENV['MYSQL_PASS'],
  'charset' => 'utf8',
];