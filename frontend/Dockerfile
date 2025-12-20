# syntax=docker/dockerfile:1.7-labs
ARG NODE_VERSION=22

FROM node:${NODE_VERSION}-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .
ENV NITRO_PRESET=node-server
RUN npm run build

FROM node:${NODE_VERSION}-slim AS prod-deps
WORKDIR /app

COPY --from=build /app/.output/server/package*.json ./.output/server/
RUN --mount=type=cache,target=/root/.npm \
    cd .output/server && \
    npm install --omit=dev --ignore-scripts --no-audit --no-fund && \
    npm cache clean --force

FROM gcr.io/distroless/nodejs20-debian12:nonroot AS runtime
WORKDIR /app
ENV NODE_ENV=production


COPY --from=build     /app/.output                         ./.output
COPY --from=prod-deps /app/.output/server/node_modules     ./.output/server/node_modules

CMD [".output/server/index.mjs"]
