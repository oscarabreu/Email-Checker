package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Constants for the SPF and DMARC record prefixes.
const (
	SPFRecordPrefix   = "v=spf1" //SPF Records start with "v=spf1"
	DMARCRecordPrefix = "v=DMARC1" // DMARC Records start with "v=DMARC1"
	DMARCDomainPrefix = "_dmarc." // DMARC Domains are prefixed with "_dmarc."
)

func main() {
	  // Define command-line flags
	  var inputFile = flag.String("input", "", "Path to the input file containing domain names.")
	  var outputFile = flag.String("output", "", "Path to the output file to write results.")
	  var workerCount = flag.Int("workers", runtime.NumCPU(), "Number of concurrent workers for domain processing.")
	  var logFile = flag.String("log", "", "Path to the log file. If not provided, logs will be printed to the console.")
	  var verbose = flag.Bool("verbose", false, "Enable verbose logging.")
	  flag.Parse()
  
	  // If logFile is provided, redirect logs there
	  if *logFile != "" {
		  f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		  if err != nil {
			  log.Fatalf("Error opening log file: %v\n", err)
		  }
		  defer f.Close()
		  log.SetOutput(f)
	  }
  
	  // Handle verbosity
	  if !*verbose {
		  log.SetOutput(io.Discard)
	  }
  
	  var input io.Reader = os.Stdin
	  var output io.Writer = os.Stdout
  
	  // If inputFile is provided, read from that file
	  if *inputFile != "" {
		  f, err := os.Open(*inputFile)
		  if err != nil {
			  log.Fatalf("Error opening input file: %v\n", err)
		  }
		  defer f.Close()
		  input = f
	  }
  
	  // If outputFile is provided, write to that file
	  if *outputFile != "" {
		  f, err := os.Create(*outputFile)
		  if err != nil {
			  log.Fatalf("Error opening output file: %v\n", err)
		  }
		  defer f.Close()
		  output = f
	  }
	// We use a Scanner to buffer the input stream and read line by line
	// This prevents loading the entire input into memory at once!
	// However, if domain line is >64k, it will return false but this is rare (I hope!)
	scanner := bufio.NewScanner(input)
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	// We introduce concurrency to speed up the process
	var wg sync.WaitGroup // WaitGroups help us wait for all goroutines to finish
	var numWorkers int
	if *workerCount > 0 {
    	numWorkers = *workerCount
	} else {
    numWorkers = runtime.NumCPU() // Default is to use the number of CPUs as the number of workers
}

	// I was advised to make an interesting decision made here & choosing the |jobs channel| =
	// the number of workers/number of CPUs. Buffering the jobs channel to numWorkers size 
	// allows for smooth job dispatch. It lets you send multiple jobs without waiting, 
	// optimizing flow between dispatching and processing. It strikes a balance between
	//  memory use and avoiding excessive blocking in the main routine.
	
	jobs := make(chan string, numWorkers) // Buffered channel to send jobs to workers

	// We use an unbuffered channel for results because we want to process them as they come in.
	results := make(chan string) // Buffered cannel to receive results from workers

	// Loops through the number of workers and starts a goroutine for each.
	for w := 0; w < numWorkers; w++ {
        go worker(jobs, results, &wg)
    }

	// This allows for asynchronous processing of the results.
	// The main program will not hang, each result is printed as it is received.
	// Additionally, this anonymous function has access to all variables in the main function.
	// So we don't need to pass anything in!
	go func() {
        for r := range results {
            fmt.Print(output, r)
        }
    }()

	// Scan each line of the input and send it to the jobs channel.
	for scanner.Scan() {
        domain := scanner.Text() // Read the domain from the input
        wg.Add(1) // Increment the WaitGroup counter
        jobs <- domain // Send the domain to the jobs channel
    }

	close(jobs) // Close the jobs channel
    wg.Wait() // Wait for all workers to finish
    close(results) // Close the results channel

	// Error handling
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}

}
// worker is a function that takes a domain from the jobs channel, 
// inspects it, and sends the result to the results channel
func worker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
for domain := range jobs {
        results <- inspectDomain(domain)
        wg.Done()
    }
}

// inspectDomain takes a domain name and prints a CSV line with the results
func inspectDomain(domain string) string {
	
	isMXPresent := detectMX(domain) // Checks if MX records exist
	isSPFPresent, detectedSPF := detectSPF(domain) // Checks if SPF records exist and returns the record
	isDMARCPresent, detectedDMARC := detectDMARC(domain) // Checks if DMARC records exist and returns the record

	return fmt.Sprintf("\nDomain: %v, \nisMXPresent: %v, \nisSPFPresent: %v, \nSPF: %q, \nisDMARCPresent: %v, \ndMARC: %q\n", domain, isMXPresent, isSPFPresent, detectedSPF, isDMARCPresent, detectedDMARC)
}

// Function to validate MX records using LookupMX from the net package.
// Returns true if MX records exist (more than 0), false if not.
func detectMX(domainName string) bool {
	mxEntries, mxErr := net.LookupMX(domainName)
	// Error handling.
	if mxErr != nil {
		log.Printf("MX Lookup Error: %v\n", mxErr)
		return false
	}
	return len(mxEntries) > 0
}

// Function to validate SPF records using LookupTXT from the net package.
// Returns true if SPF records exist (more than 0), false if not.
func detectSPF(domainName string) (bool, string) {
	txtEntries, txtErr := net.LookupTXT(domainName)
	// Error handling.
	if txtErr != nil {
		log.Printf("TXT Lookup Error: %v\n", txtErr)
		return false, ""
	}
	// Loop through the TXT entries and check if any of them start with "v=spf1"
	for _, entry := range txtEntries {
		if strings.HasPrefix(entry, SPFRecordPrefix) {
			return true, entry
		}
	}
	return false, ""
}

// Function to validate DMARC records using LookupTXT from the net package.
// Returns true if SPF records exist (more than 0), false if not.
func detectDMARC(domainName string) (bool, string) {
	// To check the DMARC record for example.com, 
	// you'd look up the TXT records for _dmarc.example.com.
	dmarcEntries, dmarcErr := net.LookupTXT(DMARCDomainPrefix + domainName)
	// Error handling.
	if dmarcErr != nil {
		log.Printf("DMARC Lookup Error: %v\n", dmarcErr)
		return false, ""
	}
	// Loop through the TXT entries and check if any of them start with "v=DMARC1"
	for _, entry := range dmarcEntries {
		if strings.HasPrefix(entry, DMARCRecordPrefix) {
			return true, entry
		}
	}
	return false, ""
}
