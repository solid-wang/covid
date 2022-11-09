# devops规范

* `Required:`Tag时在message中行头添加 `+devops` 注释块，判断是否启动流水线。
* `Optional:`当仓库中存在多个项目时，使用 `app` 注释块用于表示启动哪些项目流水线（如: `app=name1,name2` ），若没有添加 `server` 注释块时，表示启动该仓库下所有项目的流水线。
* `Optional:`如需发布，请使用 `deploy` 注释块表示部署到哪个环境（如: `deploy=dev`），若没有添加 `deploy` 注释块，表示只编译，不发布。
* ~~`Optional:`发布服务时，启动命令变更是默认关闭的，若要开启，使用 `command` 注释块 （如: `command=ture`）。 *仅当 `deploy` 存在时生效*。~~
* ~~`Optional:`发布服务时，配置文件变更是默认关闭的，若要开启，使用 `config` 注释块 （如: `config=ture`）。  *仅当 `deploy` 存在时生效*。~~

* * *
* *version = tag + short commit*
* *每个注释块之间必须以空格分割*

# 设计原因

* 1.为什么不选择push或merge来触发流水线?
  * push动作是有commit中带有信息，当前commit需要在不同时间发布不同环境时，需要修改commit message，这又会新增一次没必要的commit
  * merge动作必须存在两个分支，一些项目人数较少，git使用较为简单，或直接push主干分支
* 2.为什么要用message方式: 
  * 开发活动中，能够自主决定是否走流水线，流水线中的一些变量传递，没有找到其他更合适的方式。
* 3.为什么要添加`server`: 
  * 存在一个仓库多个项目，拆分工程量较大
* 为什么不以分支来对应发布环境？
  * 很多项目没有按照git最佳实践来运作，没有为每个环境建立分支，若devops在初期要求git规范，会降低接入效率。目前策略是先用起来和好用为导向。


# 优化 Docker 镜像
* 构建优化的 Docker 镜像，因为大型 Docker 镜像会占用大量空间，下载时间较长，连接速度较慢。如果可能，请避免对所有作业使用一个打镜像。使用多个较小的镜像，每个镜像用于特定任务，下载和运行速度更快。
* 尝试使用预装软件的自定义 Docker 镜像。下载更大的预配置镜像通常比每次使用通用镜像并在其上安装软件要快得多。Docker 的编写 Dockerfiles 的最佳实践 有更多关于构建高效 Docker 镜像的信息。
* 减小 Docker 镜像大小的方法：

  * 使用小型基础镜像，例如 debian-slim。
  * 如果不是严格需要，不要安装像 vim、curl 等便利工具。
  * 打造专属的开发镜像。
  * 禁用由软件包安装的手册页和文档以节省空间。
  * 减少 RUN 层并结合软件安装步骤。
  * 使用 multi-stage builds 将多个使用构建器模式的 Dockerfile 合并为一个 Dockerfile，可以减少镜像大小。
  * 如果使用 apt，添加 --no-install-recommends 以避免不必要的包。
  * 清理最后不再需要的缓存和文件。例如 rm -rf /var/lib/apt/lists/* 适用于 Debian 和 Ubuntu，或 yum clean all 适用于RHEL和CentOS。
  * 使用 dive 或 DockerSlim 等工具来分析和缩小镜像。
  * 为了简化 Docker 镜像管理，您可以创建一个专门的组来管理 Docker 镜像 并使用 CI/CD 流水线测试、构建和发布它们。