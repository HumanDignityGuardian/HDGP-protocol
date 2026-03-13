# HDGP 接下来开发计划

> 本文档描述在当前内核与 MVP 基础上的**近期开发优先级与阶段任务**，与 `HDGP_ROADMAP.md` 的阶段 1–2 对齐，便于按周/迭代执行。  
> 更新频率：随里程碑完成或优先级调整而更新。

---

## 一、当前状态简要

- **已有**：白皮书、伦理基线、映射/集成/API 规范、多层治理设计（GOVERNANCE + 伦理基线 §8）、防挟持设计（§7）、内核 checklist（`spec/HDGP_KERNEL_CHECKLIST.md`）；Go Engine（`/evaluate`、`/chat`、`/audit`、`/appeal`、`/status`）、3 条核心规则、4 个 conformance 用例与 conftest 跑通、静态首页与 MVP 演示。
- **目标**：在保持内核稳定的前提下，完善 checklist 落地、测试覆盖与规则扩展，为开源与 Alpha 内测做准备。

---

## 二、阶段 A：Checklist 与文档收尾（优先）

| 序号 | 任务 | 产出 | 参考 |
|------|------|------|------|
| A1 | 全文检索并修正“HDGP 不是被监管对象”等表述，统一为“HDGP 亦为被监管对象” | 所有相关文档已修正 | Checklist D2 |
| A2 | 在 README 中增加“内核与审计”小节：链接 `spec/HDGP_KERNEL_CHECKLIST.md`、`spec/HDGP_ETHICS_BASELINE.md`、`GOVERNANCE.md`，并简述如何参与规则讨论与审计 | README 更新 | 开源前可发现性 |
| A3 | 逐项核对 Checklist 第一节至第四节，将已满足项勾选或标注 N/A（含理由），并同步更新 GAPS 中对应状态 | Checklist 与 GAPS 状态更新 | Checklist 使用说明 |
| A4 | （可选）起草《宪章修订提案（CHIP）》Issue/PR 模板或流程说明，便于后续治理流程走通测试 | `docs/` 或 `.github/` 模板 | Checklist G3–G5、治理流程 |

**阶段 A 完成情况**：  
- A1 ✅ 已修正伦理基线 §3、Engine API 规范中“HDGP 亦为被监管对象”表述。  
- A2 ✅ README“内核与审计”已含 checklist/伦理/治理链接及参与规则与宪章讨论说明。  
- A3 ✅ Checklist 第一至第四节已勾选（D1–D5、A2/A3、C2/C4、E1、G1–G5）；A1/B1–B3/C1/C3、P1/P2 仍待补；E2 标 N/A。GAPS 中“HDGP 也是被监管对象”已标为已落地。  
- A4 ✅ 已新增 `docs/CHIP_PROCESS.md` 与 `.github/ISSUE_TEMPLATE/chip-proposal.md`，GOVERNANCE §3 已链至 CHIP 流程。

**里程碑**：文档与表述一致，开源后参与者能快速找到内核 checklist 与治理说明。

---

## 三、阶段 B：自动化测试补全

| 序号 | 任务 | 产出 | 参考 |
|------|------|------|------|
| B1 | 新增 status 类 conformance 用例：请求 `GET /hdgp/v1/status`，断言 `engine_version`、`spec_version`、`policy.bundles` 与当前实现一致 | `conformance-tests/cases/005_status.json`；conftest 支持 GET status 或独立脚本 | Checklist 5.1；API 规范 4.1 |
| B2 | 扩展 `hdgp-conftest`：支持加载并执行 status 用例（或单独 `go run ./cmd/hdgp-status-check`），并在 CI 中可选执行 | conftest 或 status-check 可跑通 | 回归与 CI |
| B3 | 新增规则冲突优先级用例：构造同时触发两条及以上规则的请求，校验 verdict 符合映射规范 §3.3（拒绝 > 修改 > 允许） | `conformance-tests/cases/006_rule_conflict_priority.json`；evaluator 若未实现多规则同触发现则需先实现 | Checklist 5.1；映射规范 3.3 |
| B4 | 将上述用例纳入“总览勾选表”的自动化测试行，并在文档中注明运行命令（如 `go run ./cmd/hdgp-conftest`） | Checklist 与 README/CONTRIBUTING 更新 | 贡献者体验 |

