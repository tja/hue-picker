<p align="center">
    <img width="256" src="etc/media/logo.png">
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/tja/hue-picker"><img src="https://goreportcard.com/badge/github.com/tja/hue-picker" alt="Go Report Card"></a>
    <a href="https://github.com/tja/hue-picker/blob/master/LICENSE"><img src="http://img.shields.io/badge/license-MIT-brightgreen.svg" alt="MIT License"></a>
</p>


# Philips Hue Color Picker

**Hue Picker** is a bare-bones web server and application that controls a single Philips Hue light.

It gives control over an individual Hue light without opening up the entire home automation network. For
instance, it allows children to manage the light in their own room without being able to mess around with all
other lights or home automation appliances.

The *Hue Picker* web application can be installed as an App on the iOS or Android home screen to simplify its
usage further.

<p align="center">
    <img width="480" src="etc/media/hero.gif">
</p>


## Installation

Pre-built binaries are available on the [release page](https://github.com/tja/hue-picker/releases/latest).
Download the binary, make it executable, and move it to a folder in your `PATH`:

```bash
curl -sSL https://github.com/tja/hue-picker/releases/download/v0.0.1/hue-picker-`uname -s`-`uname -m` >/tmp/hue-picker
chmod +x /tmp/hue-picker
sudo mv /tmp/hue-picker /usr/local/bin/hue-picker
```


## Setup

*Hue Picker* requires to be registered with the local Hue bridge. This is done as follows:

```bash
hue-picker register
```

The tool instructs you to press the button on the Hue bridge. Once the button is pressed, three pieces of
information will be printed; the Hue bridge's **Host** address, the **Bridge ID**, and the **User ID**.

The next step is to find the ID of the light that should be controlled. Using the Host address and Bridge ID,
simply do:

```bash
hue-picker list --host="192.168.0.40" --user="YDbjwv...4arRIk"
```

The tool will output the list of rooms and associated lights. Each light is prefixed with the light ID in
brackets (e.g. `[00:17:88:01:02:07:21:13-0b]`). Note down the light that *Hue Picker* should control.


## Usage

*Hue Picker* serves a web application via its built-in web server. Using the previously gathered information,
the server can be launched like this:

```bash
hue-picker serve --host="192.168.0.40" --user="YDbjwv...4arRIk" --light="00:17:...:21:13-0b"
```

Once started, the web application can be opened at http://localhost:80/ . Note that the port number and network
interface can be changed via the `--listen` parameter.

Run `hue-picker serve --help` to see the list of all available options.

### Configuration

*Hue Picker* will look for a configuration file `config.yaml` at several places in the following order:

- `/etc/hue-picker/config.yaml`
- `$HOME/.config/hue-picker/config.yaml`
- `$PWD/config.yaml`

Command line parameters and configuration file options are named the same.

Furthermore, *Hue Picker* can be configured via environment variables. Simply take the upper-cased command line
parameter and prefix with `HUE_PICKER_` &mdash; e.g. `--host` becomes `HUE_PICKER_HOST`, `--bridge` becomes
`HUE_PICKER_BRIDGE`, etc.

Command line parameters override configuration file options, which override environment variables.


## License

Copyright (c) 2022&ndash;23 Thomas Jansen. Released under the
[MIT License](https://github.com/tja/hue-picker/blob/master/LICENSE).

<a href="https://www.flaticon.com/free-icons/smart-light" title="smart light icons">Smart light icons created by Freepik - Flaticon</a>
