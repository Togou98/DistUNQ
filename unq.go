package DistUNQ

import (
	"net"
	"os"
	"sync"
	"time"
	"errors"
	"log"  //debug
)
const(
	Timebit40 = 40
	IPbit8 = 8
	Truncpid5 = 5
	Seqnumbit11 = 11
	rationano2milli = 1e6
)
var  defaultepoch = time.Date(2020,1,1,0,0,0,0,time.UTC)
type DistUNQ struct{
	mu *sync.Mutex
	starttime int64
	elapsedtime int64 //use milliseconds to record time about 35 years
	ipbit uint8 //8 bits 32bit ip address last 8bit
	pid uint8 // 5 bits
	seqnum uint16
}

func NewUNQ(epoch time.Time) (*DistUNQ,error){
	ret := new(DistUNQ)
	ret.mu = new(sync.Mutex)
	if epoch.IsZero() {
		return nil,errors.New("time can't be zero tiem")
	}
	if epoch.After(time.Now()){
		return nil,errors.New("Maybe you are a future-man ?")
	}
	ret.starttime = epoch.UTC().UnixNano() / rationano2milli//use millisecond as time unit
	//ret.seqnum = uint16(1 << Seqbit   - 1)   //first left shift then result -1
	ret.pid = truncatepid()
	ipnum ,err := getrealip()
	if err != nil{
		return nil,err
	}
	ret.ipbit = ipnum
	return ret,nil
}
func(d *DistUNQ)NextID()(uint64,error){
	const seqmaxmask = uint16(1 << Seqnumbit11   -1) //seq plus one if overflow then seq = 0;
	d.mu.Lock()
	defer d.mu.Unlock()
	newgap := elapsedtime(d.starttime)
	if d.elapsedtime < newgap{
		d.elapsedtime = newgap
		d.seqnum = 0

	} else {
		d.seqnum = (d.seqnum + 1) & seqmaxmask
		if d.seqnum == 0 {  //overflow
			d.elapsedtime++
			//time.Sleep(time.Duration(time.Now().UTC().UnixNano()%ration2m) * time.Nanosecond )
		}
	}
	return d.iD()
}
func(d *DistUNQ)iD()(uint64,error){
	if d.elapsedtime > (1 << Timebit40) {
		return 0, errors.New("Oh, my God! It's the end of the world! ")
	}
	return uint64(d.elapsedtime)<<(IPbit8+Truncpid5+Seqnumbit11) |
		uint64(d.ipbit)<<(Truncpid5+Seqnumbit11) |
		uint64(d.pid)<<Seqnumbit11 |
		uint64(d.seqnum), nil
}
func elapsedtime(starttime int64)int64{
	return (time.Now().UTC().UnixNano() /rationano2milli) - starttime
}
func truncatepid()uint8{
	tpid := os.Getpid()
	tpid %= 32
	//fmt.Println("TPID :",tpid)
	return uint8(tpid)
}
func getrealip()(uint8,error){
	ntwk,err := net.InterfaceAddrs()
	if err != nil{
		return 0,err
	}
	for _,i := range ntwk{
		ip,ok := i.(*net.IPNet)
		if ok &&  !ip.IP.IsLoopback(){
			fmtip := ip.IP.To4()
			//fmt.Println(fmtip)
			if checkip(fmtip){
				log.Println("IP last 8 bit :",uint8(fmtip[3]))
				return uint8(fmtip[3]),nil
			}
		}
	}
	return uint8(0),errors.New("No valid IP Address")
}
func checkip(ip net.IP) bool{
	if ip == nil {
		return false
	}
	if ip[0] == 10 {
		return false
	}
	if ip[0] == 169 && ip[1] == 254{
		return false
	}
	if ip[0] == 172 && (ip[1] >= 16 && ip[1] <= 31){
		return false
	}
	if ip[0] == 192 && ip[1] == 167{
		return false
	}
	return true
}
