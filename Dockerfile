# First stage: Get Golang image from DockerHub.
FROM golang:1.25.1-alpine3.22 AS backend-builder

# Set our working directory for this stage.
WORKDIR /app

RUN apk add --no-cache gcc musl-dev

# Copy all of our files.
COPY . .

# Get and install all dependencies.
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./migrate.go

FROM node:24.8.0-alpine3.22 AS base-frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# Get corepack
RUN npm i -g corepack@latest
RUN corepack enable
COPY ./ui /frontend/ui
WORKDIR /frontend

FROM base-frontend AS frontend-prod-deps
WORKDIR /frontend/ui
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base-frontend AS frontend-build

WORKDIR /frontend/ui

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

# Last stage: discard everything except our executables.
FROM alpine:latest AS prod

# Set our next working directory.
WORKDIR /app

RUN apk add --no-cache gcc musl-dev

# Copy our executable and our built React application.
COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/migrate .
COPY --from=backend-builder /app/config/production.toml ./config/production.toml
COPY --from=frontend-build /frontend/public ./public

ENV APP_ENV=production

# Declare entrypoints and activation commands.
EXPOSE 8000
ENTRYPOINT ["/bin/sh", "-c" , "./migrate && ./server"]
