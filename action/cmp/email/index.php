<?php
require ROOT . "/action/vendor/core/Safe.php";
use core\Taint;
use core\Db;
use core\Safe;
use core\Res;

class Input {
  public $body;
  public $email;
  public function rules() {
    return ["body" => ["text"], "email" => ["email"]];
  }
}

$input = Taint::post(new Input());
if (is_array($input)) {
  Res::error("Invalid input.", []);
  exit;
}

function email(array $curl_post_data) {
  $url = "https://api.mailgun.net/v3/rootdev.nl/messages";
  $curl = curl_init($url);
  curl_setopt($curl, CURLOPT_HTTPAUTH, CURLAUTH_BASIC);
  curl_setopt($curl, CURLOPT_USERPWD, "api:key-a68e2028194a2247743c5f9fca525025");
  curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
  curl_setopt($curl, CURLOPT_POST, true);
  curl_setopt($curl, CURLOPT_POSTFIELDS, $curl_post_data);
  curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, false);

  $res = curl_exec($curl);
  $http = curl_getinfo($curl, CURLINFO_HTTP_CODE);
  if ($http < 200 || $http > 299) {
    user_error(sprintf("HTTP(%s) => (%d) %s", $url, $http, $res));
  }
  curl_close($curl);
  $json = json_decode($res, true);
  if ($json["message"] !== "Queued. Thank you.") {
    user_error(sprintf("Mailgun::email err=$res"));
  }
}

$ignores = [
    "http://",
    "https://",
    "www.",
];
foreach ($ignores as $ignore) {
  if (strpos($input->body, $ignore) !== false) {
    Res::error(400);
    exit;
  }
}
if (strpos($input->email, "@rootdev.nl")) {
  Res::error(400);
  exit;
}

email([
  "from" => "robot@rootdev.nl",
  "to" => "rootdev@gmail.com",
  "subject" => "Contact form",
  "text" => "From: " . $input->email . "\n\nBody: " . $input->body
]);

