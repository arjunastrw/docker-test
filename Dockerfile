# Gunakan image resmi Golang sebagai base image
FROM golang:alpine

# Buat direktori kerja di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum ke direktori kerja di dalam container
COPY go.mod .
COPY go.sum .

# Download dan install dependensi proyek
RUN go mod tidy

# Salin seluruh proyek ke direktori kerja di dalam container
COPY . .

# Compile aplikasi Go
RUN go build -o employee

# Ekspor port yang digunakan oleh aplikasi
EXPOSE 8080

# Jalankan aplikasi
ENTRYPOINT [ "./employee" ]
