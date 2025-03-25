import * as esbuild from "esbuild";

const UI_BUILD_IS_DEV = process.env.UI_BUILD_IS_DEV === "1";

const context = await esbuild.context({
  bundle: true,
  entryPoints: ["./src"],
  minify: !UI_BUILD_IS_DEV,
  outfile: "./dist/index.js",
  sourcemap: true,
  target: ["chrome58", "firefox57", "safari11", "edge18"],
});

if (UI_BUILD_IS_DEV) {
  console.log("watching");
  await context.watch();
} else {
  await context.rebuild();
  await context.dispose();
}
