FROM golang:1.16
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/main ./
ENV MENU-OWNER-EMAIL owner@email.com
ENV MENU-OWNER-PASSWORD Password12345.
ENV MENU-TYPE POSTGRES
ENV MENU-USER postgres
ENV MENU-PASSWORD Password12345.
ENV MENU-PORT 5432
ENV MENU-NAME-DB postgres
ENV MENU-HOST localhost
ENV MENU-RSA None
ENV MENU-RSA-PUB None
ENV MENU-COOKIE-KEY MP}mn!=v=xw#fE_Jj{}?PFS%kB;$78hB
CMD ["./main"]  
