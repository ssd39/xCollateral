const webpack = require("webpack");
module.exports = function override(config, env) {
  //do stuff with the webpack config...
  let loaders = config.resolve;
  loaders.fallback = {
    fs: false,
    tls: false,
    net: false,
    http: require.resolve("stream-http"),
    stream: require.resolve("stream-browserify"),
    https: require.resolve("https-browserify"),
    crypto: require.resolve("crypto-browserify"),
    url: require.resolve("url/"),
  };
  config.plugins = [
    ...config.plugins,
    new webpack.ProvidePlugin({
      process: "process/browser",
    }),
    new webpack.ProvidePlugin({
      Buffer: ["buffer", "Buffer"],
    }),
  ];
  return config;
};
