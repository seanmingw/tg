const path = require('path')
const defaultSettings = require('./src/settings.js')
const JavaScriptObfuscator = require('webpack-obfuscator')

function resolve(dir) {
  return path.join(__dirname, dir)
}

const name = defaultSettings.title || 'Trojan Panel'
const port = process.env.port || 8888
module.exports = {
  publicPath: '/',
  outputDir: 'dist',
  assetsDir: 'static',
  lintOnSave: false,
  productionSourceMap: false,
  devServer: {
    port: port,
    open: true,
    overlay: {
      warnings: false,
      errors: true
    },
    proxy: {
      '/api': {
        target: `http://127.0.0.1:8081/api`,
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      }
    }
  },
  configureWebpack:
    process.env.NODE_ENV === 'production'
      ? {
          name: name,
          resolve: {
            alias: {
              '@': resolve('src')
            }
          },
          plugins: [
            // 配置项参考：https://github.com/javascript-obfuscator/javascript-obfuscator
            new JavaScriptObfuscator(
              {
                rotateStringArray: true
              },
              []
            )
          ]
        }
      : {
          name: name,
          resolve: {
            alias: {
              '@': resolve('src')
            }
          }
        },
  chainWebpack(config) {
    // set svg-sprite-loader
    config.module.rule('svg').exclude.add(resolve('src/icons')).end()
    config.module
      .rule('icons')
      .test(/\.svg$/)
      .include.add(resolve('src/icons'))
      .end()
      .use('svg-sprite-loader')
      .loader('svg-sprite-loader')
      .options({
        symbolId: 'icon-[name]'
      })
      .end()
  }
}
