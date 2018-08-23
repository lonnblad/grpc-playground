### GRPC Client

## Setup certificates

$ certstrap request-cert --common-name grpc.iot.enlight.skf.com
$ certstrap sign --CA "ca" grpc.iot.enlight.skf.com
