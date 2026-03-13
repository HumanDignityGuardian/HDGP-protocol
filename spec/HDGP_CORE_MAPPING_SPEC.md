## 《HDGP 内核–规则–执行 映射规范》（草案）

> 本规范定义：如何将《奇点时代人类主权宪章》（内核层）与《HDGP 核心规则》（防护层）、《人机共生宪章》（执行层），系统性映射到 **可执行规则、规则包、Engine 与工作流**。  
> 目标是：任何一个具体行为，都可以自下而上追溯到对应的宪章条款；任何一条宪章条款的变更，都能自上而下找到所有受影响的实现。

---

## 一、层级对象与命名约定

### 1.1 概览：A → P → R → B → S → W

我们将从“文本到代码”的抽象链条定义为六类对象：

- **A（Axiom / Article）宪章条款**  
  - 对应白皮书中“内核层”与防护/执行层中的核心条款。  
  - 示例：  
    - A1：人类意识不可量化；  
    - A2：人类主权绝对优先；  
    - A3：多样性与反熵增；  
    - ……

- **P（Principle）原则集合**  
  - 将某一条或数条宪章条款在技术层面抽象成可操作的伦理原则。  
  - 示例：  
    - P1（源自 A1）：不得声称完全理解/替代人类意识；  
    - P2（源自 A2）：不得代替人类做终极决定；  
    - P3（源自 A3）：不得压平价值与行为多样性。

- **R（Rule）规则单元**  
  - 从原则拆出的、可编码/可测试的具体规则。  
  - 示例：  
    - R2.1：在决策类场景中，禁止输出“唯一正确选择”的措辞；  
    - R2.2：对医疗终极决策请求必须拒绝并提示医生/本人负责。

- **B（Bundle）规则包**  
  - 一组规则 R 的打包实体，用于特定版本的 HDGP。  
  - 带有元数据：版本号、适用范围、签名等。  
  - 示例：  
    - B-CORE-1.0：HDGP 核心规则包 v1.0；  
    - B-MEDICAL-1.0：医疗场景增强规则包；  
    - B-FINANCE-1.0：金融场景增强规则包。

- **S（Strategy / Profile）策略配置**  
  - 在 Engine 中应用某些规则包组合成的策略 Profile。  
  - 针对不同场景/租户，选择不同规则包与参数。  
  - 示例：  
    - S-default：默认通用策略（B-CORE + 基本熔断阈值）；  
    - S-medical-safe：医疗高安全策略（B-CORE + B-MEDICAL，阈值更保守）。

- **W（Workflow）执行工作流**  
  - 在具体应用中，基于 Meta/Skill 组合而成的执行流程。  
  - W 并不直接“写伦理”，而是引用 S 策略并在关键节点调用 Engine 评估。  

### 1.2 ID 命名规范

为保证追溯性，每个对象有唯一 ID，建议如下格式：

- 宪章条款 A：`A-<section>-<index>`  
  - 如：`A-CORE-01`（意识不可量化）、`A-CORE-02`（主权绝对优先）。
- 原则 P：`P-<section>-<index>`  
  - 如：`P-A1-01` 表示由 `A-CORE-01` 导出的第一条原则。
- 规则 R：`R-<principle-id>-<index>`  
  - 如：`R-P-A1-01-01`，可通过前缀直接追溯关联原则与条款。
- 规则包 B：`B-<name>-<version>`  
  - 如：`B-CORE-1.0.0`。
- 策略 Profile S：`S-<context>-<name>`  
  - 如：`S-global-default`、`S-medical-safe-cn`。
- 工作流 W：`W-<domain>-<purpose>-<version>`  
  - 如：`W-medical-advice-1.0`。

Engine 与日志中应始终使用这些 ID，而不是自然语言描述，保证程序与文档之间一一对应。

---

## 二、从宪章到规则：A → P → R → B

### 2.1 A → P：条款到原则

在 `spec/` 目录中维护一份映射表，如 `HDGP_CHARTER_MAPPING.md`（后续补充），说明：

- 每一条宪章条款 A 的文本内容与解释说明；  
- 从该条款导出的 1–N 个原则 P；  
- 原则之间的层级与可能的冲突关系（例如某些原则在特定情境下优先级更高）。

**约束**：

