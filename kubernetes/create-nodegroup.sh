#!/bin/bash


SSH_PUBLIC_KEY='ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDBci+hXvzWxsCjbvdE580Ys+iB7yGs8Ya1qr4eQOT6vyt1Xr86/qdMcGuqXjWxs1bIbqktC2DrNDIv/XnmS35s2quX/abpHtHjvuYru7yPRCIB14NHkLWgvVkQWnqYDVOS0I1gQk4qQwaTIN8gJ6KVsgPnWoveAYCObQkzhF/8twd9WxnA+/iaV5hthp5xFWXuzIpVoUiebwCGJJ8FO3R1o9xc4MTXpQ4iEospePByNBCAuJ9fdW5VwlOXAVVN7vcY44ZkFN3HKCnlwKYLv+uiCRgMM1vjHJ3UJyVB/1mYSdxuZd1sCOe24fczWQK9geRCGPAmYiUAbXEmsoFE3BpR8jcIUcFsIiOphvsHMtLDsGU9m8lQhaHqzwGNb02g5NMg5pnNFsl3Fa4bK6Ut4RRIQVRZIOOYtp5DngE2J80339w2hE4ZOCuTmrFJsZ3WnAxu2HdB78kYHnyXIWl23QdVmQjJVp/Zz7PJcVP6HoFXC68g9yD5wpYqXtOPKawXBlV5ngA1RsQs5+H/jA+eMz9TxZRsZBqdcUSB5IQ2iVWciwTco3zhBnI+rfQG0ufEeQfmHPXS6h4bgtJAHdPgWlRZcpqlSL//Ixjsju1lcLgDe3b2QsyZTG4AVDuto0iW1zd3Fsi61BRIygjfv4CgPvoIoaxNAFbnUw9WSkfb2VCLfw== Ineat #2'

main() {
    eksctl create nodegroup \
        --cluster flapflap \
        --region us-east-1 \
        --name ec2-micro \
        --node-type t2.micro \
        --nodes 1 \
        --nodes-min 1 \
        --nodes-max 10 \
        --ssh-access \
        --ssh-public-key "$SSH_PUBLIC_KEY"
}

[[ ${BASH_SOURCE[0]} == $0 ]] && main "$@"
