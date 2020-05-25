FROM node:12.2.0-alpine AS builder

ARG APP_VERSION

WORKDIR /app

# Install dependencies
COPY _web/yarn.lock _web/package.json ./
RUN yarn install --frozen-lockfile

# Build app
COPY _web/src src
COPY _web/dist dist
COPY _web/.babelrc _web/webpack.config.js ./
RUN sed -i s/{APP_VERSION}/${APP_VERSION}/g dist/config.json
RUN yarn build

# Copy artifacts in slimmer image and start serving
FROM mhart/alpine-node

RUN yarn global add serve

WORKDIR /app

COPY --from=builder /app/dist .

CMD ["serve", "-p", "80", "-s", "."]
