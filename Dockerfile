# syntax=docker/dockerfile:1

# syntax=docker/dockerfile:1

FROM golang:1.18

RUN mkdir /forum
WORKDIR /forum
# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
# download all packages in mod file
RUN go mod download
# upload the entire 'forum' application
ADD . /forum
RUN go mod tidy
# log-in to Git before cloning the 'forum' repository
# RUN git config --global user.email "sedmakh2@gmail.com"
# RUN git config --global user.name "AthenaHTA2"
# RUN git config --global --add url."git@learn.01founders.co:".insteadOf "https://git.learn.01founders.co/"
# RUN go get git.learn.01founders.co/AthenaHTA2/forum/mainForum
# RUN cd /audit && git clone https://learn.01founders.co/git/AthenaHTA2/forum.git
# Dockerfile is in the 'mainForum' folder
RUN cd /forum
# Next build a static application binary named 'binaryForum'
RUN go build -o binaryForum
# The port that connects to docker daemon
EXPOSE 8080
LABEL version="1.0"
LABEL description="Project forum Created by Sonal, Nathan, Kingsley, Helena"
# Tell Docker to execute the 'binaryForum' command when this image is used to start a container.
ENTRYPOINT [ "/forum/binaryForum" ]

#FROM golang:1.18

#RUN mkdir /forum
#WORKDIR /forum
# Download necessary Go modules
#COPY go.mod ./
#COPY go.sum ./
# download all packages in mod file
#RUN go mod download
# upload the entire 'forum' application
#COPY . /forum
#RUN go mod tidy
# log-in to Git before cloning the 'forum' repository
# RUN git config --global user.email "sedmakh2@gmail.com"
# RUN git config --global user.name "AthenaHTA2"
# RUN git config --global --add url."git@learn.01founders.co:".insteadOf "https://git.learn.01founders.co/"
# RUN go get git.learn.01founders.co/AthenaHTA2/forum/mainForum
# RUN cd /audit && git clone https://learn.01founders.co/git/AthenaHTA2/forum.git
# Dockerfile is in the 'mainForum' folder
#WORKDIR /forum
#LABEL version="1.0"
#LABEL description="Project forum Created by Sonal, Kingsley, Nathan, Helena"
#CMD go run main.go
# Next build a static application binary named 'binaryForum'
#RUN go build -o binaryForum
# The port that connects to docker daemon
#EXPOSE 8080
# Tell Docker to execute the 'binaryForum' command when this image is used to start a container.
#ENTRYPOINT [ "/forum/binaryForum" ]



