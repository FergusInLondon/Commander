# Commander

Commander is a simple `systemd` enabled daemon written in Go, which listens on a Unix Domain Socket for JSON payloads expressing a *Command*. Commands are simply structs in Go which adhere to a specific interface and perform some lower-level action. 

Use cases for Commands include:

 - Interfacing with `dbus` for the control of a network interface;
 - Generating configuration files via template processing;
 - Interacting with `systemd` units;
 - Updating `iptables` rules.

### Rationale

This is a barebones proof of concept at the moment, but I'm hoping it'll lead to me realising a project that I've had on the backburner for around 18 months - of having a Raspberry Pi configured as a small little secured router.

Rather than manually using SSH to adjust configuration files though, I'd like to have one that's usable and easy enough for a novice to use - so naturally a web interface is required. Unfortunately I've seen a few solutions which rely upon hacky PHP applications that are littered with `exec()` commands and other nasties; so I thought I'd take a different approach and try and architect it *properly*.

With this in mind, I can write a small PHP layer which glues the UI with Commander - the underlying system management layer. The web interface needs to know little - or nothing - about the intricacies of configuring `dnsmasq`, `hostapd`, `openvpn`, `iptables`, or any of the other underlying components. All the UI layer needs to do is to be able to tell Commander what the desired end-result is.


## Status

This is currently very much a pet project, and one that's been written over the course of a few hours.. although, it's currently capable of listening for commands and executing them; it's still a bit rough around the edges.

Currently implemented commands include: (a) a demonstration echo command, (b) a notification one - demonstrating the ability to interface with `dbus`, and (c) a template one - demonstrating configuration file generation.

I aim to write two more demonstration commands: (1) an iptables interface, and (2) a `systemd` unit manager. After that I'll do some refactoring and clean up.

## System Configuration & User Privileges

This interacts with `systemd` to restart specific units, therefore beware that - depending upon your system - this may well need to be ran as root. Not cool, but not a priority to fix at the moment either.

Similarly, assumptions are made about systemd unit names - i.e `hostapd.service`. I'll try and document these, and throw in some accompanying provisioning scripts (i.e via *Ansible*) when this is done.

## Testing

There are no tests. :( I will aim to write some tests once I decide exactly how to structure them (there's a lot of moving parts.) and I'm happy with the overall design. It would be jumping the gun to have gone full TDD on a simple proof of concept!

### Manual Testing

#### Bridging TCP connections to the Unix Socket (via Socat)

[Socat](http://www.dest-unreach.org/socat/doc/socat.html) is a handy little tool that can expose the Unix socket used by Commander, to do so simply run:

    $ socat -d -d TCP4-LISTEN:8080,fork UNIX-CONNECT:/tmp/commander.sock
    >  2017/11/02 10:21:39 socat[9455] N listening on AF=2 0.0.0.0:8080

#### Example Requests with [Postman](https://www.getpostman.com)

In the `examples` directory there's a .json file which can be imported in to *Postman*, allowing you to try out the default demo Commands.

## Thanks

It's worth mentioning the quality of the [CoreOS](https://github.com/coreos) Go libraries; the `systemd` integration is managed via one of these libraries, and I aim to use [another one](https://github.com/coreos/go-iptables) for `iptables` integration.

## License (MIT)

Copyright 2017 fergus@fergus.london <https://fergus.london>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.