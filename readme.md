# Proto-Order-Data
[![Go Report Card](https://goreportcard.com/badge/github.com/insan1k/proto-order-data)](https://goreportcard.com/report/github.com/insan1k/proto-order-data)

This project's objective is to implement a client for coinbase-pro that connects to the websocket feed for matches and 
reports on the Volume Weighted Average Price for the last 200 matches. 

This project has some assumptions one should take note.
* Orders model was developed using a circular list to allow for easy iteration without having to allocate resources
to roll over the list with new orders.
* I used type decimal.Decimal for calculating the average but after measuring several times using the profiler the 
 performance impact of this library is too significant to ignore, I have since pre-allocated all 32 decimal places 
 before performing the calculations, this has cut in half the cpu time used by the program but as expected this came  
 at a cost of memory consumption
* I would have preferred to separate websocket and dialer logic from exchange logic but this is not worth to do for a 
  single exchange.
  
   