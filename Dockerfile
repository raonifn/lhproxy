FROM scratch
WORKDIR /lhproxy
ENV PATH /lhproxy
ADD ./build/out/linux-amd64/lhproxy/lhproxy .
CMD [ "lhproxy", "help" ]
