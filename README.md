# go-aws-lambda-s3

# For Linux build use below command to build executable
# setting CGO_ENABLED=0 when building on linux
# bootstrap name is madatory for binary file
GOARCH=arm64 GOOS=linux go build -o bootstrap main.go

# convert the executable to .zip file
zip awslambda.zip bootstrap

# need to provide the s3 access to aws lambda funtion while creating lambda function or after the creation

# Lambda->Functions -> myLambdaFunc -> Configuaratin -> Permissions -> click on RoleName
# ->addpermissions -> attach policies -> search with s3 and give the required access