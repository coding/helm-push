# helm-push


**已弃用**！已迁移至[新仓库](https://coding-public.coding.net/public/helm-push/helm-push/git)。


用于推送 Chart 包到 Coding 制品库的 Helm 插件。

## 安装

```shell
$ helm plugin install https://github.com/Coding/helm-push
```

## 使用

添加 Chart 仓库
```shell
helm repo add --username <username> --password <password> <repository> https://<team>-helm.pkg.coding.net/<project>/<repository>
```

推送 Chart 包
```shell
helm push mychart.tgz <repository>
```

推送 Chart 目录
```shell
helm push . <repository>
```

