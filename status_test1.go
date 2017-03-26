package main

import (
      "fmt"
      "net"
      "github.com/mitchell/cli"
      "net/url"
      "net/https"
      "host"
      "addrs"
 package cmd ()     
 type configuration struct
 {
  URL []String
  Mem string
 }
   "fmt"
   "encoding/json"
   "os"
   "github.com/openebs/mayaserver/lib/config"
   "github.com/openebs/mayaserver/lib/server"
   
func config(conf.json) 
{
  file,_ := os.Open("conf.json")
  decoder := json.newDecoder(file)
  configuration := configuration 
  {
   err:=decoder.Decode(&configuration)
  }
  if err!=nil
  {
   fmt.println("IPaddrs not found:",IPaddrs not found)
  }
  fmt.println(configuration.Users)
});

func main()
{
 s:= "IPaddress://github.com/user:status@host.com/path?K=v#f"
  u,err :=url.Parse(s)
  if err!=nil
  {
    panic(err)
  }
  // prints first the IP address
  fmt.println(U.host)
   host,_ := os.Hostname()
   addrs,_ := net.Lookup(host)
   for_,addr := rangeaddrs
   {
     if ipv4 := addr.To4();
     ipv4 != nil
     {
       fmt.println("Ipv4:",ipv4)
     }
   }
   fmt.println(U.scheme) //prints the domain name
   fmt.println(U.user) // prints the user name
   fmt.println(U.user.username())
   P,_ := U.user.password()
   fmt.println(P) //here if you give the password then your account will be opened in the url
   
   fmt.println(host)
   fmt.println(port)
   
   fmt.println(U.path)
   fmt.println(U.fragment)
   
   fmt.println(U.RawQuery) // searches in the database and accepts your username and password
   M,_ := Url.ParseQuery(U.RawQuery)
   fmt.println(M) //print sthen mayaserver
   fmt.println(M["K"][0])
 }
