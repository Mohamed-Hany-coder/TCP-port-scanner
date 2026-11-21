package scanner

import (
 "context"
 "fmt"
 "net"
 "sort"
 "strconv"
 "sync"
 "time"
)

type Status string
const ( StatusOpen Status = "open"; StatusClosed Status = "closed"; StatusError Status = "error" )
type Result struct { Host string `json:"host"`; Port int `json:"port"`; Status Status `json:"status"`; Error string `json:"error,omitempty"` }
func CheckPort(ctx context.Context, host string, port int, timeout time.Duration) Result { r:=Result{Host:host,Port:port,Status:StatusClosed}; if timeout<=0 {r.Status=StatusError;r.Error="timeout must be greater than zero";return r}; c,err:=(&net.Dialer{Timeout:timeout}).DialContext(ctx,"tcp",net.JoinHostPort(host,strconv.Itoa(port))); if err==nil {_=c.Close();r.Status=StatusOpen;return r}; if ctx.Err()!=nil {r.Status=StatusError;r.Error=ctx.Err().Error();return r};r.Error=err.Error();return r }
func Scan(ctx context.Context,host string,ports []int,timeout time.Duration) []Result { out:=make([]Result,0,len(ports));for _,p:=range ports{if ctx.Err()!=nil{break};out=append(out,CheckPort(ctx,host,p,timeout))};return out }
func ScanConcurrent(ctx context.Context,host string,ports []int,timeout time.Duration,workers int)([]Result,error){if workers<=0{return nil,fmt.Errorf("workers must be greater than zero")};if len(ports)==0{return []Result{},nil};if workers>len(ports){workers=len(ports)};jobs:=make(chan int);results:=make(chan Result,len(ports));var wg sync.WaitGroup;wg.Add(workers);for i:=0;i<workers;i++{go func(){defer wg.Done();for{select{case<-ctx.Done():return;case p,ok:=<-jobs:if !ok{return};results<-CheckPort(ctx,host,p,timeout)}}}()};go func(){defer close(jobs);for _,p:=range ports{select{case<-ctx.Done():return;case jobs<-p:}}}();go func(){wg.Wait();close(results)}();out:=make([]Result,0,len(ports));for r:=range results{out=append(out,r)};sort.Slice(out,func(i,j int)bool{return out[i].Port<out[j].Port});return out,nil}
