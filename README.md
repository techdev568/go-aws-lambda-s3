# go-aws-lambda-sample

# For Linux build use below command to build executable
# setting CGO_ENABLED=0 when building on linux

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

# convert the executable to .zip file
