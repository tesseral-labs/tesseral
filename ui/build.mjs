import { configDotenv } from 'dotenv'
import * as esbuild from 'esbuild'
import { replace } from 'esbuild-plugin-replace'

const UI_BUILD_IS_DEV = process.env.UI_BUILD_IS_DEV === '1'

if (UI_BUILD_IS_DEV) {
  configDotenv({
    path: '../.env',
  })
}

const define = {
  global: 'window',
  ...Object.fromEntries(
    Object.entries(process.env)
      .filter(([k, _v]) => k.startsWith('APP_'))
      .map(([k, v]) => [`process.env.${k}`, JSON.stringify(v)]),
  ),
}

const context = await esbuild.context({
  bundle: true,
  define,
  entryPoints: ['./src'],
  minify: !UI_BUILD_IS_DEV,
  outfile: './public/index.js',
  plugins: [
    replace({
      __API_URL__: process.env.UI_API_URL,
    }),
  ],
  sourcemap: true,
  target: ['chrome58', 'firefox57', 'safari11', 'edge18'],
})

if (UI_BUILD_IS_DEV) {
  console.log('watching')
  await context.watch()
} else {
  await context.rebuild()
  await context.dispose()
}
