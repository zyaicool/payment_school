# Stage 1: Build
FROM golang:1.22.7-alpine AS build

ENV GOPROXY=https://proxy.golang.org,direct \
    GO111MODULE=auto \
    PATH=$PATH:$HOME/go/bin \
    TZ=Asia/Jakarta

# Install necessary packages
RUN apk update && apk add --no-cache git bash wget unzip tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set up the working directory and copy the application source code
WORKDIR /main
ADD . /main/

# Copy the liquibase.properties file into the container
COPY liquibase-dev.properties /main/liquibase-dev.properties
COPY liquibase-sit.properties /main/liquibase-sit.properties
COPY liquibase-demo.properties /main/liquibase-demo.properties

# Copy the .env file into the container
COPY .env /main/.env

# Install dependencies and build the application
RUN go mod download && go build -ldflags="-w -s" -o app .

# Stage 2: Final Runtime Image
FROM alpine:3.13

WORKDIR /main

# Copy the application binary and other necessary files from the build stage
COPY --from=build /main/app /main/app
COPY --from=build /main/.env /main/.env

COPY --from=build /main/liquibase-dev.properties /main/liquibase-dev.properties
COPY --from=build /main/liquibase-sit.properties /main/liquibase-sit.properties
COPY --from=build /main/liquibase-demo.properties /main/liquibase-demo.properties

COPY --from=build /main/db /main/db
COPY --from=build /main/data /main/data

ENV TZ=Asia/Jakarta
RUN apk add --no-cache bash tzdata openjdk11 \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Optional: Copy Liquibase (if required in runtime)
RUN wget https://developer:DiGif0rm%21@nexus.digiform.co.id/repository/liquibase-package/liquibase-4.29.2.zip \
    && unzip liquibase-4.29.2.zip -d /opt/liquibase \
    && rm liquibase-4.29.2.zip
ENV PATH=$PATH:/opt/liquibase

# Set file descriptors and inotify limits
RUN echo "fs.inotify.max_user_watches=2099999999" >> /etc/sysctl.conf && \
    echo "fs.inotify.max_user_instances=2099999999" >> /etc/sysctl.conf && \
    echo "fs.inotify.max_queued_events=2099999999" >> /etc/sysctl.conf

EXPOSE 8081

ENTRYPOINT ["/main/app"]
