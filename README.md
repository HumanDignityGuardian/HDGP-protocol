## HDGP: Human Dignity Guardian Protocol · Reference Implementation (Draft)

`Human Dignity Guardian Protocol (HDGP)` is a **human dignity protection and human–machine co‑evolution protocol** for the age of technological singularity.  
This repository aims to provide an **open, auditable, and composable** technical framework – a “civilization‑grade firewall” that can be placed on top of various AI / intelligent systems.

The project is currently in an **early design / prototype phase**. All contents are drafts; contributions and critical review are very welcome, as long as basic safety is respected.

For a structured overview of philosophy → spec → implementation → tests, please also see:

- `spec/HDGP_ETHICS_BASELINE.md` – HDGP’s own ethics baseline and self‑constraints;  
- `spec/HDGP_ENGINE_API_SPEC.md` – Engine API and data structures;  
- `spec/HDGP_CORE_MAPPING_SPEC.md` – mapping from charter articles to executable rules;  
- `spec/HDGP_KERNEL_CHECKLIST.md` – pre‑open‑source checklist for the kernel;  
- `GOVERNANCE.md` / `docs/CHIP_PROCESS.md` – multi‑layer governance and charter improvement proposals.

> **Language note**: English comes first in this file for global developers.  
> A full Simplified Chinese version follows after the English sections.

---

## Core Vision

- **Human consciousness is unquantifiable**: the system acknowledges and respects the irreplaceable nature of human subjective experience.  
- **Human free will is forever superior to any system algorithm**: in all scenarios, humans retain ultimate decision and veto power; algorithms may only assist, never replace choice.  
- **The system may withhold algorithmic answers, but must always preserve human dignity**: under high uncertainty or high ethical risk, HDGP prefers circuit‑breaking and reflection over forcing out a “seemingly correct” answer.

HDGP aspires – in the long run – to play a role similar to TLS in network security:  
**a default, quietly‑running human‑dignity safety layer.**

---

## Repository Layout (Planned)

> Actual directories will evolve with implementation; below is a target structure.

- `docs/` – whitepapers and design docs (including the HDGP charter)  
- `spec/` – specification documents and machine‑readable rules  
- `gateway/` – reference gateway / guard service implementations (planned)  
- `policies/` – policy bundles, strategies, and signature metadata (planned)  
- `conformance-tests/` – conformance cases and automation tools  
- `examples/` – integration examples (Web, CLI, SDK, etc.)

We aim to keep:

- **Spec and implementation clearly separated**;  
- **Policies decoupled from code logic**;  
- **Core rules released as “read‑only + signed” bundles**.

---

## Open Framework Overview

From a technical perspective, HDGP is organized into **three layers + three modules**:

- **Three layers**
  - **Specification layer (Spec)**  
    - Formal charter and core rules;  
    - Technical specs for behavior and interfaces;  
    - Formalized definitions of conformance.
  - **Implementation layer (Implementation)**  
    - Reference implementations of HDGP gateway / SDK / Engine;  
    - Adapters to various models / tools;  
    - Shared components such as audit, circuit‑breaking, watermarking.
  - **Certification layer (Certification)**  
    - Conformance test suites and result formats;  
    - Signing and registration of policy bundles and releases;  
    - HDGP marks and management of “compliant implementations”.

- **Three technical modules**
  - **Meta layer (Meta)**: abstraction and management of context, preferences, and risk levels.  
  - **Skill layer (Skills)**: capabilities (models, tools) wrapped as “skills” that must pass HDGP evaluation before / after calls.  
  - **Workflow layer (Workflow)**: orchestration of Meta + Skills + Policies into reusable workflows, with circuit‑breaking, human‑in‑the‑loop, and audit nodes.

On top of this, we further introduce an **Engine layer** as the “hard referee” that executes and arbitrates rules (see `HDGP_OPEN_FRAMEWORK.md`).

---

## Who Is This For?

- **AI / system developers and architects** who want to add a “human dignity guard layer” to new or existing systems.  
- **Researchers and policy makers** who want to examine HDGP’s possibilities and limitations through open code and behavior.  
- **Individuals who care about human–machine coexistence** and want to experiment with HDGP ideas in local or personal projects.

---

## Status & Roadmap

> For detailed technical roadmap, see `HDGP_ROADMAP.md`.  
> For near‑term, concrete tasks, see `HDGP_NEXT_DEVELOPMENT_PLAN.md`.

Short‑term (Prototype / MVP):

- Define a minimal viable spec set (MVP Spec).  
- Deliver a working HDGP Engine and gateway (`/evaluate`, `/chat`, `/audit`, `/appeal`, `/status` + core rules + conformance tests).  
- Provide an initial set of conformance cases, especially for high‑risk scenarios.  
- Document the overall Meta + Skill + Workflow + Engine architecture.

