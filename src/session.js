const { Router } = require('express');
const passport = require('passport');
const session = require('express-session');
const RedisStore = require('connect-redis')(session);

function loginGuard(req, res, next) {
  if (!req.user) {
    res.redirect('/auth/login');
    res.end();
  } else {
    next();
  }
}


function createAuthRouter(host) {
  const router = Router();
  router.get('/ycloud', passport.authenticate(host));
  router.get('/ycloud/callback', passport.authenticate(host, {
    failureRedirect: '/auth/login?error=fail&provider=ycloud',
  }), (req, res) => {
    res.redirect('/admin');
  });
  router.get('/login', (req, res) => {
    res.redirect('/auth/ycloud');
  });
  return router;
}

function addSession(app, host) {
  app.use(session({
    store: new RedisStore({
      host: process.env.REDIS_HOST,
      prefix: `sess:${host}`,
    }),
    name: `ysitd.${host}`,
    resave: false,
    saveUninitialized: false,
  }));
  app.use(passport.initialize());
  app.use(passport.session());

  app.use('/auth', createAuthRouter(host));
  app.use(loginGuard);
}

module.exports = addSession;
