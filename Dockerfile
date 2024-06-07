FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /DiskordServer

EXPOSE 3000

ENV PORT=3000 \
    CORS_WHITELIST="localhost" \
    SECRET=324435rfsdgt43w43tsdf354t \
    DISCORD_TOKEN=MTIzNzEyOTA4MDkyMzQ4ODM0OQ.GILCPM.WtLbxZlU9qQMik3ShOZhkqBCPtC-xo7X2JhuAU \
    DISCORD_CHANEL=1237128843039604918 \
    DB="host=ep-falling-darkness-a24v5206.eu-central-1.aws.neon.tech user=main_owner password=LtPhBibyAK83 dbname=main port=5432 sslmode=require"

CMD [ "/DiskordServer" ]