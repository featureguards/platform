FROM node:16-alpine

WORKDIR /usr/src/app

COPY package.json /usr/src/app
COPY yarn.lock /usr/src/app

# Installing dependencies
RUN yarn install

COPY . .

RUN yarn build

EXPOSE 80

CMD [ "yarn", "start" ]
