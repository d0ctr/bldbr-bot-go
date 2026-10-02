FROM golang:1.26

WORKDIR /app

RUN curl -sfS https://dotenvx.sh/install.sh | sh

COPY ./go.* ./
RUN go mod download

COPY . .
RUN go build -o bot

ARG PORT=8080
EXPOSE $PORT

CMD ["dotenvx", "run", "--", "bot"]
