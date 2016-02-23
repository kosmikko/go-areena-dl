# Installation

To install, use go get and make install.

```
$ go get -d github.com/kosmikko/go-areena-dl
$ cd $GOPATH/src/github.com/kosmikko/go-areena-dl
$ make install
```

## Building for Raspberry Pi

`env GOOS=linux GOARCH=arm GOARM=5 go build -v github.com/kosmikko/go-areena-dl`

## Usage

You need to define the following environment variables (get these from http://developer.yle.fi/):
```
export YLE_APP_ID="YOUR_APP_ID"
export YLE_APP_KEY="YOUR_APP_KEY"
export YLE_SECRET="YOUR_APP_SECRET"
```

Download video(s) by defining program_ids separated by comma:

`go-areena-dl download PROGRAM_ID1,PROGRAM_ID2`
