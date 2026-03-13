## HDGP Contribution Guide

Thank you for your interest in contributing to the **Human Dignity Guardian Protocol (HDGP)**.  
The project is in an early stage; contributions of all kinds (discussion, docs, tests, code) are highly appreciated.

---

## 1. How Can I Help?

You can start with any of the following:

- **Report issues and high‑risk scenarios**  
  - Bug reports, implementation flaws, or documentation mistakes;  
  - Real or hypothetical high‑risk scenarios for improving conformance tests.

- **Improve documentation and specs**  
  - Fix typos or unclear wording;  
  - Propose stricter / clearer formulations for rules;  
  - Write or translate guides and architecture documents.

- **Contribute code and tests**  
  - Add features or fixes to reference implementations;  
  - Extend the conformance test suite;  
  - Provide SDKs / adapters / integration examples.

---

## 2. Filing Issues

- Use a clear title, and try to structure the description as:
  - **Background / scenario** – where and how you encountered the problem;  
  - **Expected behavior** – how HDGP should behave;  
  - **Actual behavior** – what currently happens;  
  - **Potential impact** – any safety / ethics risks.

- For privacy‑sensitive or abuse‑related topics, you may start with a high‑level description and follow maintainer guidance for private discussion if needed.

---

## 3. Pull Request Workflow

### 3.1 Preparation

1. Fork the repo and clone it locally;  
2. Create a feature branch, for example:

```bash
git checkout -b feat/add-medical-scenario-tests
```

3. Install dependencies and run the project locally as described in `README.md` (if applicable).

### 3.2 Coding Guidelines

- Follow existing code style (indentation, naming, comments);  
- Avoid mixing unrelated changes in a single PR (e.g., refactor + new feature);  
- For non‑obvious logic or design choices, add short explanations in PR description or code comments;  
- Whenever possible, add unit / integration tests for new features or rules.

### 3.3 Conformance & Safety

HDGP differs from typical projects: **changes that “relax constraints” may actually weaken the boundary that protects human dignity**.  
Therefore, we have extra requirements:

- If your PR:
  - modifies `spec/` or `policies/`;  
  - adjusts circuit‑breaking logic, uncertainty thresholds, or output filters;  
  - introduces new decision‑affecting features;
- then please document in the PR:
  - motivation;  
  - potential impact on different users / scenarios;  
  - tests you believe should be added or updated.

Maintainers will:

- run automated + conformance tests;  
- involve rule / ethics reviewers when necessary.

### 3.4 Submitting & Review

1. Use meaningful commit messages, for example:

```bash
git commit -m "Add uncertainty-based circuit breaker for healthcare advice"
```

2. Push your branch to your fork and open a Pull Request against the main repo:  
   - briefly describe the change in the title;  
   - explain motivation and whether rules / policies are affected in the description.

3. Wait for review from maintainers / auditors:  
   - they may request changes or additional tests;  
   - please update your PR until it is approved.

---

## 4. Special Requirements for Rule & Charter Changes

**Any change to the charter or ethics baseline** (e.g., whitepaper charter articles, constraints in `spec/HDGP_ETHICS_BASELINE.md`, or §7/§8 governance flows) **must follow the multi‑layer governance process** defined in `GOVERNANCE.md` and ethics baseline §8 (system self‑check → accountable human(s) → public notice + high‑threshold decision), and should preferably use the [Charter Improvement Proposal (CHIP)](docs/CHIP_PROCESS.md) process.  
Such changes may **not** take effect solely based on a single person or a small group signing off.

For changes that directly affect the strength of “human dignity protection” (e.g., lowering circuit‑break thresholds or allowing previously forbidden behaviors), please:

1. First open an Issue as a **proposal**, describing:  
   - limitations of current rules;  
   - your suggested change;  
   - potential benefits and risks.
2. After a public discussion period (length depends on impact), submit the corresponding PR.  
3. In the PR, reference the Issue and clearly state:  
   - final adopted proposal;  
   - feedback that has been incorporated;  
   - remaining uncertainties.
4. For major proposals, rule auditors and project founders may extend discussion or split changes into multiple smaller steps.

---

## 5. Documentation & Translation Contributions

Documentation and translation are critical. They directly influence:

- how external developers understand and use HDGP;  
- how the public perceives HDGP’s values and boundaries.

Suggestions:

- keep terminology consistent (see `README.md` and spec docs);  
- separate value judgments from implementation details where possible;  
- for potentially confusing parts, err on the side of being explicit about “what HDGP can / cannot do”.

---

## 6. Code of Conduct

HDGP naturally touches sensitive topics (ethics, values, risk assessment). We aim for a community that is:

- **Respectful** – of participants from different backgrounds, cultures, and viewpoints;  
- **Rational** – disagreements grounded in facts and logic;  
- **Safe** – avoid publishing detailed abuse vectors or attack paths in public discussions.

See `CODE_OF_CONDUCT.md` for more details.

---

## 7. Contact & Support

In early stages, the main channels are:

- GitHub Issues / Discussions;  
- future mailing lists or community spaces (TBD).

If you are unsure whether an idea is ready as a PR, feel free to start with a short Issue.  
Again, thank you for contributing your time and energy to HDGP.

---

## HDGP 贡献指南（CONTRIBUTING）

感谢你愿意为《人类尊严守护协议（HDGP）》贡献力量。  
本项目处于早期阶段，任何形式的参与（讨论、文档、测试、代码）都非常重要。

---

## 一、我可以做什么？

你可以从以下几类贡献开始：

