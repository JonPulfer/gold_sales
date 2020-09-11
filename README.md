# Gold Sales report

## Design decisions

The structure of the packages in this repository has been influenced by patterns
I have learned through DDD, Hexagonal, Onion implementation [such as described in here](https://edwardthienhoang.wordpress.com/2017/12/13/ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/)

I have treated the CSV file as a repository as that opens up the potential for defining various 
data sources. These could be other file types (JSON, XML, TSV), Database or an external service (API, Events, GraphQL).

The names and terms in places are probably not entirely accurate as more business knowledge would help
make them more relevant.

## Running

Tests can be run in the normal Go way: -

```
go test -v ./...
```

The main binary is in `cmd/gold_sales_report/` and can be built using something like: -

```
go build -o ./gold_sales_report cmd/gold_sales_report/main.go
```

Once built, you can just run it to accept the defaults and produce `output.csv`

Additionally, there are flags/envvars: -

```
/V/s/g/s/g/J/gold_sales_report (master|✚1…) $ ./gold_sales_report -h
Usage of ./gold_sales_report:
  -inputFilename="sample-transactions.csv": CSV File to read from
  -numMonths=6: Number of months
  -numTopSpenders=3: Number of top spenders per month
  -outputFilename="output.csv": Output filename
```
