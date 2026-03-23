# Pocket VM Lab

一个"小而完整"的玩具字节码虚拟机项目，用 Go 实现，用于测试和展示编码能力。

## 项目目标

构建一个闭环的简单虚拟机系统，包括：
- 字节码格式与虚拟机执行器
- 汇编器（文本 → 字节码）
- 调试器（单步执行、状态查看）
- 完整的测试覆盖

## 当前状态

已完成 **Milestone 0** 和 **Milestone 1**，实现了一个最小栈式虚拟机。

### 已实现指令

| 指令 | 功能 |
|------|------|
| `CONST` | 将常量压入栈 |
| `ADD`  | 弹出两值，压入和 |
| `PRINT`| 弹出并打印栈顶值 |
| `HALT` | 停止执行 |

### 快速开始

```bash
# 构建
go build -o pocket-vm-lab ./cmd/pocket-vm-lab

# 运行示例程序（计算 3 + 5 = 8）
./pocket-vm-lab -demo

# 运行字节码文件
./pocket-vm-lab -file testdata/sample.bin

# 运行测试
go test ./... -v
```

### 项目结构

```
pocket-vm-lab/
├── go.mod
├── README.md
├── cmd/
│   └── pocket-vm-lab/
│       ├── main.go           # CLI 入口
│       └── main_test.go      # CLI 测试
├── internal/
│   └── vm/
│       ├── opcode.go         # 操作码定义
│       ├── vm.go             # 虚拟机实现
│       └── vm_test.go        # 单元测试
└── testdata/
    └── sample.bin            # 示例字节码
```

## 技术栈

采用 **Go** 实现，适合测试：
- 清晰的工程实现
- CLI 工具组织
- 模块化设计
- 易扩展的工具链

## 开发路线

### 已完成 ✅
- [x] Go 模块初始化与 CLI 入口
- [x] 栈式虚拟机核心
- [x] 基础指令集（CONST, ADD, PRINT, HALT）
- [x] 单元测试覆盖

### 进行中 🚧
- [ ] 扩展指令集（SUB, MUL, DIV, DUP, MOV, CMP, JMP, JZ, JNZ）
- [ ] 汇编器实现（支持标签、注释、字面量）
- [ ] 控制流支持（循环、条件分支）

### 计划中 📋
- [ ] 调试器（单步执行、断点、执行追踪）
- [ ] 反汇编器
- [ ] 函数调用机制
- [ ] 标准库扩展

## 汇编示例

```asm
start:
  CONST 0
  CONST 10
loop:
  PRINT
  CONST 1
  ADD
  DUP
  CONST 10
  CMP
  JNZ loop
  HALT
```

## 贡献与扩展

如果你想参与这个项目或将其用作 coding agent 的 benchmark，核心原则是：
- **可运行性**：每次改动都保持系统可编译、可测试
- **可测试性**：为每条指令和功能补测试
- **可解释性**：保持文档与实现同步
- **可演进性**：控制复杂度，保持模块边界清晰

## License

MIT
