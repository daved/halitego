# halitego

    go get -u github.com/daved/halitego

# halitego (Install Bots)

    go get -u github.com/daved/halitego/...

## Useful Aliases

    alias buildfile="echo -n builds/\$(basename \${PWD})"
    alias gobuild="go build -i -o \$(buildfile)"
    alias haliterun="halite -t -q -d" # with "'{num} {num}' {opponent}"
    
    alias hliterun="gobuild && haliterun '300 200' \$(buildfile)"
    alias hliterunsmall="gobuild && haliterun '240 160' \$(buildfile)"
    alias hliterunlarge="gobuild && haliterun '384 256' \$(buildfile)"
    alias hliteclear="rm ./*game.log ./replay*.hlt 2>/dev/null"
