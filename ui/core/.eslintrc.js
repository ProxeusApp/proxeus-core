module.exports = {
  
  root: true,
  env: {
    node: true,
    "browser": true,
    "commonjs": true,
    "es6": true,
    "jquery": true
  },
  'extends': [
    'eslint:recommended',
    'plugin:vue/essential'
  ],
  rules: {
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-console': 'off',
    'no-undef': 'off',
    'no-empty': 'off',
    'no-unused-vars': 'off',
    'no-useless-escape': 'off',
    'vue/require-v-for-key': 'off',
    'vue/valid-template-root': 'off',
    'vue/no-mutating-props': ["error", {
      "shallowOnly": true
    }],
    'vue/multi-word-component-names': 'off',
    'standard/no-callback-literal': 'off'
  },
  parserOptions: {
    parser: '@babel/eslint-parser'
  }
}
