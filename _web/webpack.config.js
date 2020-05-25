var path = require('path')

var APP_DIR = path.resolve(__dirname + '/src')
var BUILD_DIR = path.resolve(__dirname + '/dist')

module.exports = {
  entry: APP_DIR + '/index.jsx',
  output: {
    path: BUILD_DIR,
    filename: 'bundle.js'
  },
  devServer: {
    inline: true,
    host:'0.0.0.0',
    contentBase: BUILD_DIR,
    port: 3333,
    hot: true
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader"
        }
      }
    ]
  }
};
