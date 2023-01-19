<?php
if (function_exists('shmop_open')) {
    require_once 'libs/SimpleBlock.php';
}

error_reporting(E_ALL ^ E_NOTICE);
header('X-Accel-Buffering: no');
header('Content-Type: text/event-stream');
header('Cache-Control: no-cache');

set_time_limit(0); // Prevent timeouts
ob_end_clean(); // Empty (erase) the buffer and close the output buffer
ob_implicit_flush(1); // This function forces the output to be sent to the browser as soon as it is available. This eliminates the need to use flush() to send the output to the browser after each output (echo)
$datetimeFormat = 'Y-m-d H:i:s';

/**
 * delay
 * @param int $time The number of seconds to hibernate. To hibernate for 5 seconds, pass in 5
 */
function delay($time)
{
    // Check if PHP's built-in sleep() function is disabled
    if (function_exists("sleep")) {
        // The PHP built-in sleep() function is not disabled and is preferred
        sleep($time);
    } else {
        // PHP's built-in sleep() function is disabled, save the day
        // Since there is a very small error in this way, an extra 1 second is added to ensure reliability
        $targetTime = time() + $time + 1;

        while (true) {
            if (time() == $targetTime) {
                break;
            }
        }
    }
}

if (!isset($_GET['email'])) {
    echo "Invalid request";
    exit;
}
$emailSha1 = $_GET['email'];
$timestamp = $_GET['timestamp'];
$date = new DateTime('now', new DateTimeZone('Asia/Shanghai'));
if (!empty($timestamp)) {
    $date->setTimestamp(intval($timestamp / 1000));
}
$lastData = "";

$connectedMsg = "event: connected" . PHP_EOL; // Define Event
$connectedMsg .= "data: 1" . PHP_EOL; // Push content
echo $connectedMsg . PHP_EOL;

if (function_exists('shmop_open')) {
    $time = 60;
    while (true) {
        $memory = new SimpleBlock(66);
        $data = "";
        if ($memory->exists(66)) {
            $data = $memory->read();
        }
        if ($data && $data != $lastData) {
            $lastData = $data;
            $data = json_decode($data, true);
            $c = "event: file" . PHP_EOL; // Define Event
            $c .= "data: " . join(",", $data) . PHP_EOL; // Push content
            echo $c . PHP_EOL;
        } else {
            if ($time > 59) {
                $c = "event: heartbeat" . PHP_EOL; // Define Event
                $c .= "data: 1" . PHP_EOL; // Push content
                echo $c . PHP_EOL;
                $time = 0;
            } else {
                $time++;
            }
        }

        delay(1);
    }
} else {
    // Check if PHP's built-in shmop_open() function is disabled
    require_once 'config.php';

    $time = 11;
    // Send message
    while (true) {
        $files = $database->select("file", "fileId", [
            "emailSha1" => $emailSha1,
            "updateAt[>=]" => $date->format($datetimeFormat),
            "ORDER" => ["updateAt" => "DESC"],
        ]);

        $date = new DateTime('now', new DateTimeZone('Asia/Shanghai'));

        if ($files && sizeof($files) >= 1) {
            $c = "event: file" . PHP_EOL; // Define Event
            $c .= "data: " . join(",", $files) . PHP_EOL; // Push content
            echo $c . PHP_EOL;
        } else {
            if ($time > 10) {
                $c = "event: heartbeat" . PHP_EOL; // Define Event
                $c .= "data: 1" . PHP_EOL; // Push content
                echo $c . PHP_EOL;
                $time = 0;
            } else {
                $time++;
            }
        }

        delay(5);
    }
}
