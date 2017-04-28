set +eu
if [ -z "$CGO_LDFLAGS" ]; then
      echo "CGO_LDFLAGS is not defined."
      echo "Ex: export CGO_LDFLAGS=\"-L/usr/local/lib/ -lmecab -lstdc++\""
      exit 1
fi
if [ -z "$CGO_CFLAGS" ]; then
      echo "CGO_CFLAGS is not defined."
      echo "Ex: export CGO_CFLAGS=\"-I/usr/local/include\""
      exit 1
fi
set -eu
