FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

USER vscode

# ZSH
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.5/zsh-in-docker.sh)" -- \
  -t robbyrussell \
  -p git

# DB Tools
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/gzuidhof/tygo@latest

USER root

# Install Node and PNPM
RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && nvm install 20 && nvm use 20" 
RUN su vscode -c "npm install -g pnpm"

# Persist bash history
ARG USERNAME=vscode
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory \
  && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"