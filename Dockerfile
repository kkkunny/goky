FROM archlinux/base

RUN echo "Server = https://mirrors.ustc.edu.cn/archlinux/\$repo/os/\$arch" > /etc/pacman.d/mirrorlist
RUN echo "Server = https://mirrors.tuna.tsinghua.edu.cn/archlinux/\$repo/os/\$arch" >> /etc/pacman.d/mirrorlist
RUN pacman-key --init
RUN pacman --noconfirm -Sy archlinux-keyring
RUN pacman --noconfirm -Syu make git gcc go llvm

WORKDIR /klang
COPY . .
RUN go mod download
RUN make build