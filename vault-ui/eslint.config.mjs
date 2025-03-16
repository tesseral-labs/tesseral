import * as reactHooks from "eslint-plugin-react-hooks";
import pluginJs from "@eslint/js";
import pluginReact from "eslint-plugin-react";
import globals from "globals";
import tseslint from "typescript-eslint";

/** @type {import("eslint").Linter.Config[]} */
// eslint-disable-next-line no-restricted-exports
export default [
  { files: ["**/*.{js,mjs,cjs,ts,jsx,tsx}"] },
  { languageOptions: { globals: globals.browser } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  pluginReact.configs.flat.recommended,
  reactHooks.configs["recommended-latest"],
  {
    ignores: ["src/gen/**/*", "src/components/ui/*"],
  },
  {
    settings: {
      react: {
        version: "detect",
      },
    },
    rules: {
      "react/prop-types": "off",
      "react/no-unescaped-entities": "off",
      "react-hooks/exhaustive-deps": "error",
      "func-style": ["error", "declaration"],
      "no-restricted-exports": [
        "error",
        { restrictDefaultExports: { direct: true } },
      ],
    },
  },
];
