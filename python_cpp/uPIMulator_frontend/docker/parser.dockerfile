FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONPATH="/root/upmem_linker/src:$PYTHONPATH"

RUN chmod 1777 /tmp

RUN apt update
RUN apt install -y git
RUN apt install -y wget
RUN apt install -y vim
RUN apt install -y tmux
RUN apt install -y default-jre
RUN apt install -y default-jdk

# Python 3.10
RUN apt update
RUN apt upgrade -y
RUN apt install -y software-properties-common
RUN add-apt-repository ppa:deadsnakes/ppa
RUN apt install -y python3.10
RUN apt install -y python3-pip

# ANTLR4
WORKDIR /root
RUN wget https://www.antlr.org/download/antlr-4.9.2-complete.jar
ENV CLASSPATH="/root/antlr-4.9.2-complete.jar:$CLASSPATH"
RUN echo "alias antlr4='java -Xmx500M -cp \"/root/antlr-4.9.2-complete.jar:$CLASSPATH\" org.antlr.v4.Tool'" >> /root/.bashrc
RUN echo "alias grun='java -Xmx500M -cp \"/root/antlr-4.9.2-complete.jar:$CLASSPATH\" org.antlr.v4.gui.TestRig'" >> /root/.bashrc

WORKDIR /root/upmem_compiler