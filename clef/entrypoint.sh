#!/bin/sh

ACTION="${1:-run}"
CHAINID="${CLEF_CHAINID:-1440}"
DATA=/app/data

echo "Running action '${ACTION}'"

init() {
    parse_json() { echo $1|sed -e 's/[{}]/''/g'|sed -e 's/", "/'\",\"'/g'|sed -e 's/" ,"/'\",\"'/g'|sed -e 's/" , "/'\",\"'/g'|sed -e 's/","/'\"---SEPERATOR---\"'/g'|awk -F=':' -v RS='---SEPERATOR---' "\$1~/\"$2\"/ {print}"|sed -e "s/\"$2\"://"|tr -d "\n\t"|sed -e 's/\\"/"/g'|sed -e 's/\\\\/\\/g'|sed -e 's/^[ \t]*//g'|sed -e 's/^"//' -e 's/"$//' ; }
    if [ ! -f /app/config/password ]; then
        < /dev/urandom tr -dc _A-Z-a-z-0-9 2> /dev/null | head -c32 > /app/config/password
    fi
    SECRET=$(cat /app/config/password)
    echo "Using password: '${SECRET}'"
    /usr/local/bin/clef --configdir "$DATA" --stdio-ui init 2>&1 << EOF
$SECRET
$SECRET
EOF
    if [ "$(ls -A "$DATA"/keystore 2> /dev/null)" = "" ]; then
        /usr/local/bin/clef --keystore "$DATA"/keystore --stdio-ui newaccount --lightkdf 2>&1 << EOF
$SECRET
EOF
    fi
    ls -al "$DATA"/keystore
    /usr/local/bin/clef --keystore "$DATA"/keystore --configdir "$DATA" --stdio-ui setpw 0x"$(parse_json "$(cat "$DATA"/keystore/*)" address)" 2>&1 << EOF
$SECRET
$SECRET
$SECRET
EOF
    /usr/local/bin/clef --keystore "$DATA"/keystore --configdir "$DATA" --stdio-ui attest "$(sha256sum /app/config/rules.js | cut -d' ' -f1 | tr -d '\n')" 2>&1 << EOF
$SECRET
EOF
}

run() {
    SECRET=$(cat /app/config/password)
    echo "Using password: ${SECRET}"
    rm /tmp/stdin /tmp/stdout || true
    mkfifo /tmp/stdin /tmp/stdout
    (
    exec 3>/tmp/stdin
    while read < /tmp/stdout
    do
        if [[ "$REPLY" =~ "enter the password" ]]; then
            echo '{ "jsonrpc": "2.0", "id":1, "result": { "text":"'"$SECRET"'" } }' > /tmp/stdin
            break
        fi
    done
    ) &
    /usr/local/bin/clef --stdio-ui --keystore "$DATA"/keystore --configdir "$DATA" --chainid "$CHAINID" --http --http.addr 0.0.0.0 --http.port 8550 --http.vhosts "*" --rules /app/config/rules.js --nousb --lightkdf --ipcdisable --4bytedb-custom /app/config/4byte.json --pcscdpath "" --auditlog "" --loglevel 3 < /tmp/stdin | tee /tmp/stdout
}

full() {
    if [ ! -f "$DATA"/masterseed.json ]; then
        init
    fi
    run
}

$ACTION