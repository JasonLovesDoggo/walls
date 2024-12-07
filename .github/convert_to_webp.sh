#!/usr/bin/env bash

find . -type f | egrep '.jpeg|.jpg|.tiff|.tif|.png' | parallel --progress 'cwebp -quiet -z 9 -alpha_filter best {} -o {.}.webp'
