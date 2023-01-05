FROM golang:1.18 as builder
WORKDIR /go/src/github.com/tomoconnor/golyricgraph
COPY go.mod go.sum *.go ./
RUN go get

RUN GOOS=linux go build -a -o lyricgraph-server .
RUN mkdir dotfiles
RUN mkdir imagefiles

FROM ubuntu:latest
WORKDIR /srv
COPY --from=builder /go/src/github.com/tomoconnor/golyricgraph/lyricgraph-server /srv
RUN mkdir /srv/dotfiles
RUN mkdir /srv/imagefiles

CMD [ "./lyricgraph-server", "-server" ]