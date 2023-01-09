<?php
header('X-Accel-Buffering: no');
header('Content-Type: text/event-stream');
header('Cache-Control: no-cache');

set_time_limit(0); //防止超时
ob_end_clean(); //清空(擦除)缓冲区并关闭输出缓冲
ob_implicit_flush(1); //这个函数强制每当有输出的时候，即刻把输出发送到浏览器。这样就不需要每次输出(echo)后，都用flush()来发送到浏览器了

/**
 * @param int $time 要休眠的秒数。要休眠 5 秒，就传入 5
 */
function delay($time) {
  /*
  * 检测 PHP 内置的 sleep() 函数是否被禁用
  */
  if (function_exists("sleep")) {
      // PHP 内置 sleep() 函数未被禁用，优先使用
      sleep($time);
  } else {
      // PHP 内置 sleep() 函数已被禁用，曲线救国
      // 由于这种方式存在极小的误差，所以额外增加1秒以保证可靠性
      $targetTime = time() + $time + 1;

      while (true) {
          if (time() == $targetTime) {
              break;
          }
      }
  }
}

// 发送消息
$data = 0;
while(true)
{
    $data++;
    $c = "event:message" . PHP_EOL; //定义事件
    $c .= "data: " . $data . PHP_EOL; //推送内容
    echo $c . PHP_EOL;
    delay(1);
}
