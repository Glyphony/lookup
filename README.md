# Installation Instructions

## Download GoLang 1.18
`curl -OL https://golang.org/dl/go1.18.linux-amd64.tar.gz`

## Extract Files
`sudo tar -C /usr/local -xvf go1.18.linux-amd64.tar.gz`

## Add GOPATH to your profile
`sudo nano ~/.profile`

## Insert the following
`export PATH=$PATH:/usr/local/go/bin`

## Update profile
`source ~/.profile`

- Copy main.go to any directory

## Create go.mod file
`go mod init`

## Get lookup package (created by me) from github
`go get github.com/Glyphony/lookup`
OR
`go mod download github.com/Glyphony/lookup`

## Run Program
`go run main.go <ip address>`

## Considerations
1. If this is being run from a new docker container then you will probably need to run `apt-get update` so that you can then do `apt-get install -y curl`
2. If nano is not installed then run `apt-get install -y nano`
