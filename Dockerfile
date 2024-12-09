# lightweight container for go
FROM golang:1.22-alpine

# update container's packages and install git
RUN apk update && apk add --no-cache git

# set /app to be the active directory
WORKDIR /app

# copy all files from outside container, into the container
COPY . .

COPY ./conf/app.conf ./conf/app.conf

# download dependencies
#RUN go mod tidy
COPY ./sigmatech-xyz ./sigmatech-xyz
# build binary
#RUN go build -o sigmatech-xyz

# set the entry point of the binary
ENTRYPOINT ["/app/sigmatech-xyz"]
