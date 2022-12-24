FROM golang:1.19.3-alpine

RUN mkdir -p /usr/share/man/man1 && \
  apk add --update --upgrade \
  git \
  ca-certificates \
  zsh \
  curl \
  wget \
  procps


WORKDIR /home/app

RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest 
RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest

RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.2/zsh-in-docker.sh)" -- \
  -t https://github.com/romkatv/powerlevel10k \
  -p git \
  -p git-flow \
  -p https://github.com/zdharma-continuum/fast-syntax-highlighting \
  -p https://github.com/zsh-users/zsh-autosuggestions \
  -p https://github.com/zsh-users/zsh-completions \
  -a 'export TERM=xterm-256color'

RUN echo '[[ ! -f /home/.p10k.zsh ]] || source /home/.p10k.zsh' >> ~/.zshrc && \
  echo 'HISTFILE=/home/zsh/.zsh_history' >> ~/.zshrc

CMD [ "tail", "-f" , "/dev/null" ]
