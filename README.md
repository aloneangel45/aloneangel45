<h1 align="center">Hi there, I'm Arda (<a href="https://github.com/aloneangel45">@aloneangel45</a>) 👋</h1>

<p align="center">
  <strong>Systems Engineer & Agentic AI Infrastructure Architect</strong><br>
  Engineering high-throughput, deterministic runtime tooling and low-level system utilities for the next generation of autonomous AI developers.
</p>

<p align="center">
  <a href="https://github.com/aloneangel45/toolprune"><img src="https://img.shields.io/badge/ToolPrune-v1.0.0-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="ToolPrune"/></a>
  <a href="https://github.com/aloneangel45/traceport"><img src="https://img.shields.io/badge/TracePort-v1.0.0-8b5cf6?style=for-the-badge&logo=satellite&logoColor=white" alt="TracePort"/></a>
  <a href="https://github.com/aloneangel45/FocusFlow"><img src="https://img.shields.io/badge/FocusFlow_V2-v1.0.0-cba6f7?style=for-the-badge&logo=windows&logoColor=white" alt="FocusFlow"/></a>
  <a href="https://github.com/aloneangel45/AgentFixture"><img src="https://img.shields.io/badge/AgentFixture-v1.0.0-f59e0b?style=for-the-badge&logo=githubactions&logoColor=white" alt="AgentFixture"/></a>
</p>

---

### 🚀 Featured Ecosystem: Agentic Infrastructure & Systems Tooling

A tightly integrated, zero-overhead developer toolchain designed to optimize context window economics, guarantee replay determinism, and maintain laser focus:

| Project | Role & Purpose | Key Highlights | Stack |
| :--- | :--- | :--- | :---: |
| ⚡ **[ToolPrune](https://github.com/aloneangel45/toolprune)** | **High-Throughput MCP Proxy & Schema Pruner** | Cuts verbose MCP JSON schemas by **91.4%**, drastically reducing token consumption and latency for Claude Code, Cursor, and Windsurf. | `Go` `MCP` `JSON-RPC` |
| 🛰️ **[TracePort](https://github.com/aloneangel45/traceport)** | **Deterministic Flight Recorder & Replay Debugger** | Black-box flight recorder capturing AI agent tool calls with Git shadow-tree checkpointing for **$0 zero-token error recovery**. | `Go` `Git` `SQLite WAL` |
| 📼 **[AgentFixture](https://github.com/aloneangel45/AgentFixture)** | **Zero-Cost Deterministic CI/CD Cassette Runner** | Ingests TracePort flight recordings into replayable offline cassettes. Run full AI agent regression test suites in CI for **$0 API cost**. | `Go` `CI/CD` `Proxy` |
| 🎯 **[FocusFlow V2](https://github.com/aloneangel45/FocusFlow)** | **Deep Work Discipline Engine & Desktop HUD** | High-DPI Catppuccin Mocha HUD with kernel-level SHA-256 atomic hosts rollback and low-overhead Win32/POSIX process supervision. | `Python` `Go` `Win32` |

---

### 🏗️ Ecosystem Architecture Workflow

```
   ┌─────────────────────────────────────────────────────────────┐
   │         AI Coding Agent (Claude Code / Cursor / Aider)      │
   └──────────────────────────────┬──────────────────────────────┘
                                  │
         ┌────────────────────────┴────────────────────────┐
         ▼                                                 ▼
┌──────────────────┐                              ┌──────────────────┐
│    ToolPrune     │                              │    TracePort     │
│ (Dynamic Schema  │                              │ (Flight Recorder │
│   91.4% Pruning) │                              │ & Shadow Trees)  │
└────────┬─────────┘                              └────────┬─────────┘
         │                                                 │
         ▼                                                 ▼
[Upstream MCP Servers]                            ┌──────────────────┐
(GitHub / Cloudflare / DB)                        │   AgentFixture   │
                                                  │ ($0 CI Cassette) │
                                                  └──────────────────┘
```

---

### 💻 Core Tech Stack & Engineering Focus

* **Languages:** Go, Rust, Python, C/C++, TypeScript, Modern Shell / PowerShell
* **Systems & Low-Level:** Win32 API (`Toolhelp32`, Job Objects), POSIX (`/proc`, `ps`), Signal Handling, IPC (Stdio, SSE, Unix Sockets)
* **Agentic & Protocols:** Model Context Protocol (MCP), JSON-RPC 2.0, EventSource / SSE, OpenTelemetry / OpenInference, Tree-Sitter AST
* **Storage & Determinism:** SQLite WAL Mode, SHA-256 Cryptographic Hashing, Git Tree / Blob internals, In-Memory Ring Buffers
* **DevOps & QA:** Multi-Platform GitHub Actions Matrix (Windows, Linux, macOS), Cross-Compilation, Race-Condition Auditing

---

### 📊 GitHub Activity & Metrics

<p align="center">
  <img src="https://github-readme-stats-eight-theta.vercel.app/api?username=aloneangel45&show_icons=true&theme=catppuccin_mocha&hide_border=true&title_color=cba6f7&text_color=cdd6f4&icon_color=89b4fa&bg_color=1e1e2e" alt="GitHub Stats" width="48%"/>
  <img src="https://streak-stats.demolab.com/?user=aloneangel45&theme=catppuccin_mocha&hide_border=true&stroke=cba6f7&ring=89b4fa&fire=f38ba8&currStreakLabel=cba6f7" alt="GitHub Streak" width="48%"/>
</p>

---

<p align="center">
  <em>“Build with zero bloat, verify with zero mocks, engineer for deterministic resilience.”</em><br>
  Crafted by <a href="https://github.com/aloneangel45">aloneangel45</a>
</p>
