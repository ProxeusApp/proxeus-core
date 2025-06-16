const webpack = require("webpack");
const path = require("path");

module.exports = {
  assetsDir: "static/assets/",
  outputDir: path.resolve(__dirname, "dist"),
  runtimeCompiler: true,
  productionSourceMap: false,

  css: {
    // extract CSS in components into a single CSS file (only in production)
    // can also be an object of options to pass to extract-text-webpack-plugin
    extract: true,
  },
  pages: {
    initial: {
      entry: "./src/initial.js",
      template: "public/initial.html",
      filename: "initial.html",
    },
    app: {
      entry: "./src/main.js",
      template: "public/app.html",
      filename: "app.html",
    },
    user: {
      entry: "./src/user.js",
      template: "public/user.html",
      filename: "user.html",
    },
    frontend: {
      entry: "./src/frontend.js",
      template: "public/frontend.html",
      filename: "frontend.html",
    },
  },
  devServer: {
    port: 3005,
    historyApiFallback: {
      rewrites: [
        {
          from: /^\/init$/,
          to: "/initial.html",
        },
        {
          from: /^\/$/,
          to: "/frontend.html",
        },
        {
          from: /^\/p\//,
          to: "/app.html",
        },
        {
          from: /^\/admin/,
          to: "/app.html",
        },
        {
          from: /^\/document/,
          to: "/user.html",
        },
        {
          from: /^\/user/,
          to: "/user.html",
        },
        {
          from: /.*/,
          to: "/frontend.html",
        },
      ],
    },
    proxy: {
      "/api": {
        target: "http://localhost:1323",
        ws: false,
        changeOrigin: false,
      },
      "/static": {
        target: "http://localhost:1323",
        ws: false,
        changeOrigin: false,
      },
    },
  },
  lintOnSave: false, // Disable ESLint temporarily to fix build issues
  //transpileDependencies: [
    // Force transpilation of problematic dependencies
    //"pdfjs-dist",
  //],
  chainWebpack: (config) => {
    // Disable ESLint plugin to avoid version conflicts
    config.plugins.delete("eslint");

    // Exclude babel loaders
    config.module
      .rule("js")
      .test(/\.js$/)
      .use("babel-loader")
      .loader("babel-loader")
      .options({
        exclude: /node_modules/,
        compact: false
      })
      .end();

    // Exclude worker files from normal JavaScript processing
    config.module
      .rule("js")
      .exclude.add(/\.worker\.js$/)
      .end();

    // Add specific rule for PDF.js worker files
    config.module
      .rule("worker")
      .test(/\.worker\.js$/)
      .use("worker-loader")
      .loader("worker-loader")
      .options({
        inline: "fallback",
      })
      .end();
  },
  configureWebpack: function (config) {
    config.output.globalObject = "this";

    // Disable splitChunks
    config.optimization = {
      splitChunks: false,
    };

    // Configure module resolution aliases for Node.js polyfills
    config.resolve = config.resolve || {};
    config.resolve.alias = {
      ...config.resolve.alias,
      stream: "stream-browserify",
      crypto: "crypto-browserify",
      buffer: "buffer",
      process: "process/browser",
      util: "util",
      assert: "assert",
      url: "url",
      fs: false,
      net: false,
      tls: false,
      child_process: false,
      "pdfjs-dist/build/pdf.worker.js": "pdfjs-dist/build/pdf.worker.min.js",
    };

    // Provide global variables
    config.plugins.push(
      new webpack.ProvidePlugin({
        jQuery: "jquery",
        $: "jquery",
        "window.jQuery": "jquery",
        Buffer: ["buffer", "Buffer"],
        process: "process/browser",
      })
    );

    // Define environment variables
    config.plugins.push(
      new webpack.DefinePlugin({
        "process.env.NODE_ENV": JSON.stringify(
          process.env.NODE_ENV || "development"
        ),
        global: "globalThis",
      })
    );
  },
};
