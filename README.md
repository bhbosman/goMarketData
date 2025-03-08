# GoMarketData Application


This is a CLI style application that uses golang to connect other application to collate their data into one application to view

This application connects to:
1. goLuno Application (https://github.com/bhbosman/goLuno), 
2. kraken-stream (https://github.com/bhbosman/gokraken),

 
via a TCP connection. 

# Pre-requisites
1. Golang version 1.24

# To install
1. Create a folder, for the source code
2. Clone the repo
3. Run go install



# Run application

Here is a screenshot of a running application.
![screenshot](marketdata.png)

To start the market feed, scroll to the FxServices screen, and active all the services. The application connects to the two application listed above.

