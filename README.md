# BedrockScanner

A tool for scanning Minecraft Bedrock servers across specified IP ranges or ASNs.

## Components

### 1. Scanner
The primary scanning utility that probes Minecraft Bedrock servers within specified IP ranges.

#### Usage:
```bash
./bedrockscanner \
    -what <target> \                     # What to scan (single subnet, file with subnets, or ALL; default: ALL)
    -packets-per-second <rate> \         # Packets per second to send (default: 5000)
    -write-to-file <filename> \          # Output file for results (optional)
    -num-sockets <count>                 # Number of sockets to use (default: 1)
```

#### Target options for `-what`:
- A single subnet (e.g., `192.168.1.0/24`)
- A file containing subnets (one per line)
- `ALL` to scan the entire IPv4 range (0.0.0.0 - 255.255.255.255)

### 2. ASN to Ranges Converter (`cmd/asn_to_ranges`)
A helper utility that converts an Autonomous System Number (ASN) to its corresponding IP ranges.

#### Usage:
```bash
./asn_to_ranges \
    -as <ASN> \             # Autonomous System Number to query
    -save <filename>        # Output file (default: as.txt)
```

## Requirements
- Go 1.24 or later
- Network connectivity with ability to send/receive UDP packets

## Building
Compile both tools with:
```bash
go build -o bedrockscanner
go build -o asn_to_ranges cmd/asn_to_ranges.go
```

## Typical Workflow
1. Identify ASNs of interest
2. Convert ASNs to IP ranges using `asn_to_ranges`
3. Scan the resulting ranges with the scanner
4. Process results with custom tools

The scanner outputs results in a SQLite table.
![image](https://github.com/user-attachments/assets/4d760781-3dac-472c-93cc-4e0bcb7902e0)

## Notes
- The scanner respects the specified packets-per-second rate across all sockets
- When scanning large ranges, consider using lower PPS values or multiple machines
- Results include successfully responding Bedrock servers with their basic information
