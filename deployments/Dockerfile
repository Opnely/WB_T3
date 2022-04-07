FROM golang

WORKDIR /WB_T3_docker

COPY config.toml database* go.* main.go model* router* ./

RUN go build

CMD [ "./service" ]
