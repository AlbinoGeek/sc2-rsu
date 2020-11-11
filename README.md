# sc2-rsu

**Unofficial SC2ReplayStats Uploader** by AlbinoGeek

## Features

- Cross-Platform: Full Linux and Windows support, Mac OS X untested
- Bronze Friendly: `login` command accepts an API Key or Email Address
- Optimized Build: Requires no vespene, and less than 10MB of RAM

Missing a feature? [Request it!](/issues/new)

## Usage

1. [Download a Release](/AlbinoGeek/sc2-rsu/releases) or [Get the Sources](#building-from-source)
2. Issue the `login` command to perform one-time setup (get your SC2ReplayStats API Key)
3. Run the program without a command to automatically upload games as they are played

```
$ sc2-rsu login john@doe.com
  Password for sc2ReplayStats account john@doe.com: 
  Success! Logged in to account #9001
  API Key set in configuration!

$ sc2-rsu
  Starting Automatic Replay Uploader...
  Ready!
```

For full usage instructions, consult the `--help` output.

## Building From Source

1. Install Golang
2. Clone the repository
3. Run `make`

```
# Install Golang [apt / yum / dnf install golang]
$ git clone https://github.com/AlbinoGeek/sc2-rsu
$ cd sc2-rsu
$ make all
```

The binary will be built into the `_dist` directory.
