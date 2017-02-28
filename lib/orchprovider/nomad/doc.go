// Package nomad provides Nomad implementation of orchestration provider
// that aligns by the interfaces suggested by mayaserver's orchprovider.
//
// This package primarily consists of below files:
// 1. nomad.go
// 2. api.go
// 3. client.go
//
// The dependencies can be depicted as shown below:
//
//    nomad   ==    aligns with   ==>   orchprovider
//
//    nomad   ==    depends on    ==>   api
//    nomad   ==    depends on    ==>   client
//    api     ==    depends on    ==>   client
package nomad