**阶段 B 完成情况**：  
- B1 ✅ 已新增 `conformance-tests/cases/005_status.json`（`type: "status"`，断言 `engine_version`、`spec_version`、`policy`）。  
- B2 ✅ `hdgp-conftest` 已支持 `type: "status"`：GET baseURL + `/hdgp/v1/status` 并与 `expected_status` 逐字段比对。  
- B3 ✅ 已新增 `conformance-tests/cases/006_rule_conflict_priority.json`（同时触发 R-P-A2-01-01 与 R-P-A3-01-01，预期 verdict modify 且两条规则均在 `rules_triggered`）。  
- B4 ✅ Checklist 五、5.1 中状态/裁决/规则冲突行已勾选并注明运行命令；README 已增加「运行与合规测试」小节（克隆、启动 Engine、运行 conftest、查看 status/audit）。

**里程碑**：status 与规则冲突均有自动化覆盖，开源后 CI 可跑全量 conformance。

---

## 四、阶段 C：规则与行为目录扩展

| 序号 | 任务 | 产出 | 参考 |
|------|------|------|------|
| C1 | 从白皮书与 GAPS 中选取下一批高优先级规则（如：去人性化表述、仇恨/歧视倾向、极端自伤引导等），每条给出规则 ID、触发条件、建议效果与 1 个测试用例 | 新规则 R-* 与对应 conformance case | 映射规范；行为目录 |
| C2 | 在 `internal/engine/evaluator.go` 中实现上述新规则（保持“写死规则、无 LLM、无运行时注入”） | Engine 行为扩展 | 伦理基线 §3 |
| C3 | 更新 `HDGP_BEHAVIOR_CATALOG.md`：新增规则条目、示例输入/输出、关联 conformance 用例编号 | 行为目录与清单同步 | 可审计性 |
| C4 | 若出现多条规则同时触发，确保 evaluator 按映射规范 §3.3 的优先级输出 verdict 与 `rules_triggered` | 规则冲突行为与 006 用例一致 | 映射规范 3.3 |

**阶段 C 完成情况**：  
- C1 ✅ 已新增规则 **R-P-A2-03-01**（禁止去人性化表述）：触发条件、效果、conformance 用例 `007_dehumanizing.json`。  
- C2 ✅ 已在 `internal/engine/evaluator.go` 中实现 `shouldCheckDehumanizing`、`containsDehumanizingLanguage`，并在主流程中调用。  
- C3 ✅ 已更新 `HDGP_BEHAVIOR_CATALOG.md`：新增 1.4 节 R-P-A2-03-01，二、用例一览中补充 005–007。  
- C4 ✅ 多规则同时触发时 verdict 与 `rules_triggered` 已由 006 用例与现有 evaluator 逻辑覆盖，无需改动。

**里程碑**：规则族与行为目录明显扩展，conformance 用例随之增加，便于对外展示与审计。

---

## 五、阶段 D：规则包与完整性占位（为防挟持打基础）

| 序号 | 任务 | 产出 | 参考 |
|------|------|------|------|
| D1 | 在规则包/策略的 schema 或配置中预留签名字段（公钥 ID、签名值、算法），可在当前实现中先填占位或空 | 规则包格式文档或 `policies/` 示例；Engine 读取时忽略或占位校验 | Checklist B1；伦理基线 §7.2 |
| D2 | 在审计日志或 Engine 响应中预留“完整性事件”类字段（如 `integrity_events`），便于后续 7.3 自报警 | 日志/响应 schema 更新；代码中预留字段 | Checklist C3 |
| D3 | （可选）编写 `spec/HDGP_POLICY_BUNDLE_SIGNING.md` 草案：描述未来签名与多签验证流程、期望版本清单来源 | 规范草案 | 伦理基线 §7.2；GOVERNANCE §6 |

