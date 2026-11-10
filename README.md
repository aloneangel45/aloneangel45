<h1 align="center">Hi there, I'm Arda (<a href="https://github.com/aloneangel45">@aloneangel45</a>)</h1>

<p align="center">
  <strong>Building reliable, deterministic infrastructure for AI agents and developer tooling.</strong><br>
  Engineering high-throughput, deterministic runtime proxies, checkpointing engines, and low-level system utilities for autonomous AI development workflows.
</p>

<p align="center">
  <a href="https://github.com/aloneangel45/toolprune/releases/tag/v1.2.0"><img src="https://img.shields.io/badge/ToolPrune-v1.2.0-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="ToolPrune v1.2.0"/></a>
  <a href="https://github.com/aloneangel45/traceport/releases/tag/v1.1.0"><img src="https://img.shields.io/badge/TracePort-v1.1.0-8b5cf6?style=for-the-badge&logo=git&logoColor=white" alt="TracePort v1.1.0"/></a>
  <a href="https://github.com/aloneangel45/AgentFixture/releases/tag/v1.1.0"><img src="https://img.shields.io/badge/AgentFixture-v1.1.0-f59e0b?style=for-the-badge&logo=githubactions&logoColor=white" alt="AgentFixture v1.1.0"/></a>
  <a href="https://github.com/aloneangel45/FocusFlow/releases/tag/v2.1.0"><img src="https://img.shields.io/badge/FocusFlow-v2.1.0-cba6f7?style=for-the-badge&logo=windows&logoColor=white" alt="FocusFlow v2.1.0"/></a>
</p>

---

### Featured Ecosystem: Agentic Infrastructure & Systems Tooling

A tightly integrated, zero-bloat developer toolchain engineered to compress context-window overhead, guarantee deterministic replay, and protect deep-work focus:

| Project | Role & Purpose | Key Highlights | CI | Stack |
| :--- | :--- | :--- | :---: | :---: |
| **[ToolPrune](https://github.com/aloneangel45/toolprune)** | **Adaptive MCP Proxy & Token-Budget Stream Filter** | Progressive `--max-schema-tokens` compaction and dynamic tool hydration cutting MCP schema tokens by **93.3%–96.5%** at **<0.4ms** cold-pass latency across Stdio, SSE, and Streamable HTTP (`/mcp`). | [![CI](https://github.com/aloneangel45/toolprune/actions/workflows/ci.yml/badge.svg)](https://github.com/aloneangel45/toolprune/actions/workflows/ci.yml) | `Go` `MCP` `JSON-RPC` |
| **[TracePort](https://github.com/aloneangel45/traceport)** | **Deterministic Flight Recorder & Zstd Shadow-Tree Engine** | Content-addressed workspace checkpointing with Zstd preset/trained dictionaries (**86.1% footprint reduction**, **557us** stat-cached scan) and O(1) JSONL flight logs for **$0 zero-token rollback**. | [![CI](https://github.com/aloneangel45/traceport/actions/workflows/ci.yml/badge.svg)](https://github.com/aloneangel45/traceport/actions/workflows/ci.yml) | `Go` `Zstd` `Git` |
| **[AgentFixture](https://github.com/aloneangel45/AgentFixture)** | **Zero-Cost Deterministic CI/CD Cassette Runner** | Native Anthropic & OpenAI SSE stream recording/replay (`instant`/`scaled`/`exact` timing), zero-outbound-network guard, 16MB MCP stdio proxy, Go `agenttest` SDK, and TracePort `v1.1` importer. | [![CI](https://github.com/aloneangel45/AgentFixture/actions/workflows/ci.yml/badge.svg)](https://github.com/aloneangel45/AgentFixture/actions/workflows/ci.yml) | `Go` `SSE` `CI/CD` |
| **[FocusFlow](https://github.com/aloneangel45/FocusFlow)** | **Deep Work Discipline Engine & Desktop HUD** | Unified Go CLI and zero-dependency Python Catppuccin Mocha HUD with crash-safe SHA-256 dual-stack IPv4/IPv6 hosts sinkholing, Win32 DACL preservation, and low-overhead Win32/POSIX process supervision. | [![CI](https://github.com/aloneangel45/FocusFlow/actions/workflows/ci.yml/badge.svg)](https://github.com/aloneangel45/FocusFlow/actions/workflows/ci.yml) | `Go` `Python` `Win32` |

---

### Ecosystem Architecture Workflow

```
   ┌───────────────────────────────────────────────────────────────┐
   │        AI Coding Agent (Claude Code / Cursor / Windsurf)      │
   └───────────────────────────────┬───────────────────────────────┘
                                   │
          ┌────────────────────────┴────────────────────────┐
          ▼                                                 ▼
┌─────────────────────┐                           ┌─────────────────────┐
│      ToolPrune      │                           │      TracePort      │
│  (Adaptive Stream   │                           │  (Zstd Shadow-Tree  │
│  93.3%-96.5% Prune) │                           │  & Flight Recorder) │
└──────────┬──────────┘                           └──────────┬──────────┘
           │                                                 │
           ▼                                                 ▼
 [Upstream MCP Fleet]                             ┌─────────────────────┐
 (Stdio / SSE / /mcp)                             │    AgentFixture     │
                                                  │ (SSE + MCP Offline  │
                                                  │   $0 CI Cassettes)  │
                                                  └─────────────────────┘
```

---

### Core Tech Stack & Engineering Focus

* **Languages:** Go, Rust, Python, C/C++, TypeScript, Modern Shell / PowerShell
* **Systems & Low-Level:** Win32 API (`Toolhelp32`, `ReplaceFileW` DACLs, Job Objects), POSIX (`/proc`, process groups, signal handling), IPC (Stdio, SSE, Streamable HTTP)
* **Agentic & Protocols:** Model Context Protocol (MCP), JSON-RPC 2.0, Anthropic & OpenAI Server-Sent Events (SSE), `cl100k_base` Token Calibration, Pre-Bind JSON Schema Validation
* **Storage & Determinism:** Dictionary-Trained Zstandard (Zstd) Content-Addressable Storage, Append-Only JSONL Event Logs, Crash-Safe Write-Ahead Journals, SHA-256 Integrity Verification
* **DevOps & QA:** Multi-OS GitHub Actions Matrices (Linux, Windows, macOS), Go Race Detector Auditing, Reproducible `CGO_ENABLED=0` Cross-Compilation (`amd64` / `arm64`)

---

### GitHub Activity & Metrics

<p align="center">
  <img src="https://github-readme-stats-eight-theta.vercel.app/api?username=aloneangel45&show_icons=true&theme=catppuccin_mocha&hide_border=true&title_color=cba6f7&text_color=cdd6f4&icon_color=89b4fa&bg_color=1e1e2e" alt="GitHub Stats" width="48%"/>
  <img src="https://streak-stats.demolab.com/?user=aloneangel45&theme=catppuccin_mocha&hide_border=true&stroke=cba6f7&ring=89b4fa&fire=f38ba8&currStreakLabel=cba6f7" alt="GitHub Streak" width="48%"/>
</p>

---

<p align="center">
  <em>"Build with zero bloat, verify with zero mocks, engineer for deterministic resilience."</em><br>
  Crafted by <a href="https://github.com/aloneangel45">aloneangel45</a>
</p>
