<?php
class Aes
{
    protected static $method = 'AES-256-CBC';

    public static function encrypt($data, $key)
    {
        $ivLen = openssl_cipher_iv_length(static::$method);
        $iv = openssl_random_pseudo_bytes($ivLen);
        $text = openssl_encrypt($data, static::$method, $key, OPENSSL_RAW_DATA, $iv);
        return self::safetyBase64Encode($iv . $text);
    }

    public static function decrypt($text, $key)
    {
        $cipherText = self::safetyBase64Decode($text);
        $ivLen = openssl_cipher_iv_length(static::$method);
        $iv = substr($cipherText, 0, $ivLen);
        $cipherText = substr($cipherText, $ivLen);
        $data = openssl_decrypt($cipherText, static::$method, $key, OPENSSL_RAW_DATA, $iv);
        return $data;
    }

    public static function safetyBase64Encode($text)
    {
        $text = base64_encode($text);
        $text = str_replace(['+', '/'], ['-', '_'], $text);
        return $text;
    }

    public static function safetyBase64Decode($text)
    {
        $text = str_replace(['-', '_'], ['+', '/'], $text);
        $text = base64_decode($text);
        return $text;
    }
}
