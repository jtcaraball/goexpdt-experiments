# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21.6 AS build-stage

RUN apt-get update && apt-get install unzip

# Get kissat repo
ADD https://github.com/arminbiere/kissat/archive/refs/tags/rel-3.1.1.zip rel-3.1.1.zip
RUN unzip rel-3.1.1.zip
RUN mv kissat-rel-3.1.1 /kissatdir
RUN rm rel-3.1.1.zip

# Compile solver binary
WORKDIR /kissatdir
RUN ./configure
RUN make

# Move binary to root and remove left overs
WORKDIR /
RUN mv kissatdir/build/kissat kissat
RUN rm -r kissatdir

# Compile go application
COPY . .
RUN go build -o main .

# Set entrypoint
ENTRYPOINT ["./main"]
