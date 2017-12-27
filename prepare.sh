#!/bin/bash

pdir="prep"
sdir="${pdir}/src/hlt"
hdir="${sdir}/github.com/daved/halitego"
ifile="install.sh"
zfile="sub_halitego.zip"

mkdir -p ${hdir}

cp ./main.go ${pdir}/MyBot.go
cp -a ./vendor/* ${sdir}
cp -a ./bot ./geom ./ops ${hdir}

pushd ${pdir}

echo '#!/bin/bash' > ${ifile}
echo 'export GOPATH="${PWD}/src/hlt"' > ${ifile}
chmod +x ${ifile}

zip -r ${zfile} ./*
mv ${zfile} ../

popd

rm -rf ${pdir}
