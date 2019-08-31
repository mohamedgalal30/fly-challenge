## Start the App
- To start the app just run ```docker-compose up```
- After docker get everything ready, you can retrieve the payment transaction data from api ``` http://localhost:3000/api/payment/transaction```
- As you see in the url, the app will be bounded on port 3000, so you need to get that port available.
## Testing
To test the app just run ```docker exec -it fly-challenge go test ./... -v```
## Benchmark
To show the benchmark just run ```docker exec -it fly-challenge go test ./search/ -v -run="none" -bench="BenchmarkSearchRun" -benchtime="3s" -benchmem```
## Note
In case you want to add new provider, all you have to do is
- add new file in providers package, define its scheme and its logic.
- add provider data files in data directory and refer to its path in the provider file.
- And of course write unit test for it. just use functions in providers_test package in your unit test(don't write test logic again)
