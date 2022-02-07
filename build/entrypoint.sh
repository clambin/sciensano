#!/bin/sh

if [ -z "$MOUNTPOINT" ]; then
  MOUNTPOINT=/data
fi

if [ -z "$FILENAME" ]; then
  FILENAME=TF_SOC_POP_STRUCT_2021
fi


URL=https://statbel.fgov.be/sites/default/files/files/opendata/bevolking%20naar%20woonplaats%2C%20nationaliteit%20burgelijke%20staat%20%2C%20leeftijd%20en%20geslacht

if [ ! -f "$MOUNTPOINT"/population/"$FILENAME".txt ]; then
  echo "Downloading population file ..."
  mkdir -p "$MOUNTPOINT"/download
  mkdir -p "$MOUNTPOINT"/population
  cd "$MOUNTPOINT"/download || exit 1
  rm -f "$FILENAME".zip
  wget "$URL"/"$FILENAME".zip || exit 1
  unzip "$FILENAME".zip -d ../population -o || exit 1
  echo "Done. Starting sciensano server"
fi

cd - >/dev/null || exit 1
exec ./sciensano "$@"
