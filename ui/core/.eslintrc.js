module.exports = {
  
  root: true,
  env: {
    node: true,
    browser: true,
    commonjs: true,
    es6: true,
    jquery: true
  },
  'extends': [
    'plugin:vue/essential',
    '@vue/standard',
    'eslint:recommended',
    'plugin:prettier/recommended'
  ],
  rules: {
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-undef': 'off',
    'no-useless-escape': 'off',
    'vue/require-v-for-key': 'off',
    'vue/valid-template-root': 'off',
    'vue/no-mutating-props': ['error', {
      shallowOnly: true
    }],
    'vue/multi-word-component-names': 'off',
    'standard/no-callback-literal': 'off'
  },
  parserOptions: {
    parser: '@babel/eslint-parser'
  }
}
