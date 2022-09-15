#!/bin/bash

NODE_WALLET_PASS_FILE=$VEGA_HOME/node-vega-wallet-pass.txt
ETH_WALLET_PASS_FILE=$VEGA_HOME/ethereum-vega-wallet-pass.txt

NODE_RPC_ADDR=localhost:3002
NODE_REST_ADDR=localhost:3003
DATA_NODE_REST_ADDR=localhost:3029
CLEF_ADDR=127.0.0.1:8550
WALLET_NAME=created-wallet
VEGA_WALLET_PASS_FILE=$VEGA_HOME/vega-wallet-pass.txt
VEGA_HOME=~/.vegacapsule/testnet/vega/node0


VEGA_PUB_KEY=$(vega wallet --home $VEGA_HOME key list --wallet $WALLET_NAME --passphrase-file $VEGA_WALLET_PASS_FILE --output json | jq -r '.keys[0].publicKey')

choosen_validator=$(curl -s http://$NODE_REST_ADDR/validators | jq -r '.validators[] | select(.vegaPubKey=="'$VEGA_PUB_KEY'")')

NODE_ID=$(echo $choosen_validator | jq -r .nodeId)
ETH_ADDR=$(echo $choosen_validator | jq -r .ethereumAddress)
current_block=$(curl -s http://$NODE_REST_ADDR/statistics | jq -r .statistics.blockHeight)
target_block=$(expr $current_block + 20)

if ($USE_CLEF)
then
    echo 'Generating Clef Ethereum wallet';
    ETH_NEW_ADDR=$(curl -s -H 'Content-Type: application/json' -X POST --data '{"id": 1, "jsonrpc": "2.0", "method": "account_new"}' http://$CLEF_ADDR | jq -r .result)
    vega nodewallet --home $VEGA_HOME --passphrase-file $NODE_WALLET_PASS_FILE import -c ethereum --eth.clef-address http://$CLEF_ADDR  --clef-account-address $ETH_NEW_ADDR --force
else
    echo "Generating regular Ethereum wallet";
    ETH_NEW_ADDR=$(vega nodewallet --home $VEGA_HOME --passphrase-file $NODE_WALLET_PASS_FILE generate --chain ethereum --wallet-passphrase-file $ETH_WALLET_PASS_FILE --force --output json | jq -r .walletFilePath | sed 's/.*--/0x/')
fi


echo NODE_ID "$NODE_ID"
echo VEGA_PUB_KEY "$VEGA_PUB_KEY"
echo ETH_ADDR "$ETH_ADDR"
echo ETH_NEW_ADDR "$ETH_NEW_ADDR"
echo target_block $target_block

vega wallet --home $VEGA_HOME command send --wallet $WALLET_NAME --passphrase-file $VEGA_WALLET_PASS_FILE --node-address $NODE_RPC_ADDR --pubkey $VEGA_PUB_KEY '{"ethereumKeyRotateSubmission": {"targetBlock": "'$target_block'", "newAddress": "'$ETH_NEW_ADDR'", "currentAddress": "'$ETH_ADDR'", "submitterAddress":""}}'

while true; do
sleep 1;
current_block=$(curl -s http://$NODE_REST_ADDR/statistics | jq -r .statistics.blockHeight);
blocks_left=$(expr $target_block - $current_block);
echo "Wating for target block: $target_block, current block: $current_block";
if [[ $blocks_left -lt 5 ]];
then vega nodewallet --home $VEGA_HOME reload -c ethereum; break;
fi;
done

sleep 4
echo "/api/v2/vega/keys/ethereum/rotations response"
curl -s $DATA_NODE_REST_ADDR/api/v2/vega/keys/ethereum/rotations
echo ""
echo "/api/v2/erc20/multisigcontrol/signer/added/bundle response:"
curl -s $DATA_NODE_REST_ADDR/api/v2/erc20/multisigcontrol/signer/added/bundles
echo ""
echo "/api/v2/erc20/multisigcontrol/signer/removed/bundle response:\n"
curl -s $DATA_NODE_REST_ADDR/api/v2/erc20/multisigcontrol/signer/removed/bundles

