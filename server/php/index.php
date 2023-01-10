<?php
require_once 'config.php';
require_once 'libs/Toro.php';

if (function_exists('shmop_open')) {
    require_once 'libs/SimpleBlock.php';
}

// routes
require_once 'routes/HomeHandler.php';
require_once 'routes/DeviceAddHandler.php';
require_once 'routes/DeviceListHandler.php';
require_once 'routes/DeviceRemoveHandler.php';
require_once 'routes/FileConfigsHandler.php';
require_once 'routes/FileConfigHandler.php';
require_once 'routes/FileCheckHandler.php';
require_once 'routes/FileSyncHandler.php';

Toro::serve(
    array(
        "/" => "HomeHandler",
        "/device/add" => "DeviceAddHandler",
        "/device/list" => "DeviceListHandler",
        "/device/remove" => "DeviceRemoveHandler",
        "/file/configs" => "FileConfigsHandler",
        "/file/config" => "FileConfigHandler",
        "/file/check" => "FileCheckHandler",
        "/file/sync" => "FileSyncHandler",
    ),
    $server_options
);
