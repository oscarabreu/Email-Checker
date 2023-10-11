package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Entry point of the program.
func main() {
	// Create a new scanner to read from standard input.
	scanner := bufio.NewScanner(os.Stdin)
	// Print column headers for the results.
	fmt.Printf("domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")

	// Iterate through every line of input and check the domain details.
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	// Handle any errors encountered while reading from the input.
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

// Function to check the details of a domain.
func checkDomain(domain string) {
	// Variables to store the status and details of the domain.
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// Look up the MX records of the domain.
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	// Check if there are any MX records.
	if len(mxRecords) > 0 {
		hasMX = true
	}

	// Look up the TXT records of the domain.
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	// Check for an SPF record in the domain's TXT records.
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// Check for a DMARC record associated with the domain.
	dmarcRecordList, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error%v\n", err)
	}
	// Check for DMARC presence in the DMARC TXT records.
	for _, record := range dmarcRecordList {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	// Print the domain details.
	fmt.Printf("%v %v %v %v %v %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
