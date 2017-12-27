#!/bin/bash

pdir="prep"
sdir="${pdir}/src"
hdir="${sdir}/github.com/daved/halitego"
zfile="sub_halitego.zip"

rm -rf ${pdir}
mkdir -p ${hdir}

cp ./main.go ${pdir}/MyBot.go
cp -a ./vendor/* ${sdir}
cp -a ./bot ./geom ./ops ${hdir}

pushd ${pdir}

zip -r ${zfile} ./*
mv ${zfile} ../

popd

rm -rf ${pdir}