- **反馈问题与风险场景**  
  - 报告 Bug、实现缺陷或文档错误；  
  - 提供真实或假想的高风险使用场景，用于完善合规测试。

- **改进文档与规范**  
  - 修正错别字、语义不清的地方；  
  - 建议更严格/更清晰的条款表述；  
  - 翻译或撰写使用指南、架构说明。

- **贡献代码与测试**  
  - 为参考实现添加新特性或修复问题；  
  - 为合规测试套件增加新的用例；  
  - 提供 SDK / 适配器 / 集成示例。

---

## 二、提 Issue 的建议

- 使用清晰的标题，并尽量按照以下结构描述：
  - **背景 / 场景**：在哪种情境下遇到问题；  
  - **预期行为**：你认为 HDGP 应该如何表现；  
  - **实际行为**：当前实现的行为；  
  - **可能的影响**：是否存在安全/伦理风险。

- 对于涉及隐私或潜在滥用风险的内容，可以先用概括描述，必要时根据维护者建议改为私下沟通。

---

## 三、提交代码变更（Pull Request）流程

### 1. 准备工作

1. Fork 仓库，并克隆到本地；  
2. 为你的工作创建一个新分支，例如：

```bash
git checkout -b feat/add-medical-scenario-tests
```

3. 安装依赖并按 `README.md` 中的说明在本地运行项目（如果已有实现）。

### 2. 编码规范与基本要求

- 保持代码风格与现有代码一致（缩进、命名、注释风格等）；  
- 避免在单个 PR 中引入过多不相关的修改（例如同时重构和添加功能）；  
- 对复杂逻辑和设计意图，可以在 PR 描述或代码注释中简要说明；
- 尽量为新增功能/规则编写相应的单元测试或集成测试。

### 3. 合规测试与安全考虑

HDGP 与传统项目不同之处在于：**很多看似“放松限制”的修改，实质上可能削弱人类尊严保护的边界**。  
因此，我们对合规性有额外要求：

- 当你的 PR：
  - 修改 `policies/` 或 `spec/` 相关内容；  
  - 调整熔断逻辑、不确定性阈值、输出过滤规则；  
  - 引入影响决策行为的新特性；
- 请务必在 PR 中说明：
  - 变更动机；  
  - 对不同用户/场景可能的影响；  
  - 你认为需要新增或更新的测试用例。

项目维护者会：

- 运行自动化测试与合规测试套件；  
- 在必要时请求规则审计员参与 Review。

### 4. 提交与 Review

1. 提交代码时采用有意义的 Commit 信息，例如：

```bash
git commit -m "Add uncertainty-based circuit breaker for healthcare advice"
```

2. 将分支推送到你的 Fork，并在主仓库创建 Pull Request：  
   - 标题简要描述变更；  
   - 描述中说明：变更内容、动机、是否涉及规则/策略变动。

3. 等待维护者/审计员 Review：  
   - 他们可能会提出修改建议或请求补充测试；  
   - 请根据 Review 意见进行更新，直到 PR 获得通过。

---

## 四、规则与规范变更的特别要求

**涉及宪章条款或伦理基线的变更**（如修改白皮书宪章、`spec/HDGP_ETHICS_BASELINE.md` 中的伦理约束或 §7§8 治理流程），须走 **GOVERNANCE.md** 与 **伦理基线 §8** 所规定的多层治理流程（系统自检 → 责任方指定 → 公告时间窗与高门槛），并优先使用 [宪章修订提案（CHIP）](docs/CHIP_PROCESS.md) 流程。不得仅凭单人或少数人签署即生效。

对于会直接影响“人类尊严保护强度”的变更（例如：降低熔断阈值、允许某些原先被禁止的行为），请遵循以下额外流程：

1. 先在 Issue 中发起“提案（Proposal）”，说明：
   - 当前规则的不足之处；  
   - 你建议的改动方案；  
   - 可能带来的好处与风险。

2. 经过一段时间的公开讨论（视影响大小而定），再提交对应的 PR。

3. PR 中需引用该 Issue，并明确：
   - 最终采用的方案；  
   - 已经纳入的反馈；  
   - 你认为仍存在的不确定性。

4. 对于重大提案，规则审计员与项目发起人可能要求延长讨论期或拆分为多个小步骤实施。

---

## 五、文档与翻译贡献

文档和翻译非常重要，它们直接影响到：

- 外部开发者能否正确理解与使用 HDGP；  
- 公众如何看待 HDGP 的价值与边界。

文档类贡献建议：

- 保持术语一致（参考 `README.md` 与规范文档）；  
- 尽量将价值判断与实现细节分开描述；  
- 对可能被误解的地方，宁可多写一两句，讲清楚“能做什么 / 不能做什么”。

---

## 六、行为准则

HDGP 讨论的主题天然带有敏感性（伦理、价值观、风险评估等），我们希望社区氛围保持：

- **尊重**：尊重不同背景、文化、立场的参与者；  
- **理性**：以事实和逻辑为基础讨论分歧；  
- **安全**：避免在公开场合详细讨论具体滥用手段和攻击路径。

详细内容请参见 `CODE_OF_CONDUCT.md`（待补充）。

---

## 七、联系与支持

在项目早期，主要沟通渠道包括：

- 仓库 Issue / Discussions；  
- 未来可能建立的邮件列表或社区空间。

如果你不确定某个想法是否适合直接以 PR 形式提交，可以先开一个简短的 Issue 进行讨论。  
再次感谢你为 HDGP 付出的时间与心力。

