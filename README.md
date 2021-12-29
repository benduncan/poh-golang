# Proof of History - Concepts in Go

```
 ________  ________  ___  ___         
|\   __  \|\   __  \|\  \|\  \        
\ \  \|\  \ \  \|\  \ \  \\\  \       
 \ \   ____\ \  \\\  \ \   __  \      
  \ \  \___|\ \  \\\  \ \  \ \  \ ___ 
   \ \__\    \ \_______\ \__\ \__\\__\
    \|__|     \|_______|\|__|\|__\|__| go!
                           
```

## Objective

Provide source-code concepts from the [Solana whitepaper](https://github.com/solana-labs/whitepaper/blob/master/solana-whitepaper-en.pdf) in Golang to understand proof-of-history and the mechanics behind Solana further.

Used as an exercise to understand Proof of History further by writing a simple implementation as described in the whitepaper. 

## Usage

```
make build
```

Build from source

```
./bin/poh-golang
```

Run a proof-of-history to generate the hashrate on a single core, and validate the results on all available CPU cores

```
CPU Cores 8
Generate Hashrate 1482625 p/sec (1-core)
Verify Hashrate 4701530 p/sec (8-cores)
Verify Hashrate 587691 p/core
```

## Benchmarks

```
make bench
```

Simulate benchmarks to generate proof-of-history (POH) using a single core, and verify the output on multiple cores.

## Roadmap

- [X] Calculate hashrate for POH generation (single core)
- [X] Calculate hashrate for POH validation (multiple cores)
- [ ] Implement Solana hash table indexes for user addresses and packet payloads described in the whitepaper.
- [ ] Add GPU support via Cuda implementation of SHA256
- [ ] Synchronize multiple POH validators on the local network
- [ ] Implement event signatures
- [ ] Add basic Proof of Stake implementation as per the Solana whitepaper
- [ ] Implement vote support for PoS implementation
- [ ] Implement election for PoS and simulate PoH generator failure
- [ ] Add PoH election support and secondary validator promotion
- [ ] Implement Streaming proof of Replication using basic CBC encryption and XOR inputs