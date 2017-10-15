const passport = require('passport');
const OAuth2Strategy = require('passport-oauth2');
const axios = require('axios');

passport.serializeUser((user, done) => {
  done(null, user);
});

passport.deserializeUser((user, done) => {
  done(null, user);
});

const { OAUTH_HOST } = process.env;

const authorizationURL = `https://${OAUTH_HOST}/oauth/authorize`;
const tokenURL = `https://${OAUTH_HOST}/oauth/token`;

function addProvider(host, idKey, secretKey) {
  passport.use(host, new OAuth2Strategy({
    authorizationURL,
    tokenURL,
    clientID: process.env[idKey],
    clientSecret: process.env[secretKey],
    callbackURL: `https://${host}/auth/ycloud/callback`,
  }, (accessToken, refreshToken, profile, cb) => {
    axios.get(`${OAUTH_HOST}/api/v1/user/info`, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    })
      .then((response) => {
        if (response.status !== 200) {
          throw new Error('Fail to get user info');
        } else {
          cb(null, {
            provider: 'ycloud',
            id: response.data.username,
            displayName: response.data.display_name,
            emails: [
              { value: response.data.email },
            ],
            photos: [
              { value: response.data.avatar_url },
            ],
            oauth: {
              accessToken,
              refreshToken,
            },
          });
        }
      })
      .catch(cb);
  }));
}
module.exports = addProvider;
