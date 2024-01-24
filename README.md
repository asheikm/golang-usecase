# golang-usecase
MyFavArtist App - This application will use last.fm to get popular artist based on region.

# Getting started

To get this app running locally:

- Clone this repo
- Install Golang programming language ([instructions](https://golang.org/doc/install))
- Install Docker ([instructions](https://docs.docker.com/engine/install/))


# Basic Configuration

Configuration 
   
   - This app uses listening port 8080, need to update the lastfm auth key to run

# Build and run

# Compile and generate binary or docker image:

- How to compile and generate binary
   - Put the code in GOPATH which is usually ($HOME/go/src)
   - Optional: Run `go install ` will get the dependency packages (Go mod take care of this anyway)
   - cd to myfavartist directory, Run `go build -o bin/myfavartist main.go `(This will generate binary file name myfavartist in bin directory)
   - 

# Building and runnning app from docker image:

   - cd to myfavartist directory
   - To build into docker image `docker build . -t myfavartist:1.0 `
   - To list images from the local docker hub repo `docker images`
   - To run docker image ` docker run -it -d -p 8080:8080 myfavartist:1.0` (This will run myfavartist app and listen on port 8080)


# Unit Test

Unit Testing
    
   - This app uses go programming languages default test framework 
   - cd to myfavartist directory, Run `go test ./...` from project directory, this command will get the result something similar below 
     
     `[user@dev utils]$ go test -v ./...`

 
# Application Usage
   
Curl command to get artist detail based on region 
   
  - Run `curl -X GET -H "Content-type: application/json" \-H "Accept: application/json"  "http://localhost:8080/api/v1/artist/italy"` (This will fectch if the artist and track details)


