FROM golang:latest
# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Update package list and install required packages
RUN apt-get update && apt-get install -y \
    openssh-server \
    git \
    && rm -rf /var/lib/apt/lists/*

# Create git user with normal bash shell
RUN useradd -m -d /home/git -s /bin/bash git

# Create .ssh directory for git user (not needed for anonymous access but keeping for completeness)
RUN mkdir -p /home/git/.ssh && \
    chown git:git /home/git/.ssh && \
    chmod 700 /home/git/.ssh

# Create authorized_keys file (not needed for anonymous access)
RUN touch /home/git/.ssh/authorized_keys && \
    chown git:git /home/git/.ssh/authorized_keys && \
    chmod 600 /home/git/.ssh/authorized_keys

# Create git-shell-commands directory for reference (not needed with bash shell)
RUN mkdir -p /home/git/git-shell-commands && \
    chown git:git /home/git/git-shell-commands

COPY go.mod go.sum ./
RUN go mod download
WORKDIR /
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ccg ./cmd

COPY .env /.env
COPY .data /.data

COPY scripts/git-receive-pack /usr/local/bin/git-receive-pack
COPY scripts/git-upload-pack /usr/local/bin/git-upload-pack
# Make wrapper scripts executable
RUN chmod +x /usr/local/bin/git-receive-pack && \
    chmod +x /usr/local/bin/git-upload-pack

# Configure SSH daemon for anonymous access with bash shell
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication no/' /etc/ssh/sshd_config && \
    sed -i 's/#PermitEmptyPasswords no/PermitEmptyPasswords yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config

# Set empty password for git user to allow anonymous access
RUN passwd -d git

# Allow SSH access for git user (removed shell restrictions)
RUN echo "Match User git" >> /etc/ssh/sshd_config && \
    echo "    AllowTcpForwarding no" >> /etc/ssh/sshd_config && \
    echo "    X11Forwarding no" >> /etc/ssh/sshd_config && \
    echo "    PasswordAuthentication yes" >> /etc/ssh/sshd_config && \
    echo "    PermitEmptyPasswords yes" >> /etc/ssh/sshd_config

# Create git-shell-commands with custom git commands
RUN echo '#!/bin/bash' > /home/git/git-shell-commands/git-receive-pack && \
    echo 'exec /usr/local/bin/git-receive-pack "$@"' >> /home/git/git-shell-commands/git-receive-pack && \
    chmod +x /home/git/git-shell-commands/git-receive-pack && \
    chown git:git /home/git/git-shell-commands/git-receive-pack

RUN echo '#!/bin/bash' > /home/git/git-shell-commands/git-upload-pack && \
    echo 'exec /usr/local/bin/git-upload-pack "$@"' >> /home/git/git-shell-commands/git-upload-pack && \
    chmod +x /home/git/git-shell-commands/git-upload-pack && \
    chown git:git /home/git/git-shell-commands/git-upload-pack

# Create repositories directory
RUN mkdir -p /home/git/repositories && \
    chown git:git /home/git/repositories && \
    chown -R git:git /.data

# Generate SSH host keys
RUN ssh-keygen -A

# Expose SSH port
EXPOSE 22

# Start SSH daemon
CMD ["/usr/sbin/sshd", "-D"]