const { Router } = require('express');
const addProxy = require('./proxy');
const addProvider = require('./oauth');
const addSession = require('./session');

function createVirtualHost(host, backend, idKey, secretKey) {
  const app = Router();

  addProvider(host, idKey, secretKey);
  addSession(app, host);

  app.use(addProxy(backend));
  return app;
}

module.exports = createVirtualHost;
