FROM node:24-alpine AS build

WORKDIR /app
COPY ./ .

RUN yarn install && \
    yarn build

FROM busybox:1.37.0 AS serve

WORKDIR /app

COPY --from=build /app/build .

EXPOSE 5173
ENTRYPOINT ["busybox", "httpd", "-f", "-v", "-p", "5173"]

