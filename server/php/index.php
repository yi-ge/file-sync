<?php

require_once 'Toro.php';

class HelloHandler
{
  function get()
  {
    echo "Hello, world";
  }
}

Toro::serve(
  array(
    "/" => "HelloHandler",
  )
);