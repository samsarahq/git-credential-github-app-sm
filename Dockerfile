FROM golang:1.25.0 AS build

RUN mkdir /code

COPY . /code

WORKDIR /code

ENV GOBIN /code/build
RUN mkdir build && go install


FROM scratch
COPY --from=build --chmod=555 /code/build/git-credential-github-app-sm /