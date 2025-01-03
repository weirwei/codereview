# codereview

`codereview` 是一个用于代码审查的工具。它可以对代码审查，并生成 `markdown` 格式的审查报告。

报告示例：

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
### 20241125
**feature**
1. 支持目录树规范检查，并新增配置项 `tree_standard`，见 ([配置](#配置))
2. 支持代码风险评级
3. 新增 `-v` 选项，支持查看版本
4. 优化 `-d` 选项，打印的日志更方便查看
5. 优化prompt 和输出排版，输出更稳定

### 20241030
**feature**
1. 审核方式从原本的单文件审核，改成多文件合并审核，最大程度满足 `maxTokens`。（review 进度功能受影响）

**bug fix**
1.  修复 `gptserver 报错优化，目前报错会当做空接口，直接 LGTM`。 现在会提示 “llm 结果为空”

### 20241022
**feature**
1. 修改默认模型为 `claude-3-5-sonnet-20240620`

### 20241018
**feature**
1. 新增 `-m` 选项，支持自定义模型。
2. 新增 `-d` 选项，支持控制台打印详细日志。
3. 优化 prompt ，输出更加规范。

**bug fix**
1. 修复 diff 删除文件报错的问题。

## 安装

使用以下命令安装 `codereview`：

```bash
go install  
```

## 使用

### Quick Start

切换到工程根目录下，使用以下命令执行代码审查：

```bash
codereview  
```

### Options

- `-h` 或 `--help`: 显示帮助信息
- `-o`: 指定输出的文件，如：`codereview -o result.md`
- `-p`: 指定审查的 `package`，可批量，用 `,` 分割。如：`codereview -p controller/common,codereview/reviewexample`
- `-m`: 指定模型。
- `-d`: debug 模式。感觉结果不太对劲的时候可以加上看看日志。

## 配置

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

1. ~~Debug 模式，优化报错提示；~~
2. ~~规范分级；~~
3. ~~目录树检查；~~
4. 补充通用规范；
5. 优化输出效果，减少无意义的建议；
6. 支持大文件 `review`，大文件拆函数；
7. 支持参与到 `CI/CD` 流程中；
8. review 指定某个文件；
9. review 未提交文件；

## BUGS
1. ~~删除的文件运行 git diff 会报错~~
2. ~~gptserver 报错优化，目前报错会当做空接口，直接 LGTM~~ 现在会提示 “llm 结果为空”
   

## 模型测评

- o1-preview: 效果还行。速度慢。并且因为不支持流，请求ptserver可能造成504超时
- o1-mini: 无意义输出太多。速度慢
- gpt-4o-mini: 效果还行，可能会无法联系所有内容。速度很快
- gpt-4o: 效果不错。速度快
- glm-4v: 输出没法看。速度快
- glm-4-flash: 无意义输出太多。速度快
- Ernie-4.0-8K: 效果还行，但是输出格式不正确。速度中等
- Ernie-3.5-8K: 效果差。速度慢