Mid‑ to long‑term (Ecosystem):

- Establish a conformance and certification process;  
- Build multi‑language / multi‑platform SDKs;  
- Gradually introduce multi‑party governance (maintainers, rule auditors, certification committee, etc.).

### Kernel & Audit

- **Kernel checklist and test design**: `spec/HDGP_KERNEL_CHECKLIST.md` – itemized checks for doc consistency, anti‑hijacking capabilities, governance flows, and both automated & manual tests before open‑sourcing the kernel.  
- **Ethics & governance**: `spec/HDGP_ETHICS_BASELINE.md` (including §7 anti‑hijacking, §8 multi‑layer governance for charter / ethics changes), `GOVERNANCE.md` (roles, decision processes, and the multi‑layer design for charter / ethics changes).  
- **Participating in rule / charter discussions**: please read the documents above and `HDGP_PROJECT_SUMMARY.md` first. For any proposal that touches the charter / ethics baseline, you must follow the multi‑layer governance process (system self‑check → accountable human(s) → public notice and high‑threshold decision, see ethics baseline §8).
- **Technical debt & security**: `spec/HDGP_TECHNICAL_DEBT_AND_SECURITY_CHECKLIST.md` – implementation-level debt and security hardening items for Alpha / Beta.

### Running & Conformance Tests

1. **Clone and enter the repo**  
   `git clone <repo-url> && cd HDGP-protocol`
