# goBACnet 
[![Build Status](https://travis-ci.com/alexbeltran/gobacnet.svg?token=pGXqDCNsbwsP7nyfos9q&branch=master)](https://travis-ci.com/alexbeltran/gobacnet)
[![GoDoc](https://godoc.org/github.com/alexbeltran/gobacnet?status.svg)](https://godoc.org/github.com/alexbeltran/gobacnet)

gobacnet is a client for bacnet written exclusively with go. The goal is to
only offer a client and a test server application.

**NOTE:** This code is very experimental and therefore should not be used in
anything you want working. There are many changes being made and I cannot
guarantee compatibility between versions. Many features of the protocol are
missing and will be added overtime. 

# Installation
Part of this library is a command line client. To install from source run:
```
go get -u github.com/alexbeltran/gobacnet/baccli
```

For usage run:

```
baccli --help
```

# Features
Below are features that are intended to be completed for the library and for
the command line interface ordered in no particular manner. If you would like
a feature not in the list below add a discussion on github for it.

## Library
- [x] Who Is
- [x] Read Property
- [x] Read Multiple Property
- [ ] Read Range
- [ ] Write Property
- [ ] Write Property Multiple
- [ ] Who Has
- [ ] Change of Value Notification
- [ ] Event Notification
- [ ] Subscribe Change of Value
- [ ] Atomic Read File
- [ ] Atomic Write File

## Command Line Interface
- [x] Who Is
- [x] Read Property
- [x] Read Multiple Property
- [ ] Read Range
- [ ] Write Property
- [ ] Write Property Multiple
- [ ] Who Has
- [ ] Atomic Read File
- [ ] Atomic Write File

# Contributing
Contributions are more then welcome for this project. Use golint for
formatting and be sure to include test coverage on any new additions. 

# License
This library is heavily based on the BACnet-Stack library originally written by
Steve Karg and therefore is released under the same license as his project.
This includes the exception which allows for this library to be linked by
proprietary code without that code becoming GPL. This exception was taken
from the original BACnet stack and is included in every file.

The exception is as follows:
```
    "As a special exception, if other files instantiate
     templates or use macros or inline functions from
     this file, or you compile this file and link it
     with other works to produce a work based on this file,
     this file does not by itself cause the resulting work
     to be covered by the GNU General Public License.
     However the source code for this file must still be
     made available in accordance with section (3) of the
     GNU General Public License."
```
