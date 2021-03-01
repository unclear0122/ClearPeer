# ClearPeer Tool

ClearPeer is a tool for disconnecting the remote peers from the blockchain daemon using the `disconnectnode` command.

This is a primitive approach and some parameters are hardcoded, they can be changed in the code before building.
The command execution is done via OS command intensionally to be sure that code is executed on the host with access to the daemon dataDir.

The benefit of using this program is allowing to release some connection slots and serve not fully synced nodes as a block source without having to re-start the host daemon.

The logic of the program is:
1. Check if the host node is fully synced (headers = blocks)
2. Check if the number of connections is above the `minConnections` threshold
3. Check connected peers and disconnect it if all the below conditions are met:
 - the connection is incoming (`inbound`)
 - the remote peer is at the same block height with the host
 - the total number of disconnected peers is below the `maxDisconnect` threshold
 - the peer is randomly selected for the disconnecting

Disconnected peers are not getting any ban scores and are not limited to re-establish a connection to the host.

## Limitations

Pre-build binary hardcoded to 
- work with `raven-cli` and should be placed to the same directory with `raven-cli` binary
- has minimum 110 connected peers to trigger any disconnection
- disconnect maximum 32 peers at once

Compiled as self-containing binary and tested with Ubuntu 18.04 x64 LTS, Ubuntu 20.04 LTS x64, but also should work with other Linux platforms.

Requires Go 1.15 to be installed to build from sources.

## License
The MIT License (MIT)
Copyright (c) 2021 unclear0122, https://solus.cryptoscope.io