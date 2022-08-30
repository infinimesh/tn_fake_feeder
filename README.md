# tn_fake_feeder

[![Release tool and DB](https://github.com/infinimesh/tn_fake_feeder/actions/workflows/release.yml/badge.svg)](https://github.com/infinimesh/tn_fake_feeder/actions/workflows/release.yml)

internal tool to feed fake data to platform

## Installation

1. Navigate to [Releases page](https://github.com/infinimesh/tn_fake_feeder/releases)
2. Find an Archive for your OS and CPU Arch
3. Create directory like `tn_feeder` and navigate into it ( like `mkdir tn_feeder && cd tn_feeder`)
4. Download Archive into that directory (e.g. `wget https://github.com/infinimesh/tn_fake_feeder/releases/download/v0.0.0/tn-feeder-v0.0.0-linux-amd64.tar.gz`)
5. Unpack it(e.g. `tar xvf tn-feeder-v0.0.0-linux-amd64.tar.gz`)
    You must see three files:
     - This README
     - track.db - database with real world roads waypoints
     - tn-feeder - executable file
6. Congrats, You are ready to fake!

## Configuration and running feeder

1. Perform login via [`inf`](https://github.com/infinimesh/inf)(if you don't have it installed, get it, that's way easier than this)
   `inf login api.your.infinimesh:8000 user password`
2. Pick/Create Namespace of your choice and copy its UUID (like `aaaaaa-bbbb-cccc-dddddd` kind of thing)
3. Run `./tn-feeder <namespace-uuid> <amount-of-trucks>`
    Amount of trucks stands for number of devices `tn_feeder` will create and simulate
4. Press `Ctrl+C` once you want to stop. Programm won't exit immediately(may take up to 15 sec), don't worry, it's cleaning up the devices from `infinimesh`, let it finish.
