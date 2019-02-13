# co2mini

Simple CO2-mini library for Go.
Currently only Linux(386/amd64) is supported.

## Hardware
CO2 Monitor - CO2-mini
https://www.kk-custom.co.jp/emp/CO2-mini.html

## example

See `example/main.go`

Build
```
go build -o co2mini example/main.go
```

Run
```
sudo ./co2mini
# {"co2":0,"temp":23.2,"time":1550069838}
# {"co2":2433,"temp":23.2,"time":1550069843}
```
