FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y git
RUN apt install -y wget
RUN apt install -y cmake
RUN apt install -y ninja-build
RUN apt install -y pkg-config
RUN apt install -y python3
RUN apt install -y python3-pip

# UPMEM SDK
WORKDIR /root
RUN wget sdk-releases.upmem.com/2023.2.0/ubuntu_22.04/upmem-2023.2.0-Linux-x86_64.tar.gz
RUN tar -zxvf upmem-2023.2.0-Linux-x86_64.tar.gz
RUN echo "source /root/upmem-2023.2.0-Linux-x86_64/upmem_env.sh" > /root/.bashrc

WORKDIR /root/uPIMulator
