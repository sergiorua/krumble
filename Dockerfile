FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN apk add --no-cache \
	ca-certificates \
	--virtual .build-deps \
	git \
	gcc \
	libc-dev \
	libgcc \
        make

COPY . /go/src/github.com/sergiorua/krumble
WORKDIR /go/src/github.com/sergiorua/krumble
ARG BUILD_TAGS="netgo osusergo"
RUN make VK_BUILD_TAGS="${BUILD_TAGS}" build
RUN cp bin/krumble /usr/bin/krumble


FROM scratch
COPY --from=builder /usr/bin/krumble /usr/bin/krumble
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
ENTRYPOINT [ "/usr/bin/krumble" ]
CMD [ "--help" ]
