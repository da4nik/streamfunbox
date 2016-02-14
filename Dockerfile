FROM busybox:glibc


MAINTAINER Maksim Stepanov <stepanov.ms@mail.ru>

CMD mkdir -p /StreamFunBox
WORKDIR /StreamFunBox

ADD ./streamfunbox ./
ADD images/ images/
ADD sounds/ sounds/
ADD public/ public/

EXPOSE 8998

CMD "./streamfunbox"