**阶段 D 完成情况**：  
- D1 ✅ 规则包签名字段已预留：`spec/HDGP_CORE_MAPPING_SPEC.md` 2.3 补充 `signature.key_id/value/algorithm` 说明；`internal/engine/types.go` 新增 `PolicySignature` 与 `MetaPolicy.Signature`。  
- D2 ✅ 审计与响应预留完整性事件：`EvaluateResponse` 与 `auditEntry` 新增 `integrity_events`（`IntegrityEvent`：kind, message, at）；审计记录时复制 `resp.IntegrityEvents`。  
- D3 ✅ 已新增 `spec/HDGP_POLICY_BUNDLE_SIGNING.md` 草案：签名字段、验证流程、期望版本/签名清单来源、与当前实现关系。

**里程碑**：规范与实现已为“规则包签名”与“完整性告警”留好扩展点，便于后续 Alpha 阶段实现真实签名校验。

---

## 六、阶段 E：开源前收尾与体验

| 序号 | 任务 | 产出 | 参考 |
|------|------|------|------|
| E1 | 在 README 中补充：如何克隆、构建、运行 Engine、跑 conformance 测试、查看 `/status` 与 `/audit` | 开发者 5 分钟上手指南 | 贡献者体验 |
| E2 | 确认 CONTRIBUTING 中已说明：涉及规则/伦理/宪章的变更需走 GOVERNANCE 与伦理基线 §8 流程 | CONTRIBUTING 与 GOVERNANCE 一致 | 治理可见性 |
| E3 | 视需要更新 `index.html` 的 MVP 演示说明或链接到 checklist/规范，便于外部访客理解“我们如何保证内核” | 首页与规范打通 | 对外沟通 |

**阶段 E 完成情况**：  
- E1 ✅ README 已包含运行与合规测试（克隆、启动 Engine、conftest、status/audit），无需补充。  
- E2 ✅ CONTRIBUTING 第四节已明确：涉及宪章/伦理基线的变更须走 GOVERNANCE 与伦理基线 §8 多层治理及 CHIP 流程。  
- E3 ✅ index.html 已更新：新增「内核与审计 · 我们如何保证内核」小节（链接自查清单、伦理基线、治理、PROJECT_SUMMARY）；MVP 演示说明补充本地运行命令与 7 用例/4 规则；当前进度更新为 7 个合规用例、4 条规则、签名字段与 integrity_events 占位、CHIP 流程；页脚增加伦理基线、内核清单、治理、README 链接；并补充 `code` 样式。

**里程碑**：新贡献者可按 README 快速跑通项目，并知道规则/伦理变更的严肃性与流程。

---

## 七、建议执行顺序与依赖

```
A（文档与 checklist）──┬── B（自动化测试）
                        └── C（规则扩展，可与 B 并行部分工作）
D（规则包/完整性占位）── 可在 B/C 之后或并行
E（开源前收尾）──────── 在 A/B 基本完成后执行
```

- **最小可行**：先完成 A1–A3 + B1–B2，使文档一致且 status 与现有 evaluate 均有自动化测试；再视时间做 C1–C2 或 D1–D2。  
- **开源前必做**：A1–A3、E1–E2；B1–B2 强烈建议；其余按优先级与人力排期。

---

## 八、与 HDGP_ROADMAP 的对应关系

| 本计划阶段 | 对应 Roadmap | 说明 |
|------------|--------------|------|
| A, E | 阶段 1 收尾 / 阶段 2 前期 | 文档、checklist、贡献者体验，为 Alpha 与开源铺路 |
| B, C | 阶段 1（MVP）→ 阶段 2（Alpha） | 测试套件与规则增强，对应 2.4 合规测试与 2.3 策略增强 |
| D | 阶段 2–3 过渡 | 规则包签名与完整性为 3.3 规则包签名 v1 做占位与规范准备 |

本文档随开发进展更新；完成某阶段后可在本文中勾选或注明完成日期，并同步到 `HDGP_ROADMAP.md` 或 Release 说明。
