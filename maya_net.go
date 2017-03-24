package main

import (
             "flag"
             "io"
             "strings"           
             "net/http"
)

type Status int

const (
      UNCHECKED Status = iota
      DOWN
      UP
)


// Meta contains the meta-options and functionality that nearly every
// Maya server command inherits.
type Meta struct 
{
  Ui cli.Ui

  // Whether to not-colorize output
               noColor bool
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet 
{
     f := flag.NewFlagSet(n, flag.ContinueOnError)

  // FlagSetClient is used to enable the settings for specifying
  // client connectivity options.
   if fs&FlagSetClient != 0 
   {
     f.BoolVar(&m.noColor, "no-color", false, "")
   }

  // Create an io.Writer that writes to our UI properly for errors.
  // This is kind of a hack, but it does the job. Basically: create
 // a pipe, use a scanner to break it into lines, and output each line
 // to the UI. Do this forever.
     errR, errW := io.Pipe()
     errScanner := bufio.NewScanner(errR)
     go func()
     {
       for errScanner.Scan() 
       {
         m.Ui.Error(errScanner.Text())
       } 
     }()

     f.SetOutput(errW)

     return f
}

  func (m *Meta) Colorize() *colorstring.Colorize 
  {
    return &colorstring.Colorize{
    Colors: colorstring.DefaultColors,
    Disable: m.noColor,
    Reset: true,
   }
 }

 // The Site struct encapsulates the details about the site being monitored.
 type Site struct
 {
    url string
    last_status Status
 }

 // Site.Status makes a GET request to a given URL and checks whether or not the
 // resulting status code is 200.
  func (s Site) Status() (Status, error) 
  {
    resp, err := http.Get(s.url)
    status := s.last_status

    if (err == nil) && (resp.StatusCode == 200) 
    {
        status = UP
    }
    else
    {
        status = DOWN
    }

    return status, err
}
