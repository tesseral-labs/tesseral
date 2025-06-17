// @ts-expect-error esbuild replaces this at build time, but tsc-check doesn't know that
export const API_URL = __REPLACED_BY_ESBUILD_API_URL__;
// @ts-expect-error ibid
export const DOGFOOD_PROJECT_ID = __REPLACED_BY_ESBUILD_DOGFOOD_PROJECT_ID__;
