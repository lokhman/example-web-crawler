# Example Web Crawler
This is the example web crawler written in Go, that concurrently outputs a simple site map for the given URL.

- It will stream and output the URLs extracted from all available pages within one host in the real time.
- It will not revisit the visited pages.
- It will run workers concurrently, so the order of streamed URLs is not defined.
- It will stop and return error if the number of visited pages reaches the limit, specified in the arguments
  (defaults to `1000`).
- It will log all possible errors from the workers to `stderr` stream.
- It was implemented to be production ready and cover as many edge cases as possible, although it still has room for
  improvement (technically and functionally).
- It has 100% code coverage with tests.

## How to build
To build the program you will require Go compiler (ideally 1.11+) and Dep tool:

    $ cd $GOPATH/src/github.com/lokhman/example-web-crawler
    $ dep ensure -vendor-only
    $ go build .

## How to start
To start the program you should pass the base URL as the first argument to the executable:

    $ ./example-web-crawler https://example.com/

You may also specify the maximum pages limit with `-max-pages` flag:

    $ ./example-web-crawler -max-pages=100 https://example.com/

## How to test
The unit tests can be run with:

    $ go test -v ./...

The code coverage report can be generated with:

    $ go test -v -coverprofile cover.out ./...
    $ go tool cover -html=cover.out -o cover.html
    $ open cover.html
