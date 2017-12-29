#!/bin/bash

pdir="prep"
sdir="${pdir}/src"
hdir="${sdir}/github.com/daved/halitego"
mfile="MyBot.go"
zfile="sub_halitego.zip"

rm -rf ${pdir}
mkdir -p ${hdir}

cp ./cmd/gopherbot/main.go ${pdir}/${mfile}
cp -a ./vendor/* ${sdir}
cp -a ./geom ./internal/* ./ops ${hdir}

pushd ${pdir} > /dev/null

sed -i 's#/internal/#/#g' ${mfile}

zip -r ${zfile} ./* > /dev/null
mv ${zfile} ../

popd > /dev/null

rm -rf ${pdir}
