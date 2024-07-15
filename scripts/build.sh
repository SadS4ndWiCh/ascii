rm -fr dist

echo 'building darwin-amd64...'
GOOS=darwin GOARCH=amd64 go build -o dist/darwin-amd64/ascii
echo 'building darwin-arm64...'
GOOS=darwin GOARCH=arm64 go build -o dist/darwin-arm64/ascii

echo 'building linux-amd64...'
GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/ascii
echo 'building linux-arm64...'
GOOS=linux GOARCH=arm64 go build -o dist/linux-arm64/ascii

echo 'building windows-amd64...'
GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/ascii
echo 'building windows-386...'
GOOS=windows GOARCH=386 go build -o dist/windows-386/ascii

cd dist
for dir in $(ls -d *); do
    tar cfzv "$dir".tgz $dir
    rm -rf $dir
done

echo 'done!'