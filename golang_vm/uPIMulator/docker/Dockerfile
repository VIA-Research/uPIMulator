FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y git
RUN apt install -y wget
RUN apt install -y git-lfs
RUN apt install -y cmake
RUN apt install -y ninja-build
RUN apt install -y pkg-config
RUN apt install -y python3
RUN apt install -y python3-pip

# UPMEM SDK
WORKDIR /root
RUN git clone https://github.com/kagandikmen/upmem-sdk.git /opt/upmem-sdk-archive
RUN cd /opt/upmem-sdk-archive && git lfs install && git lfs pull
RUN tar -zxvf /opt/upmem-sdk-archive/2023.1.0/upmem-2023.1.0-Linux-x86_64.tar.gz -C /root
RUN mv /root/upmem-2023.1.0-Linux-x86_64 /root/upmem
RUN rm -rf /opt/upmem-sdk-archive
RUN echo "source /root/upmem/upmem_env.sh" > /root/.bashrc
ENV UPMEM_HOME=/root/upmem
ENV PATH="${UPMEM_HOME}/bin:${PATH}"

WORKDIR /root/uPIMulator
