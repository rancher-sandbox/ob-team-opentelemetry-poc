# let us download a file with curl on Linux command line #

VERSION="0.125.0" # otel collector version
ARCH="amd64" # go archicture
URL="https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${VERSION}/ocb_${VERSION}_linux_${ARCH}"

wget -L $URL
ls -l
mv "./ocb_${VERSION}_linux_$ARCH" ocb
chmod +x ocb
sudo chown -R root:root ./ocb
sudo mv -v ocb /usr/bin