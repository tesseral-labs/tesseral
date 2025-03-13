import eslint from "@eslint/js";
import tseslint from "typescript-eslint";

export default tseslint.config(
  {
    ignores: ["public/index.js", "tailwind.config.js", "build.mjs", "src/gen"],
  },
  {
    rules: {
      "func-style": ["error", "declaration"],
    },
  },
);
