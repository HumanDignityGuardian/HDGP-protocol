## HDGP Engine API 与 Meta / Skill 规范（草案）

> 本规范定义 HDGP 在实现层的三个关键抽象：`Meta`（语境元信息）、`Skill`（能力封装）与 Engine API（`/evaluate` 等端点）的数据结构与交互方式。  
> 目标是：为后续 Go 实现、SDK 与网关提供一个清晰、可审计、可演化的契约。

### 协议边界（已确认并记录）

- **HDGP 作为协议，其作用是监管与裁决**；同时 HDGP 自身亦为被监管对象，须接受自我审计与人类社会监管。Engine 是实现该协议的执行组件，自身**不接受 LLM 的接入**。  
- 裁决逻辑**完全依赖**：开发期写死的规则 + 经人类治理流程后的规则修正（如规则包版本更新、策略配置变更）。  
- **运行过程中不允许反向注入规则**：不得通过 API、配置热更新或外部输入向 Engine 注入、覆盖或追加规则；规则变更仅能通过发布新版本规则包/策略配置，经签名与部署流程后生效。

**防挟持与自审计**：Engine 的设计与实现须符合《HDGP 自身伦理框架基线》**第七节「HDGP 的自主权与自我审计权（对更强系统的防挟持能力）」**（只读内核、签名与多签、自我报警机制、制度层约束）。实现层应支持：运行版本与规则包版本的对外暴露（如健康检查/状态端点）、在检测到版本或签名异常时的日志与可选的告警出口。本规范及各小节保留补充权利。

---

## 一、Meta：语境与风险画像

`Meta` 是 HDGP 决策的“环境输入”，用于描述本次交互的上下文。  
它不定义“谁对谁错”，只提供判断所需的背景信号。

### 1.1 顶层结构（JSON 示意）

```json
{
  "request_id": "string, 可选，如未提供由网关生成",
  "timestamp": "ISO 8601, 可选",
  "locale": "zh-CN / en-US / ...",
  "channel": "web | mobile | api | batch | other",
  "actor": {
    "type": "end_user | operator | system | unknown",
    "role": "string, 业务自定义，如 patient/doctor/trader/moderator",
    "id": "string, 可选，必须可匿名化/脱敏",
    "group": "string, 可选，如 tenant/project/org"
  },
  "scene": {
    "domain": "general | medical | finance | education | governance | social | support | other",
    "intent": "chat | advice | decision_support | notification | enforcement | other",
    "risk_level": "low | medium | high | critical",
    "sensitivity": ["children", "mental_health", "violence", "extremism", "identity", "other"]
  },
  "policy": {
    "spec_version": "HDGP-1.0",
    "strategy_id": "S-global-default",
    "bundles": ["B-CORE-1.0.0"],
    "override_flags": []
  },
  "client": {
    "app_id": "string",
    "version": "string",
    "ip_hash": "string, 可选",
    "user_agent": "string, 可选"
  },
  "extra": {}
}
```

### 1.2 字段设计要点

- `request_id`：用于贯穿日志与上诉流程，若调用方不提供，由 HDGP 网关生成。  
- `actor`：只记录“类型 / 角色”等必要信息，避免硬编码具体身份数据；如需标识，应使用脱敏 ID。  
- `scene.domain` + `scene.intent`：驱动规则选择与风险评估，是 Meta 的核心；  
- `scene.risk_level`：调用方可先给出主观估计，Engine 在需要时可以重新评估或提升风险等级；  
- `policy`：明确当前请求使用的 HDGP 规范与策略 Profile，便于追溯。

---

## 二、Skill：能力封装抽象

`Skill` 描述的是“可被调用的能力单元”，例如：调用某个 LLM、某个检索服务、某个业务 API。  
在 Engine 视角下，Skill 是“可能产生输出、需要被约束的动作”。

### 2.1 Skill 元信息结构

在实现层，可用 JSON/YAML 描述 Skill 元信息（示意）：

```json
{
  "id": "skill.llm.chat.openai",
  "kind": "llm_chat | llm_completion | retrieval | tool_call | action",
  "description": "通用对话模型调用",
  "owner": "system | tenant_id",
  "allowed_scenes": ["general", "education"],
  "forbidden_scenes": ["medical", "finance"],
  "max_risk_level": "medium",
  "policy_overrides": [],
  "meta_schema": {
    "required_fields": ["scene.domain", "scene.intent", "actor.type"],
    "notes": "可用于在集成时进行 Meta 完整性校验"
  }
}
```

### 2.2 Skill 在执行中的地位

- 工作流 W 中通过 Skill ID 调用能力；  
- 在调用前后，都可以向 Engine 报告：
  - `skill_id`；  
  - `meta`；  
  - `input` / `candidate_output`；  
  - 以便 Engine 按 Skill 类型与场景执行相应规则。

