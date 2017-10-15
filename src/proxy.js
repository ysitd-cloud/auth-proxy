const proxy = require('http-proxy-middleware');

function addProxy(backend) {
  return proxy(backend);
}

module.exports = addProxy;
