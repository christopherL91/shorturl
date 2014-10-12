#!/bin/bash

printf "Installing...\n"
pandoc -s -t man shorturl.md -o shorturl.1
mv shorturl.1 /usr/local/share/man/man1
gzip /usr/local/share/man/man1/shorturl.1
printf "Done!\n"
