const webpack = require("webpack");

module.exports = {
    target: "web",
    mode: "production",
    plugins: [
        new webpack.ProvidePlugin({
            process: "process/browser",
            Buffer: ["buffer", "Buffer"],
        }),
    ],
    resolve: {
        fallback: {
            stream: require.resolve("stream-browserify"),
            process: false,
        },
        alias: {
            process: "process/browser",
        },
    },
    output: {
        filename: "xrpl.js",
        library: {
            name: "xrpl",
            type: "var",
            export: "default",
        },
    },
};
