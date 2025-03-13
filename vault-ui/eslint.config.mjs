import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import preferArrow from 'eslint-plugin-prefer-arrow';

export default tseslint.config(
  {
    languageOptions: {
      parserOptions: {
        tsconfigRootDir: process.cwd(),
        project: ['./tsconfig.json'],
      },
    },
    plugins: {
      'prefer-arrow': preferArrow,
    },
  },
  {
    ignores: [
      'node_modules',
      'src/gen',
      'src/components/ui',
      'src/hooks',
      '**/*.{js,mjs,jsx}',
    ],
  },
  eslint.configs.recommended,
  tseslint.configs.recommended,
  tseslint.configs.recommendedTypeChecked,
  {
    rules: {
      'prefer-arrow/prefer-arrow-functions': [
        'error',
        {
          disallowPrototype: true,
          singleReturnOnly: false,
          classPropertiesAllowed: false,
        },
      ],
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unsafe-return': 'error',
      '@typescript-eslint/no-misused-promises': [
        'error',
        {
          checksVoidReturn: false,
        },
      ],
      '@typescript-eslint/no-unsafe-assignment': 'error',
      '@typescript-eslint/no-unsafe-member-access': 'error',
      '@typescript-eslint/restrict-template-expressions': 'off',
    },
  },
);
