FROM golang AS builder
WORKDIR /code
COPY . .
RUN make

#---

FROM ubuntu:latest
WORKDIR /code
COPY --from=builder /code/bin/ministub .
ENTRYPOINT ["./ministub"]