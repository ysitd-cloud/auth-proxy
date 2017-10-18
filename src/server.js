const express = require('express');
const morgan = require('morgan');
const helmet = require('helmet');
const createVirtualHost = require('./vhost');
const proxies = require('../proxy.json');

const app = express();

app.use(helmet.hidePoweredBy());
app.use(helmet.frameguard());
app.use(helmet.xssFilter());
app.use(morgan(':req[host] - ":method :url HTTP/:http-version" :status ":referrer" ":user-agent"'));

proxies.forEach((proxy) => {
  app.use(createVirtualHost(proxy.host, proxy.backend, proxy.idKey, proxy.secretKey));
});

app.listen(80);
