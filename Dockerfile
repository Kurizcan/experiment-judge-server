FROM lcanboom/experiment_judge_base:latest
LABEL maintainer="kurizcan"
# 设置 Go 环境变量
ENV PATH=$PATH:/usr/local/go/bin GOPATH=/home/godev/gopath GO111MODULE=on GOPROXY=https://goproxy.io,direct
# 设置工作目录
WORKDIR /home/godev/gopath/src/experiment-judge-server/
# 复制源文件
COPY . /home/godev/gopath/src/experiment-judge-server/
RUN go build .
ENV PORT 9000
EXPOSE 9000
CMD ["sh", "start.sh"]