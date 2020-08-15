FROM golang:1.14.7
RUN sh -c "$(wget -O- https://raw.githubusercontent.com/deluan/zsh-in-docker/master/zsh-in-docker.sh)" -- \
    -t random \
    -p git \
    -p ssh-agent \
    -p https://github.com/zsh-users/zsh-autosuggestions \
    -p https://github.com/zsh-users/zsh-completions \
    -p https://github.com/zsh-users/zsh-history-substring-search
RUN apt-get -y update
RUN apt-get -y remove docker docker-engine docker.io
RUN apt install -y docker.io
RUN apt install -y docker-compose