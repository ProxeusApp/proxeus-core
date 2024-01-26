const webpack = require('webpack')
const path = require('path')

module.exports = {
  assetsDir: 'static/assets/',
  outputDir: path.resolve(__dirname, 'dist'),
  runtimeCompiler: true,
  productionSourceMap: false,

  css: {
    // extract CSS in components into a single CSS file (only in production)
    // can also be an object of options to pass to extract-text-webpack-plugin
    extract: true
  },
  pages: {
    initial: {
      entry: './src/initial.js',
      template: 'index.html',
      filename: 'index.html'
    },
    app: {
      entry: './src/main.js',
      template: 'public/app.html',
      filename: 'app.html'
    },
    user: {
      entry: './src/user.js',
      template: 'public/user.html',
      filename: 'user.html'
    },
    frontend: {
      entry: './src/frontend.js',
      template: 'public/frontend.html',
      filename: 'frontend.html'
    }
  },
  devServer: {
    port: 3005,
    historyApiFallback: {
      rewrites: [
        {
          from: /^\/init$/,
          to: '/index.html'
        },
        {
          from: /^\/$/,
          to: '/frontend.html'
        },
        {
          from: /^\/p\//,
          to: '/app.html'
        },
        {
          from: /^\/admin/,
          to: '/app.html'
        },
        {
          from: /^\/document/,
          to: '/user.html'
        },
        {
          from: /^\/user/,
          to: '/user.html'
        },
        {
          from: /.*/,
          to: '/frontend.html'
        }
      ]
    },
    proxy: {
      '/api': {
        target: 'http://localhost:1323',
        ws: false,
        changeOrigin: false
      },
      '/static': {
        target: 'http://localhost:1323',
        ws: false,
        changeOrigin: false
      }
    }
  },
  chainWebpack: config => {
    config.module
      .rule('eslint')
      .use('eslint-loader')
      .loader('eslint-loader')
      .tap(options => {
        options.configFile = path.resolve(__dirname, '.eslintrc.js')
        options.fix = true
        return options
      })
    // optionally replace with another progress output plugin
    // `npm i -D simple-progress-webpack-plugin` to use
    // config.plugin('simple-progress-webpack-plugin').use(require.resolve('simple-progress-webpack-plugin'), [
    //   {
    //     format: 'compact', // options are minimal, compact, expanded, verbose
    //   },
    // ])
  },
  configureWebpack: function (config) {
    config.output.globalObject = 'this'
    this.optimization = {
      splitChunks: false
    }

    config.plugins.push(
      new webpack.ProvidePlugin({ jQuery: 'jquery', $: 'jquery', 'window.jQuery': 'jquery' }))
  }
}
