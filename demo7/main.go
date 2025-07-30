package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(111)

	for scanner.Scan() {
		CheckDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("1")
	}
}

func CheckDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, damrcRecord string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Fatal("2")
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecord, err := net.LookupTXT(domain)
	if err != nil {
		log.Println("3")
	}
	for _, record := range txtRecord {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Println("4")
	}
	for _, record1 := range dmarcRecords {
		if strings.HasPrefix(record1, "v=DMARC1") {
			hasDMARC = true
			damrcRecord = record1
			break
		}
	}

	fmt.Println(hasMX, hasSPF, hasDMARC, spfRecord, damrcRecord)
}
