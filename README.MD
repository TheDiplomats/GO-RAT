TODO List

-Encrypt Communications
Especially between implant and server
should use TLS by modyfying grpc.DialOptions

-Connection Retries
Implant needs to retry in event it can't connect to server...if implant
dies access is lost in implant.go grpc.Dial() again to success

-Register Implant
Need to identify different implants for different machines. Admin client
only expects one, perhaps give each a unique id #

-Add Database Persistence
Pretty obvious

-Implement multiple implants

-Improve Implant Functionality

-Chain OS Commands
