module.exports = {
    root: true,
    env: {
        node: true,
        es2022: true
    },
    extends: [
        // 'plugin:vue/vue3-essential',
        // '@vue/typescript'
        "eslint:recommended",
        "plugin:vue/vue3-recommended",
        "prettier"
    ],
    rules: {
        "vue/require-default-prop": "off",
        "vue/first-attribute-linebreak": [
            "error",
            {
                singleline: "ignore",
                multiline: "ignore"
            }
        ]
        // "@typescript-eslint/no-explicit-any": "off"
    },
    parserOptions: {
        parser: "@typescript-eslint/parser",
        ecmaVersion: 2020
    }
}
