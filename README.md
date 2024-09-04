![MemHunter](https://socialify.git.ci/b3nguang/MemHunter/image?description=1&font=Inter&forks=1&issues=1&language=1&logo=https%3A%2F%2Favatars.githubusercontent.com%2Fu%2F121670274%3Fs%3D400%26u%3D686132087f2e2324958b610f905a1b388478295b%26v%3D4&name=1&owner=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Dark)

## ✈️ 一、工具特性

**MemHunter** 是一个用于在 Windows 系统中搜索进程内存中指定字符串的工具。该工具旨在帮助应急响应人员在日志缺失或加密的情况下，定位潜在的恶意行为。主要特性如下：

- **多进程并发搜索**：支持并发处理，能够快速遍历系统中的多个进程，加快检索速度。
- **基于 Windows API**：利用 Windows 提供的底层 API（如 `VirtualQueryEx`、`ReadProcessMemory` 等）直接与进程内存交互，确保工具的可靠性和高效性。
- **精准定位**：能够在进程的内存空间中精确找到用户指定的字符串，并输出相关的进程信息，包括 PID、进程名称和进程路径。
- **内存块遍历**：逐块遍历进程内存空间，确保所有潜在的内存区域都得到检查。

## 🚨 二、配置

**MemHunter** 基于 Go 语言开发，运行在 Windows 操作系统上。以下是运行该工具的前提条件和配置步骤：

- **操作系统**：Windows 7 及以上版本。
- **Go 语言环境**：确保已安装 Go 1.16 及以上版本。
- **必要权限**：为了能够读取其他进程的内存，需以管理员身份运行该工具。

**安装步骤**：

1. 安装 Go 语言环境，并将其添加到系统路径中。

2. 克隆或下载 **MemHunter** 源代码到本地：
   ```bash
   git clone https://github.com/b3nguang/MemHunter.git
   ```
   
3. 在终端中导航到项目目录并执行以下命令以编译可执行文件：
   ```bash
   go build -o MemHunter main.go
   ```

## 🚀 三、使用

1. 打开终端并导航到 **MemHunter** 所在目录。
2. 运行工具并输入需要搜索的目标字符串：
   ```bash
   MemHunter
   ```
   程序启动后，会提示输入需要检索的字符串：
   ```
   [*]请输入检索的字符串:
   ```
   输入目标字符串后，工具将自动开始在所有运行的进程内存中搜索该字符串。
3. 搜索结果将以以下格式输出：
   ```
   [+]在 PID 1234 的地址 0x7ffdf0000 处找到字符串: example
   [+]进程名称: example.exe
   [+]进程文件路径: C:\Program Files\Example\example.exe
   ```
4. 程序将持续运行直到所有进程都被遍历完成。搜索结果将按进程逐一输出。

**注意事项**：
- 该工具需要以管理员权限运行，以确保能够访问所有进程的内存空间。
- 如果目标进程的内存被保护或不可读，工具可能无法成功检索到相关信息。

**MemHunter** 是一个强大且易用的内存检索工具，适用于在复杂应急响应场景下快速定位潜在威胁。

## 🖐 四、免责声明

1. 如果您下载、安装、使用、修改本工具及相关代码，即表明您信任本工具
2. 在使用本工具时造成对您自己或他人任何形式的损失和伤害，我们不承担任何责任
3. 如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任
4. 请您务必审慎阅读、充分理解各条款内容，特别是免除或者限制责任的条款，并选择接受或不接受
5. 除非您已阅读并接受本协议所有条款，否则您无权下载、安装或使用本工具
6. 您的下载、安装、使用等行为即视为您已阅读并同意上述协议的约束

## 🙏 五、参考项目

https://github.com/Fheidt12/Windows_Memory_Search

https://www.freebuf.com/sectool/408673.html

在这里向 `Windows_Memory_Search` 项目献上最诚挚的敬意。
