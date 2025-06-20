FROM golang:1.24.2-bookworm AS builder

RUN useradd -m builder
WORKDIR /src

ADD ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN chown -R builder:builder /src
USER builder

RUN go build -a -o main .

# Use eclipse-temurin as the base image for Java
FROM eclipse-temurin:21-jdk

# Install Go, Supervisor, and wget
RUN apt-get update && apt-get install -y \
    supervisor \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Copy Supervisor configuration
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /app

# Download Aeron all-in-one JAR
RUN wget -O aeron-all.jar https://repo1.maven.org/maven2/io/aeron/aeron-all/1.48.2/aeron-all-1.48.2.jar

# Copy Go source code
COPY --from=builder /src/main .
COPY low-latency.properties /app/md.properties


# Run Supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]


