FROM node:16-alpine

WORKDIR /usr/src/app

COPY . .

# Installing dependencies
CMD "yarn" "install"

EXPOSE 3000