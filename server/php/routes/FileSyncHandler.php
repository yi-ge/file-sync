<?php
class FileSyncHandler
{
    public function post_xhr($json)
    {
        global $database;

        $datetimeFormat = 'Y-m-d H:i:s';

        if (
            !array_key_exists('email', $json) ||
            !array_key_exists('machineId', $json) ||
            !array_key_exists('timestamp', $json) ||
            !array_key_exists('fileId', $json) ||
            !array_key_exists('token', $json)
        ) {
            echo json_encode([
                "status" => -1,
                "msg" => "Missing required parameters",
                "result" => null,
            ]);
            return;
        }

        $emailSha1 = $json['email'];
        $token = $json['token'];

        $user = $database->get("user", "*", [
            "emailSha1" => $emailSha1,
        ]);

        if (!$user) {
            echo json_encode([
                "status" => -2,
                "msg" => "Invalid email",
                "result" => null,
            ]);
            return;
        }

        unset($json["token"]);
        ksort($json);
        $sign = '';
        foreach ($json as $k => $v) {
            $sign .= $k . "=" . $v . "&";
        }
        $sign = $sign . $user["verify"];

        $token = Aes::safetyBase64Decode($token);

        $publicKey = $user["publicKey"];
        $publicKeyId = openssl_pkey_get_public($publicKey);

        if (openssl_verify($sign, $token, $publicKeyId, OPENSSL_ALGO_SHA1) != 1) {
            echo json_encode([
                "status" => -3,
                "msg" => "Invalid token",
                "result" => null,
            ]);
            return;
        }

        if (!array_key_exists('updateAt', $json)) { // Download
            $file = $database->get("file", "*", [
                "email" => $user['email'],
                "fileId" => $json['fileId'],
            ]);

            echo json_encode([
                "status" => 1,
                "msg" => "OK",
                "result" => $file,
            ]);
        } else { // Upload
            if (
                !array_key_exists('content', $json) ||
                !array_key_exists('sha256', $json) ||
                !array_key_exists('fileName', $json)
            ) {
                echo json_encode([
                    "status" => -4,
                    "Missing required parameters",
                    "result" => null,
                ]);
                return;
            }

            $date = new DateTime('now', new DateTimeZone('Asia/Shanghai'));
            $date->setTimestamp(intval($json['updateAt'] / 1000));

            $last_file_id = $database->insert("file", [
                "email" => $user['email'],
                "emailSha1" => $emailSha1,
                "fileId" => $json['fileId'],
                "fileName" => $json['fileName'],
                "content" => $json['content'],
                "sha256" => $json['sha256'],
                "fromMachineId" => $json['machineId'],
                "updateAt" => $date->format($datetimeFormat),
            ]);

            if (function_exists('shmop_open')) {
                $shmid = shmop_open(66, "c", 0755, 320);
                shmop_read($shmid, 0, 320);
                shmop_write($shmid, $json['fileId'], 0);
            }

            echo json_encode([
                "status" => 1,
                "msg" => "OK",
                "result" => [
                    "lastId" => $last_file_id,
                    "fileId" => $json['fileId'],
                ],
            ]);
        }
    }
}
