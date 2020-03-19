#  this package it's a simple useful distributed id generator  

##  install  
	`go get install github.com/togou98/DistUNQ`
	
**in $GOROOT/src you can see that code files  
	![pic](https://github.com/Togou98/DistUNQ/blob/master/bitpattern.png)	
## UNQ generate 64bit integer  
		* 40bit length of timebit 40bits
		* 8 bit length of machine ip address last 8bits
		* 5 bit length of currency process pid%32 so it at least need 5bits
		* 11 bit length of sequence number MAX(2047) 
