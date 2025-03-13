// @ts-check

import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  eslint.configs.recommended,
  tseslint.configs.recommended,
  {
    ignores: ['public/index.js', 'tailwind.config.js', 'build.mjs', 'src/gen'],
  },
);
