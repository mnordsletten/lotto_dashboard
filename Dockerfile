from golang:1.11

WORKDIR /go/src/lotto_dashboard

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . .
RUN dep ensure
RUN go install
ENTRYPOINT ["lotto_dashboard"]