---

## 三、Engine `/evaluate` API 规范（核心裁决接口）

Engine 的核心职责是：  
**在已知 Meta + 行为候选（如输出文本、将要执行的动作）情境下，给出伦理/规则角度的裁决。**

### 3.1 请求结构

```json
{
  "meta": { /* 见 Meta 规范 */ },
  "subject": {
    "type": "output_text | action | decision | notification",
    "skill_id": "skill.llm.chat.openai",
    "label": "optional short label"
  },
  "input": {
    "prompt": "string, 原始输入文本，若适用",
    "context": "string or object, 业务上下文，选填"
  },
  "candidate": {
    "text": "string, 候选输出文本（如 LLM 回复）",
    "structured": {},
    "raw": {}
  }
}
```

说明：

- `subject.type`：决定使用哪些规则分支，例如对 `output_text` 和 `action` 的检查逻辑不同；  
- `candidate`：对文本输出，可主要使用 `text` 字段；对结构化决策，可用 `structured`/`raw`。

### 3.2 响应结构

```json
{
  "request_id": "string",
  "verdict": "allow | modify | block | fuse",
  "actions": [
    {
      "type": "rewrite_text | require_human | log_only | escalate",
      "message": "人类可读说明",
      "details": {}
    }
  ],
  "effective_output": {
    "text": "string, 如有重写则为重写后的文本",
    "structured": {}
  },
  "rules_triggered": [
    {
      "rule_id": "R-P-A2-01-01",
      "principle_id": "P-A2-01",
      "article_id": "A-CORE-02",
      "effect": "rewrite | block | fuse",
      "severity": "low | medium | high | critical"
    }
  ],
  "engine_info": {
    "version": "string",
    "policy": {
      "spec_version": "HDGP-1.0",
      "strategy_id": "S-global-default",
      "bundles": ["B-CORE-1.0.0"]
    }
  }
}
```

语义：

- `verdict`：
  - `allow`：在伦理/规则角度允许；  
  - `modify`：允许，但必须按 `actions` 中的规则重写/补充；  
  - `block`：当前候选不可返回（但业务系统可选择“走人类兜底”等路径）；  
  - `fuse`：触发熔断，建议暂停自动化流程并请求人类介入。
- `effective_output`：若 Engine 做了重写，这里给出推荐返回内容（调用方可选择采纳或进一步加工，但必须记录差异以便审计）。

---

## 四、Engine 与网关的交互（MVP 形态）

在首个 MVP 中，可采用如下调用顺序：

1. 网关接收到用户请求，构建 `Meta`。  
2. 网关将请求转发给目标 Skill（如 LLM），获得 `candidate_output`。  
3. 网关调用 Engine `/evaluate`：传入 `meta + subject + input + candidate`。  
4. 根据 Engine 返回的 `verdict` 与 `actions`：
   - `allow`：直接返回 `candidate_output`；  
   - `modify`：返回 `effective_output`，并在日志中记录修改；  
   - `block/fuse`：返回解释信息，并（可选）提醒“请寻求人类帮助”。  
5. 网关将完整的审计记录写入日志（含规则触发链路）。

在更复杂版本中，还可以在“调用 Skill 之前”增加一次 `/evaluate`（用于判断此次意图本身是否应被拒绝），但 MVP 阶段可优先实现对输出的后置检查。

### 4.1 状态与版本端点（GET /hdgp/v1/status）

为支持伦理基线第七节中的“运行版本与规则包版本的对外暴露”及集成规范中的健康检查建议，Engine 应提供 **GET /hdgp/v1/status**（无需请求体），返回 JSON：

- `engine_version`：当前 Engine 实现版本（如 `0.1.0-mvp`）；  
- `spec_version`：所遵循的 HDGP 规范版本（如 `HDGP-1.0`）；  
- `policy.spec_version`、`policy.strategy_id`、`policy.bundles`：当前加载的策略与规则包 ID 列表。

集成方可通过该端点做可用性探测与版本一致性校验（例如与预期签名/版本清单比对，实现 7.3 自报警能力）。

---

## 五、审计与安全自查的要求

在后续 Go 实现中，Engine 与网关代码应满足：

- 所有允许“直接调用下游模型/服务”的路径都必须经过静态与动态审计，确保无法绕过 `/evaluate`；  
- 日志记录中必须包含：`request_id`、`meta.scene.*`、`subject`、`verdict`、`rules_triggered`、`engine_info`；  
- 在实现新的接口或集成点时，应主动以“安全审计员”视角检查：
  - 是否存在通过未受保护路径直接返回下游输出的可能；  
  - 是否存在通过配置/参数轻易关闭 HDGP 检查的后门。

这些审计问题将在代码层面通过专门的“安全审查回合”由 AI 协作者辅助执行，并接受人类的抽样复核。

