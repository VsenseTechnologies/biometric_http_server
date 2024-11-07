FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

ENV DB_URL="postgresql://biometric_o6ol_user:QIrtAnXPQW5VlQk5wRCJEq6I9fV9XlZI@dpg-csm4qfggph6c73acqof0-a.oregon-postgres.render.com/biometric"
ENV REDIS_URL="redis://default:E7K1lapur0i7f5J6IruKl661wVBSmMLb@redis-17676.c267.us-east-1-4.ec2.redns.redis-cloud.com:17676"
ENV SERVER_PORT=":8080"
ENV SMTP_USERNAME="cshubhanga@gmail.com"
ENV SMTP_PASSWORD="yvzggalzaxjlwotv"
ENV SMTP_SERVICE="smtp.gmail.com"


CMD ["./main"]
