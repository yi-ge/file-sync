<?php

class DeviceRemoveHandler
{
    public function post_xhr($json)
    {
        global $database;

        if (
            !array_key_exists('email', $json) ||
            !array_key_exists('machineId', $json) ||
            !array_key_exists('timestamp', $json) ||
            !array_key_exists('removeMachineId', $json) ||
            !array_key_exists('machineKey', $json) ||
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
        $removeMachineId = $json['removeMachineId'];
        $machineId = $json['machineId'];
        $machineKey = $json['machineKey'];
        $token = $json['token'];

        $user = $database->get("user", "*", [
            "emailSha1" => $emailSha1,
        ]);

        if (!$user) {
            $database->insert("log", [
                "email" => $user['email'],
                "machineId" => $machineId,
                "action" => "DeviceRemove",
                "content" => json_encode([
                    "success" => false,
                    "status" => -2,
                    "msg" => "Invalid email",
                    "emailSha1" => $emailSha1,
                    "machineId" => $machineId,
                    "machineKey" => $machineKey,
                    "removedMachineId" => $removeMachineId,
                ]),
            ]);
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
            $database->insert("log", [
                "email" => $user['email'],
                "machineId" => $machineId,
                "action" => "DeviceRemove",
                "content" => json_encode([
                    "success" => false,
                    "status" => -3,
                    "msg" => "Invalid token",
                    "email" => $user['email'],
                    "machineId" => $machineId,
                    "machineKey" => $machineKey,
                    "removedMachineId" => $removeMachineId,
                ]),
            ]);

            echo json_encode([
                "status" => -3,
                "msg" => "Invalid token",
                "result" => null,
            ]);
            return;
        }

        $machine = $database->get("device", "*", [
            "email" => $user['email'],
            "machineId" => $machineId,
        ]);

        if (sha1($machine['machineKey']) !== $machineKey) {
            $database->insert("log", [
                "email" => $user['email'],
                "machineId" => $machineId,
                "action" => "DeviceRemove",
                "content" => json_encode([
                    "success" => false,
                    "status" => -4,
                    "msg" => "Invalid machineKey",
                    "email" => $user['email'],
                    "machineId" => $machineId,
                    "machineKey" => $machineKey,
                    "removedMachineId" => $removeMachineId,
                ]),
            ]);

            echo json_encode([
                "status" => -4,
                "msg" => "Invalid machineKey",
                "result" => null,
            ]);
            return;
        }

        $res = $database->delete("device", [
            "machineId" => $removeMachineId,
        ]);

        if ($res->rowCount() == 1) {
            $database->insert("log", [
                "email" => $user['email'],
                "machineId" => $machineId,
                "action" => "DeviceRemove",
                "content" => json_encode([
                    "success" => true,
                    "email" => $user['email'],
                    "machineId" => $machineId,
                    "machineKey" => $machineKey,
                    "removedMachineId" => $removeMachineId,
                ]),
            ]);

            echo json_encode([
                "status" => 1,
                "msg" => "OK",
                "result" => [
                    "removedMachineId" => $removeMachineId,
                ],
            ]);
        } else {
            $database->insert("log", [
                "email" => $user['email'],
                "machineId" => $machineId,
                "action" => "DeviceRemove",
                "content" => json_encode([
                    "success" => false,
                    "status" => -5,
                    "msg" => "Failed to remove",
                    "email" => $user['email'],
                    "machineId" => $machineId,
                    "machineKey" => $machineKey,
                    "removedMachineId" => $removeMachineId,
                ]),
            ]);

            echo json_encode([
                "status" => -5,
                "msg" => "Failed to remove",
                "result" => null,
            ]);
        }
    }
}
