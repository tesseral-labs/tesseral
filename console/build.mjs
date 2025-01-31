import * as esbuild from 'esbuild'
import { configDotenv } from 'dotenv'

const APP_BUILD_IS_DEV = process.env.APP_BUILD_IS_DEV === '1'

if (APP_BUILD_IS_DEV) {
  configDotenv({
    path: '../.env',
  })
}

const context = await esbuild.context({
  bundle: true,
  define: {
    __REPLACED_BY_ESBUILD_API_URL__: JSON.stringify(process.env.APP_API_URL),
  },
  entryPoints: ['./src'],
  minify: !APP_BUILD_IS_DEV,
  outfile: './public/index.js',
  sourcemap: true,
  target: ['chrome58', 'firefox57', 'safari11', 'edge18'],
})

if (APP_BUILD_IS_DEV) {
  console.log('watching')
  await context.watch()
} else {
  await context.rebuild()
  await context.dispose()
}
