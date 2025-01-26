go build -o goweasel -ldflags "-s -w" main.go
upx --ultra-brute goweasel