convert icon-original.png -resize 16x16 icon16.png
convert icon-original.png -resize 48x48 icon48.png
convert icon-original.png -resize 128x128 icon128.png
cd host
go build

