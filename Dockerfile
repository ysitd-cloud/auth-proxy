FROM node:6-alpine

WORKDIR /app

ADD . /app

RUN yarn install --production

ENTRYPOINT ["yarn", "start"]
