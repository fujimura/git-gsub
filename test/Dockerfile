FROM golang:1.14

RUN apt-get update
RUN apt-get install -y ruby
RUN gem install bundler

WORKDIR /app

COPY ./ ./

RUN make
RUN bundle

CMD ["bundle", "exec", "rspec"]
