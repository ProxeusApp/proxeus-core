const webpack = require('webpack')
const path = require('path')
const ESLintPlugin = require('eslint-webpack-plugin')
const { VueLoaderPlugin } = require('vue-loader')

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
      template: 'public/initial.html',
      filename: 'initial.html'
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
          to: '/initial.html'
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
  configureWebpack: function (config) {
    config.output.globalObject = 'this'
    this.optimization = {
      splitChunks: false
    }

    config.plugins.push(
      new webpack.ProvidePlugin({ jQuery: 'jquery', $: 'jquery', 'window.jQuery': 'jquery' }))

    config.plugins.push(
      new ESLintPlugin({
        extensions: ['js'],
        exclude: '/node_modules/',
        options: {
            fix: true,
            configFile: path.resolve(__dirname, '.eslintrc.js')
          }
      }))

    config.plugins.push(
      new VueLoaderPlugin())

    // Is this a good way to configure webpack?
    config.module.rules = [
        {
          test: /\.vue$/,
          loader: 'vue-loader',
          exclude: '/node_modules/'
        },
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: [{
            loader: "babel-loader",
            options: { presets: ['@babel/preset-env'] }
          }]
        }
    ]
  }
}