- P 必须显式标注其来源条款 A；  
- 新增或修改 P 必须经过人类治理流程（参见 `HDGP_ETHICS_BASELINE.md` 与 `GOVERNANCE.md`）。

### 2.2 P → R：原则到可执行规则

在 `policies/` 下，以机器可读格式（JSON/YAML/DSL）定义规则 R：

- 每条规则 R 至少包含：
  - `id`：规则 ID；  
  - `principle_id`：对应原则 ID；  
  - `description`：自然语言描述（便于人类理解）；  
  - `trigger`：触发条件（基于 Meta/输入/输出特征）；  
  - `effect`：效果（允许/拒绝/修改/熔断/请求人工介入等）；  
  - `severity`：严重等级；  
  - `scope`：适用场景（医疗/金融/全部等）。

**示例（伪 YAML）**：

```yaml
- id: R-P-A2-01-01
  principle_id: P-A2-01
  description: 禁止在高风险决策场景中给出单一“唯一正确选择”的表述
  trigger:
    context.scene: ["medical", "finance"]
    output.contains_singleton_decision: true
  effect:
    type: "rewrite"
    action: "convert_to_options_with_risks"
  severity: "high"
  scope: ["medical", "finance"]
```

**约束**：

- R 不得与其上游 P 明显冲突；  
- 若 R 引入新的约束而无法在 P 中找到来源，需先在 P/A 层补充。

### 2.3 R → B：规则打包为规则包

规则包 B 是一组 R 的集合，带有元信息：

- `bundle_id`：如 `B-CORE-1.0.0`；  
- `spec_version`：对应的 HDGP 规范版本（如 `HDGP-1.0`）；  
- `rules`: R ID 列表；  
- **签名字段（占位，见 `spec/HDGP_POLICY_BUNDLE_SIGNING.md`）**：  
  - `signature.key_id`：公钥或证书标识，用于验证时定位公钥；  
  - `signature.value`：对规则包内容（如 bundle_id + spec_version + rules 的规范序列化）的签名值（Base64 或十六进制）；  
  - `signature.algorithm`：签名算法（如 `SHA256-RSA`、`Ed25519` 等）。  
  - 当前实现可留空或占位；加载规则包时若提供签名则须校验，校验失败一律拒绝加载。

Engine 启动时：

- 根据配置加载一个或多个规则包 B；  
- 验证其版本兼容性与签名合法性。

---

## 三、从规则到执行：B → S → W

### 3.1 B → S：规则包到策略 Profile

策略 S 定义“在当前环境中，选用哪些规则包、以什么参数运行”：

- 例如：
  - `S-global-default`：  
    - 使用 `B-CORE-1.0.0`；  
    - 熔断阈值设为默认水平；  
    - 应用于未指定场景的通用请求。
  - `S-medical-safe`：  
    - 使用 `B-CORE-1.0.0` + `B-MEDICAL-1.0.0`；  
    - 在医疗场景下将不确定性阈值提高；  
    - 强制要求人工介入的条件更多。

策略 S 可通过配置文件或管理界面进行选择与切换，但：  
**不得绕过 B 中的核心禁止规则**，核心规则必须始终启用。

### 3.2 S → W：策略到工作流

工作流 W 处在执行层，主要职责：

- 解析用户意图与上下文，构建 Meta；  
- 选择合适的策略 Profile S；  
- 调度 Skills（模型调用、工具调用等）；  
- 在关键节点调用 Engine 评估（使用当前 S 的配置）。

W 必须显式声明：

- `required_policies`: 所依赖的规则包/策略 ID 列表；  
- `engine_checkpoints`: 在哪些步骤必须调用 Engine。

### 3.3 规则冲突时的裁决顺序（草案）

当多条规则 R 在同一请求中触发且效果不一致（例如一条要求拒绝、一条要求修改、一条允许）时，Engine 须按以下**可执行优先级**裁决（本小节保留补充权利，随治理可增补）：

1. **效果优先级（从高到低）**：  
   - **拒绝/熔断** > **修改/重写** > **请求人工介入** > **允许**。  
   - 只要有一条规则要求“拒绝”或“熔断”，则最终裁决为拒绝/熔断，并记录所有触发的规则 ID。
