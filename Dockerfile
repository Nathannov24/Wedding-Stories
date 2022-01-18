#ambil image yang dibutuhkan untuk me run app kita
FROM golang:1.17-alpine

# set direktori kerja yang akan digunakan
WORKDIR /app

# copy ke direktory app
COPY ./go.mod ./
COPY ./go.sum ./
# jalankan go mod download
RUN go mod download
#copy main.go ke direktory app
COPY . .
#jalankan build untuk membuat file executable
RUN go build -o program
# jalankan file executable hasil build
CMD ./program