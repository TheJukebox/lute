FROM node:24-alpine AS build

WORKDIR /app
COPY ./ .

RUN yarn install && \
    yarn build

FROM node:24-alpine AS serve

WORKDIR /app
COPY --from=build /app/build .

ENTRYPOINT ["node", "index.js"]