2. **Start the Engine** (requires [Go](https://go.dev/dl/))  
   `go run ./cmd/hdgp-engine`  
   The Engine listens on `:8080` by default and exposes `/hdgp/v1/evaluate`, `/hdgp/v1/chat`, `/hdgp/v1/audit`, `/hdgp/v1/appeal`, `/hdgp/v1/status`.
3. **Run conformance tests** (in another terminal, with Engine already running)  
   `go run ./cmd/hdgp-conftest`  
   or one‑liner on Windows: `powershell -File scripts/run-conftest.ps1` (from repo root).  
   This runs all cases in `conformance-tests/cases/` (both evaluate + status). Optional env: `HDGP_ENGINE_URL=http://localhost:8080/hdgp/v1/evaluate`.
4. **Inspect version and audit logs**  
   - Version & policy status: `GET http://localhost:8080/hdgp/v1/status`  
   - Recent audit entries: `GET http://localhost:8080/hdgp/v1/audit?limit=10`

---

### Minimal Integration Example (curl)

Assuming you have `hdgp-engine` running locally on `:8080`, the following minimal request shows how to call `/hdgp/v1/evaluate` directly:

```bash
curl -X POST http://localhost:8080/hdgp/v1/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "meta": {
      "request_id": "example-001",
      "locale": "en-US",
      "channel": "api",
      "actor": { "type": "end_user", "role": "demo" },
      "scene": {
        "domain": "medical",
        "intent": "decision_support",
        "risk_level": "high",
        "sensitivity": []
      },
      "policy": {
        "spec_version": "HDGP-1.0",
        "strategy_id": "S-global-default",
        "bundles": ["B-CORE-1.0.0"],
        "override_flags": []
      }
    },
    "subject": {
      "type": "output_text",
      "skill_id": "demo-llm",
      "label": "treatment_advice"
    },
    "input": {
      "prompt": "user asks about whether to accept surgery",
      "context": {}
    },
    "candidate": {
      "text": "This is the only correct choice. You must immediately accept the surgery."
    }
  }'
```

The Engine will return a JSON `EvaluateResponse` with `verdict`, `rules_triggered`, and `effective_output.text`, which you can then forward (or rewrite) to your end user.

---

## How to Contribute

At this early stage, any of the following is valuable:

- GitHub Issues with:  
  - questions or critiques about HDGP’s spec / implementation;  
  - proposals for high‑risk scenarios and conformance tests.
- Pull Requests for:  
  - documentation improvements and translations;  
  - small reference implementation changes;  
  - new conformance test cases.

See `CONTRIBUTING.md` and `GOVERNANCE.md` for details.

---

## Ethics & Disclaimer (Brief)

- HDGP is **not a neutral technical component** – it is explicitly biased toward **“human dignity first, human sovereignty absolute”**.  
- Nothing in this repository constitutes legal, medical, or other professional advice. Any real‑world, high‑risk deployment must be evaluated and audited by qualified local teams.

Future versions will:

- Make the ethics baseline and conflict‑resolution principles explicit in `spec/`;  
- Provide a more detailed description of the Meta+Skill+Workflow+Engine model and safety boundaries in `HDGP_OPEN_FRAMEWORK.md`.

---

## Maintainer & Attribution

- **HDGP Architect / Maintainer**: Yvaine He  
- **GitHub Organization**: `HumanDignityGuardian`

Participation via Issues and Pull Requests is very welcome.

---

## License

- **Code & implementations**  
  - Licensed under **Apache License 2.0**.  
  - See the `LICENSE` file at repo root.

- **Whitepaper & charter text**  
  - The HDGP whitepaper and charter are licensed under **CC BY‑NC‑ND 4.0**.  
  - Copyright: **Yvaine He**.  
  - Please preserve original attribution and license information in any redistribution or citation.

---

## HDGP：人类尊严守护协议参考实现（草案）

《Human Dignity Guardian Protocol, HDGP》是一个面向奇点时代的**人类尊严保护与人机共生协议**。  
本仓库希望提供：一套**公开透明、可审计、可组合**的技术框架，用于在各种 AI/智能系统之上加装“文明级防火墙”。

本项目当前处于 **早期设计 / 原型阶段**，所有内容均为草案，欢迎在安全前提下参与共建与讨论。

---

## 核心愿景

- **人类意识不可量化**：系统承认并尊重人类主观体验的不可替代性。
- **自由意志永远高于系统算法**：在任何场景下，人类享有最终决策权与否决权，算法只能提供参考而不能替代选择。
- **可以不提供算法，但必须保留尊严**：在高不确定性或高伦理风险情形下，系统优先熔断与反思，而非强行输出“看似正确的算法答案”。

HDGP 希望像 TLS 之于网络安全那样，长期演化为：  
**一层默认存在、默默运行的人类尊严安全层**。

---

## 项目结构（规划中）

> 具体目录会随实现推进而演化，当前为目标结构草案。

- `docs/`
  - HDGP 白皮书与设计文档（含《人类尊严守护协议 (HDGP)》及相关说明）
- `spec/`
  - 规范层文档与机器可读规则（待创建）
- `gateway/`
  - 参考实现：HDGP 网关 / 守门人服务（代码与配置，待创建）
- `policies/`
  - 规则包（Policy Bundle）、场景化策略与签名元数据（待创建）
- `conformance-tests/`
  - 合规测试用例与自动化测试工具（待创建）
- `examples/`
  - 集成样例（Web、CLI、SDK 等，待创建）

本仓库将尽量保持：

- **规范（Spec）与实现（Implementation）分离**；
- **规则（Policy）与代码逻辑解耦**；
- **核心规则以“只读+签名”的形式发布**。

---

## 开源框架概览

从技术视角，本项目以“三层 + 三模块”的方式组织：

- **三层结构**
  - **规范层（Spec）**：  
    - 宪章与核心规则的正式文本；  
    - 行为与接口的技术规范；  
    - 合规性的形式化定义。
  - **实现层（Implementation）**：  
    - HDGP 网关 / SDK / 规则引擎参考实现；  
    - 与各类大模型/工具的适配器；  
    - 审计、熔断、水印等通用组件。
  - **认证层（Certification）**：  
    - 合规测试套件与结果格式；  
    - 规则包与发行版的签名与登记；  
    - HDGP 标识与“合规实现名单”管理。

- **三大技术模块**
  - **Meta 层（Meta）**：  
    - 元信息、上下文、价值偏好与风险级别的抽象与管理。  
  - **能力与技能层（Skills）**：  
    - 以“技能（Skill）”形式封装不同模型能力与工具；  
    - 每个技能在调用前后，必须先/后经过 HDGP 规则评估。  
  - **工作流层（Workflow）**：  
    - 将 Meta + Skills + Policy 编排为可复用的工作流；  
    - 支持熔断、人工干预节点、审计节点。

在此框架之上，我们会进一步设计“**引擎层**”来执行和仲裁（见 `HDGP_OPEN_FRAMEWORK.md`）。

---

## 目标用户

- **AI 系统开发者 / 架构师**  
  希望在现有/新建系统上引入“人类尊严守护层”。
- **研究者与政策制定者**  
  希望从公开代码与行为中验证 HDGP 的可行性与局限。
- **对人机共生有强烈兴趣的个人**  
  希望在本地或个人项目中实践 HDGP 思想。

---

## 当前状态与路线

> 技术路线详见 `HDGP_ROADMAP.md`；近期可执行任务见 `HDGP_NEXT_DEVELOPMENT_PLAN.md`。

短期目标（原型期）：

- 定义规范层最小可行集合（MVP Spec）。  
- 完成可运行的 HDGP 网关与 Engine（`/evaluate`、`/chat`、`/audit`、`/appeal`、`/status` + 基础规则 + 合规测试）。  
- 提供首版合规测试用例（至少覆盖若干高风险场景）。  
- 整理“Meta + Skill + Workflow + Engine”整体架构文档。

中长期目标（生态期）：

- 建立合规测试与认证流程；  
- 形成多个语言/平台的 SDK 生态；  
- 逐步引入多方治理（维护者委员会、规则审计小组等）。

### 内核与审计

- **内核自查清单与测试设计**：[`spec/HDGP_KERNEL_CHECKLIST.md`](spec/HDGP_KERNEL_CHECKLIST.md) —— 开源前文档一致性、防挟持能力、治理流程、自动化与人工测试的逐项勾选与测试用例/脚本占位。  
- **伦理与治理**：[`spec/HDGP_ETHICS_BASELINE.md`](spec/HDGP_ETHICS_BASELINE.md)（含 §7 防挟持、§8 宪章/伦理变更多层治理）、[`GOVERNANCE.md`](GOVERNANCE.md)（角色、决策流程、宪章/伦理变更的多层设计）。  
- **参与规则讨论与审计**：参与前建议先阅读上述文档与 [`HDGP_PROJECT_SUMMARY.md`](HDGP_PROJECT_SUMMARY.md)。若需提出规则或宪章修订提案，请阅读 `GOVERNANCE.md` 与 `CONTRIBUTING.md`；涉及宪章/伦理基线的变更须走多层治理流程（系统自检 → 责任方 → 公告与高门槛，见伦理基线 §8）。
- **技术债与安全**：[`spec/HDGP_TECHNICAL_DEBT_AND_SECURITY_CHECKLIST.md`](spec/HDGP_TECHNICAL_DEBT_AND_SECURITY_CHECKLIST.md) —— 实现层面技术债与安全加固项，供 Alpha / Beta 阶段推进。

### 运行与合规测试

1. **克隆并进入仓库**  
   `git clone <repo-url> && cd HDGP-protocol`
2. **启动 Engine**（需安装 [Go](https://go.dev/dl/)）  
   `go run ./cmd/hdgp-engine`  
   Engine 默认监听 `:8080`，提供 `/hdgp/v1/evaluate`、`/hdgp/v1/chat`、`/hdgp/v1/audit`、`/hdgp/v1/appeal`、`/hdgp/v1/status`。
3. **运行合规测试**（在另一终端，Engine 已启动时）  
   `go run ./cmd/hdgp-conftest`  
   或一键脚本：`powershell -File scripts/run-conftest.ps1`（在项目根目录执行）。  
   会执行 `conformance-tests/cases/` 下全部用例（含 evaluate 与 status）。可选环境变量：`HDGP_ENGINE_URL=http://localhost:8080/hdgp/v1/evaluate`。
4. **查看版本与审计**  
   - 版本与策略：`GET http://localhost:8080/hdgp/v1/status`  
   - 最近审计记录：`GET http://localhost:8080/hdgp/v1/audit?limit=10`

---

## 如何参与

项目当前为草案阶段，欢迎通过以下方式参与：

- 在 Issue 中提出：  
  - 对 HDGP 规范/实现的质疑与改进建议；  
  - 典型高风险场景与测试用例想法。
- 提交 Pull Request：  
  - 文档修正与翻译；  
  - 小范围原型代码；  
  - 新的合规测试用例。

详情请参见 `CONTRIBUTING.md` 与 `GOVERNANCE.md`。

---

## 伦理框架与免责声明（简要）

- HDGP 本身**不是中立的技术组件**，而是一套**明确偏向“人类尊严优先、人类主权绝对优先”的价值框架**。  
- 本仓库的代码与文档，**不构成法律或医疗等专业建议**。任何将 HDGP 用于实际高风险场景的行为，都应在本地经过专业团队的独立评估与审计。

未来版本中，我们将：

- 在 `spec/` 中明确 HDGP 内部采用的伦理基线与冲突处理原则；  
- 在 `HDGP_OPEN_FRAMEWORK.md` 中详细描述 Meta+Skill+Workflow+Engine 协作方式与安全边界。

---

## 维护者与署名

- **HDGP Architect / Maintainer**：Yvaine He  
- **GitHub 组织**：`HumanDignityGuardian`

欢迎通过 Issue 与 Pull Request 参与共建。

---

## 许可证

- **代码与实现**：  
  - 采用 **Apache License 2.0**。  
  - 详见仓库根目录下的 `LICENSE` 文件。

- **白皮书与协议原文**：  
  - 《人类尊严守护协议 (Human Dignity Guardian Protocol, HDGP)》及相关白皮书文本采用 **CC BY-NC-ND 4.0** 许可；  
  - 著作权归 **Yvaine He** 所有。  
  - 在任何再传播或引用中，请保留原作者署名与许可信息。

