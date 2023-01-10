<?php
// $shmid = shmop_open(66, "w", 0, 0);
$shmid = shmop_open(66, "c", 0755, 40);
$data = shmop_read($shmid, 0, 40);
echo $data;
