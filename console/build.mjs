import * as esbuild from "esbuild";
import * as fs from "fs";
import { configDotenv } from "dotenv";

const CONSOLE_BUILD_IS_DEV = process.env.CONSOLE_BUILD_IS_DEV === "1";

if (CONSOLE_BUILD_IS_DEV) {
  configDotenv({
    path: "./.env",
  });

  const config = {
    CONSOLE_API_URL: process.env.CONSOLE_API_URL,
  };
  fs.writeFileSync("./public/config.json", JSON.stringify(config, null, 2));
}

const context = await esbuild.context({
  bundle: true,
  entryPoints: ["./src"],
  minify: !CONSOLE_BUILD_IS_DEV,
  outfile: "./public/index.js",
  sourcemap: true,
  target: ["chrome58", "firefox57", "safari11", "edge18"],
});

if (CONSOLE_BUILD_IS_DEV) {
  console.log("watching");

  await context.watch();
} else {
  await context.rebuild();
  await context.dispose();
}
