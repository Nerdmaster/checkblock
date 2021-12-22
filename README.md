# checkblock

External app to check and cache the latest block from DMO nodes.

This can be used to ease the burden on a node that has many computers sending
data to it. You will have to update your miner to check the file instead of
reading the network, but file reads are *much* faster than RPC calls.
