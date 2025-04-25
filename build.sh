go test -v
if [ $? -ne 0 ]; then
    echo "tests failed. Exiting."
    exit 1
fi
go build -o goweasel -ldflags "-s -w" main.go
upx --ultra-brute goweasel