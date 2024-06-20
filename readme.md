# Description
This program helps you parse cron expression and expands each part to show the times when it will run and prints the output to the terminal with the field name taking 14 columns 

## How to run
This program is written in golang. You need to have `go` installed to run this program. Follow the below instruction to install go in mac os.
```
brew install golang
```
Once installed, you can check the version with the below command
```
go version
```

Next is to install dependencies required for this program. Run the below command to install dependencies.
```
go mod vendor
```

Once the modules installed, you can run the program with the below command.
```
go run main.go "5 0 * 8 * /usr/bin/find"
```

The above command will print the below ouput to the terminal
```
minute         5
hour           0
day of month   1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month          8
day of week    0 1 2 3 4 5 6
command        /usr/bin/find
```


## Running tests
The below command runs test. 
```
go test ./...
```
