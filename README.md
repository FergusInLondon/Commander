# Commander - A Work in Progress


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

#### [Pre-Refactor Tag](https://github.com/FergusInLondon/Commander/tree/46cbe22e40a9c4bb9a27804dc5eab70709ee4a6a)
This actually had a working proof of concept, but it had some very ugly logic and was overly complex.

## License (MIT)

Copyright 2017 fergus@fergus.london <https://fergus.london>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.