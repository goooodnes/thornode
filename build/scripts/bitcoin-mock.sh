#!/bin/sh

SIGNER_NAME="${SIGNER_NAME:=thorchain}"
SIGNER_PASSWD="${SIGNER_PASSWD:=password}"
MASTER_ADDR="${BTC_MASTER_ADDR:=bcrt1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynawhcf2xa}"
BLOCK_TIME=5

bitcoind -regtest -rpcuser=$SIGNER_NAME -rpcpassword=$SIGNER_PASSWD -rpcallowip=0.0.0.0/0 -rpcbind=127.0.0.1 -rpcbind=$(hostname) &

# give time to bitcoind to start
sleep 5
bitcoin-cli -regtest -rpcuser=$SIGNER_NAME -rpcpassword=$SIGNER_PASSWD generatetoaddress 500 $MASTER_ADDR

# mine a new block every BLOCK_TIME
while true
do
	bitcoin-cli -regtest -rpcuser=$SIGNER_NAME -rpcpassword=$SIGNER_PASSWD generatetoaddress 1 $MASTER_ADDR
	sleep $BLOCK_TIME
done
