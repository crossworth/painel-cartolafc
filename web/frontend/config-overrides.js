const { override, fixBabelImports, addLessLoader } = require('customize-cra')

module.exports = override(
  addLessLoader({
    lessOptions: {
      javascriptEnabled: true
    }
  }) //fixme: do we need this
  // fixBabelImports('antd', {
  //   libraryDirectory: 'es',
  //   style: 'css',
  // }),
)
