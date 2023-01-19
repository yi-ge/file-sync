<?php
if (function_exists('mb_strlen')) {
  echo "support";
} else {
  echo "not support";
}