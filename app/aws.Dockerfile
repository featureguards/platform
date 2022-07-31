FROM node:16-alpine

WORKDIR /usr/src/app

ARG MIXPANEL_ID
ENV NEXT_PUBLIC_MIXPANEL_ID=${MIXPANEL_ID}

COPY package.json /usr/src/app
COPY yarn.lock /usr/src/app

# Installing dependencies
RUN yarn install

COPY . .

RUN yarn build

EXPOSE 3000

CMD [ "yarn", "start" ]