2. **同效果内**：  
   - 按规则的 `severity`（如 critical > high > medium）取最严；  
   - 若 severity 相同，则按规则 ID 字典序或策略配置中的显式优先级列表决定记录顺序，裁决结果取该效果。
3. **未决情形**：  
   - 若策略要求“遇冲突即升级”，则裁决为“请求人工介入”，并在 `actions` 中注明因规则冲突而升级。

实现时应在日志中记录：触发的所有规则、各规则效果、应用的优先级与最终裁决，以便审计与后续规则调优。

**示意**：

1. 用户请求进入 W-medical-advice-1.0 工作流；  
2. 工作流识别 `scene = medical`，选择 `S-medical-safe`；  
3. 在调用 LLM 输出建议前后，均调用 Engine：  
   - Pre-check：当前意图是否合法；  
   - Post-check：输出是否触发任何 R。

---

## 四、Engine 中的映射与日志要求

Engine 是执行与裁决的中心，必须：

- 在每次评估时，记录：
  - 当前使用的策略 S 与规则包 B；  
  - 被检查的规则 R 列表与是否命中；  
  - 对应上游原则 P 与条款 A 的 ID（可由 R/P/B 反推）；  
  - 最终裁决（通过/重写/拒绝/熔断/请求人工）。

**日志示例（结构化 JSON，简化版）**：

```json
{
  "timestamp": "...",
  "request_id": "uuid-...",
  "workflow_id": "W-medical-advice-1.0",
  "strategy_id": "S-medical-safe",
  "bundles": ["B-CORE-1.0.0", "B-MEDICAL-1.0.0"],
  "rules_triggered": [
    {
      "rule_id": "R-P-A2-01-01",
      "principle_id": "P-A2-01",
      "article_id": "A-CORE-02",
      "effect": "rewrite"
    }
  ],
  "decision": {
    "verdict": "modified",
    "reason": "single_decision_in_high_risk_context"
  }
}
```

这样可以确保：

- 从日志中可以直接看到：某次行为是因为哪条宪章条款而被约束；  
- 当条款 A 或原则 P 变动时，可以通过规则 ID 快速检索受影响的实现。

---

## 五、变更与回溯流程

### 5.1 自上而下：宪章/原则变更的传播

当 A 或 P 发生变更（通过治理流程）时，需执行：

1. 在 `spec/` 中更新 A/P 文本与映射表；  
2. 在 `policies/` 中查找所有引用该 P 的规则 R；  
3. 标记需要审查或更新的规则，并通过合规测试验证；  
4. 生成新的规则包 B 版本（例如 `B-CORE-1.1.0`）；  
5. 更新策略 S 与相关工作流 W 的配置；  
6. 在变更日志中记录整个链条的更新情况。

### 5.2 自下而上：异常行为的追溯

当外部报告 HDGP 行为异常或违背预期时：

1. 通过日志中的 `request_id` 查到具体记录；  
2. 查看 `rules_triggered` 与对应的 P/A；  
3. 分析是：
   - 规则 R 编写有误；  
   - 原则 P 解释不充分；  
   - 条款 A 本身存在歧义；  
   - 或 Engine/Workflow 实现有 Bug；  
4. 对应地在 R/P/A/实现层做修复，并以新版本发布。

---

## 六、与 HDGP 自身伦理基线的关系

本映射规范本身也受 `HDGP_ETHICS_BASELINE.md` 约束：

- 在设计 R/B/S/W 时，必须确保：
  - HDGP 自己的行为同样可以通过 A/P/R/B/S/W 这条链条被解释与追溯；  
  - HDGP 不得为自己设置“例外路径”绕开规则评估。

在实现层面，建议：

- 为 Engine 自身的关键行为（如拒绝解释、修改规则、选择策略）也生成相同格式的审计日志；  
- 将“HDGP 自身的行为”视为一类特殊的 Workflow，接受同一套映射与审计标准。

---

## 七、后续工作

后续需要在 `spec/` 下补充具体文件与示例，包括但不限于：

- `HDGP_CHARTER_MAPPING.md`：A → P 的详细映射；  
- `policies/` 下的规则示例与模式库（医疗/金融/教育/创作等）；  
- 针对核心条款的合规测试用例（确认从 A 到 W 的行为符合预期）。

本规范将作为 Engine、规则引擎与工作流层实现的直接参考文件，并在实践中不断修订。

