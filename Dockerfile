FROM golang AS builder
WORKDIR /code
ENV USER=ministub
ENV UID=10001
RUN adduser --disabled-password --gecos "" --home "/nonexistant" --shell "/sbin/nologin" --no-create-home --uid "${UID}" "${USER}"
COPY . .
RUN make

#---

FROM scratch
COPY --from=builder /code/bin/ministub /code/ministub
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER ministub:ministub
CMD ["/code/ministub"]