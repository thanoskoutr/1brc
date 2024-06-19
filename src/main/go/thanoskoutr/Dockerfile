FROM golang:1.22 AS build
WORKDIR /app
COPY . ./
RUN make all

FROM scratch AS run
WORKDIR /
COPY --from=build /app/1brc /
