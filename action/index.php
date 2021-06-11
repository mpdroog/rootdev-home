<?php
namespace prj;
require __DIR__ . "/vendor/core/init-browser.php";
require __DIR__ . "/vendor/core/Taint.php";
require __DIR__ . "/vendor/core/Env.php";
require __DIR__ . "/vendor/core/Res.php";

trait ProjectValidators {}
\core\Taint::init();
\core\Env::init();

$d = sprintf("%s/cmp/%s/index.php", __DIR__, \core\Env::hfastPath());
if (! file_exists($d)) {
  \core\Res::error(404);
  exit("Page not found - action.");
}

require $d;
