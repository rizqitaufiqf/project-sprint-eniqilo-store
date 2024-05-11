FROM golang:1.21.0

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy

EXPOSE 8080

##### Uncomment this for build and pushing to docker hub #####
# ## Build
# FROM golang:1.21.10-bullseye AS build
# WORKDIR /app

# COPY . .
# RUN go mod download

# RUN apt-get update && apt-get install -y dumb-init
# RUN CGO_ENABLED=0 GOOS=linux go build -o /eniqilo-store

# ## Deploy
# FROM gcr.io/distroless/base-debian11
# WORKDIR /

# COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
# COPY --from=build /eniqilo-store /eniqilo-store

# USER nonroot:nonroot

# EXPOSE 8080
# ENTRYPOINT ["/usr/bin/dumb-init", "--"]
# CMD ["/eniqilo-store"]
