const { defineConfig } = require("@vue/cli-service");

module.exports = defineConfig({
  configureWebpack: {
    resolve: {
      fallback: {
        fs: false,
        tls: false,
        net: false,
        path: false,
        zlib: false,
        http: false,
        https: false,
        stream: false,
        crypto: false,
      },
      extensions: [".js", ".jsx", ".json", ".ts", ".tsx"],
    },
  },

  pluginOptions: {
    electronBuilder: {
      nodeIntegration: true,
    },
    quasar: {},
  },
});
