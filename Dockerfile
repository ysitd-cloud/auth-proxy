FROM node:6-alpine

RUN yarn install --production

ENTRYPOINT ["yarn", "start"]
