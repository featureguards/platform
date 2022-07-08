FROM node:16-alpine as builder

WORKDIR /usr/src/app

COPY . .

RUN yarn install
RUN yarn build && yarn export

FROM nginx

EXPOSE 80

COPY --from=builder /usr/src/app/out /usr/share/nginx/html
COPY --from=builder /usr/src/app/nginx.conf /etc/nginx/nginx.conf