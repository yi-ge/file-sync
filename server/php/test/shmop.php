<?php
require_once '../libs/SimpleBlock.php';

$results = [
    '123',
    '456',
    '789',
];

$data = json_encode($results);

$memory = new SimpleBlock(66);
$memory->write($data);
$stored_array = json_decode($memory->read(), true);

print_r($stored_array);
