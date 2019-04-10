#!/bin/sh

rm -rf build
mkdir -p build/webui
cd webui
GOOS=linux go build
cd ..
mv webui/webui build/webui
cp -R webui/templates build/webui/templates
cd build
tar czfv webui.tgz webui
rm -rf build/webui
