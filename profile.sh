go run main.go
cp /var/folders/rc/k0d6f9l55p724tb3ns1y3y600000gp/T/profile456027223/cpu.pprof .
go tool pprof --pdf main cpu.pprof > callgraph.pdf
open callgraph.pdf