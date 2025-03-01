# codereview

`⁠codereview` is a tool for code review. It can review code and generate a review report in `⁠markdown` format.

report example：

````markdown
---

> codereview/reviewexample/bad_review_example.go:10-15

```go
func Atoi(s string) int {
        var err error
        _ = err
        atoi, err := strconv.Atoi(s)
        return atoi
}
```

- 隐患: 未处理错误。
  评级: P1
  建议: 添加错误处理逻辑，例如返回错误或使用默认值。

- 隐患: 声明了未使用的变量 err
  评级: P2
  建议: 删除未使用的变量声明

- 隐患: 函数命名与标准库冲突
  评级: P3
  建议: 重命名函数，避免与 strconv.Atoi 冲突

````

## RELEASE NOTES

## 安装

使用以下命令安装 `codereview`：

```bash
go install github.com/weirwei/codereview@latest
```

## 使用

### Quick Start

切换到工程根目录下，使用以下命令执行代码审查：

```bash
codereview  
```

### 可用命令

`codereview` 提供了以下主要命令：

1. 主命令：
   - `codereview`: 执行代码审查。
     选项：
     - `--pkg`, `-p`: 指定要审查的包，可以用逗号分隔多个包。
     - `--version`, `-v`: 显示版本信息。
     - `--debug`, `-d`: 设置日志级别为 DEBUG。


2. 配置命令：
   - `codereview config`: 管理配置
     - `set <key> <value>`: 设置配置项
     - `get <key>`: 获取配置项的值
     - `list`: 列出所有配置项

3. 版本命令：
   - `codereview version`: 显示 CodeReview 的版本号

### 命令示例

1. 执行代码审查：
   ```bash
   codereview
   ```

2. 审查特定包：
   ```bash
   codereview -p package1,package2
   ```

3. 以调试模式运行：
   ```bash
   codereview -d
   ```
   或
   ```bash
   codereview --debug
   ```

4. 显示版本信息：
   ```bash
   codereview -v
   ```
   或
   ```bash
   codereview version
   ```

4. 配置管理：
   ```bash
   codereview config set model gpt-3.5-turbo
   codereview config get model
   codereview config list
   ```


## 配置文件

在项目根目录下创建 `.codereview.yml` 文件。无配置情况下，会根据默认配置运行程序。

以下是一个 `.codereview.yml` 的示例：

```yaml
languages: # 启用 code review 的语言，目前只支持 go  
  - go  

code:  
  git:  
    review_branch: HEAD # 代码review的分支，默认为当前分支  
    compare_branch: origin/master # 对比的分支，默认为 origin/master  
  files:  
    ignore: # 忽略的文件  
      - .*_test.go  
      - .*_mock.go  

# 知识库。自定义  
knowledge:  
  enable: true
  tree_standard:
    /api: 
      - 封装请求依赖方的api
      - 不应该有复杂的业务逻辑
  custom:  
    go:  
      - regexp: goroutine\.NewMulti # 匹配对应的代码片段  
        rules:  
          - goroutine.NewMulti() 返回的对象必须调用 Wait() 函数  
      - regexp: \.\(.+\)  
        rules:  
          - 类型断言的时候需要使用 ", ok" 来检查结果。  
```

### 配置参数

#### languages

`languages` 配置的是进行代码审查的语言，`array` 类型。会根据语言的匹配包含对应后缀的文件。  
默认情况下，会审查所有文件。

支持语言：`go`

#### code

`code` 是进行审查的代码内容的配置。

| 参数  | 说明  | 参数选项 |
| --- | --- | --- |
| git | git相关的配置 | `review_branch`: 代码审查的分支，`string` 类型，默认当前分支<br/>`compare_branch`: 对比的分支，`string` 类型，默认 `origin/master` |
| files | 审查的文件配置 | `ignore`: 忽略的文件，`array-string` 类型，支持正则表达式语法 |

#### knowledge

`knowledge` 是知识库配置。根据正则表达式判断该内容的审查是否需要应用指定的规则。

| 参数  | 说明  | 参数选项 |
| --- | --- | --- |
| enable | 是否启动 `knowledge`，`true` or `false` | -   |
| custom | 自定义的 `knowledge`，kv 形式 | k 为对应的 `language`，v 为配置内容，`array` 类型。v 包含：<br/>`regexp`: 匹配生效的代码片段，`string` 类型，支持正则表达式；<br/>`rules`: 指导代码审查的规则，`array-string` 类型。 |
| tree_standard | 目录树检查的 `knowledge`，kv 形式 | k 为对应的 `文件路径`，v 为配置内容，`array-string` 类型。 |

## ToDo
1. 支持加入 ci/cd 流程;
