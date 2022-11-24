<?php
require_once '../libs/Aes.php';

$data = "1234567891111111";
$aes = new Aes();
$key = "abcdabcdabcdabcdabcdabcdabcdabcd";
$encrypted = $aes->encrypt($data, $key);

echo $encrypted;