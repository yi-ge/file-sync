<?php
require_once '../libs/Aes.php';

$data = "1234567891111111";
$key = "abcdabcdabcdabcdabcdabcdabcdabcd";
$encrypted = Aes::encrypt($data, $key);

echo $encrypted . "\n";

$src = "123456789, abc, 中文，中文";
echo Aes::safetyBase64Encode($src);