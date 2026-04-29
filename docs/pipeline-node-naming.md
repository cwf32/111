# Pipeline 节点命名规范

本文档用于规范 MDA 项目中 MaaFramework Pipeline 节点的命名方式，提升节点可读性、可维护性与跨模块一致性。

Pipeline 节点名是 JSON 根对象中的 key，会被 `next`、`on_error`、`target`、`anchor`、`And` / `Or` 识别条件、`pipeline_override` 等字段引用。因此节点重命名必须同步更新所有引用。

## 基本原则

节点名必须使用 **PascalCase**。

推荐格式为：

```text
<Domain><ActionOrObject><Role>
```

其中：

- `Domain` 表示所属功能域或模块，例如 `Shop`、`Arena`、`SimulationRoom`、`DailyRewards`。
- `ActionOrObject` 表示节点处理的动作、页面、对象或业务目标。
- `Role` 表示节点在流程中的功能角色，例如 `Flow`、`Enter`、`OnPage`、`Visible`、`Detected`、`Confirm`、`Selected`、`Claim` 等。

示例：

```text
ShopEnterCommonShop
ShopOnCommonShopPage
SimulationRoomSelectBuff
SimulationRoomConfirmBuff
ArenaQuickBattleAvailable
DailyRewardsClaimMission
```

## 禁止事项

节点名不得使用以下形式：

```text
_BeginSimulation1
25Check
clickReward
shop_enter
Flag_In_Shop
```

具体禁止规则：

1. 不得以下划线 `_` 开头。
2. 不得以数字开头。
3. 不得使用 `snake_case`、`camelCase` 或混合分隔符。
4. 不得使用无业务语义的临时编号，例如 `Node1`、`Check2`。
5. 不得仅用过于泛化的名称，例如 `Confirm`、`Check`、`Click`，除非是明确的全局通用节点。
6. 不推荐使用 `FlagInX` 作为新节点命名，除非是为了兼容旧节点。

## 节点角色命名

### 入口节点

任务或模块入口节点使用：

```text
<Domain>Main
```

或在兼容现有任务入口时保留：

```text
<Domain>
```

推荐：

```text
ShopMain
ArenaMain
SimulationRoomMain
DailyRewardsMain
```

现有入口节点如 `Shop`、`Arena`、`SimulationRoom` 可逐步迁移，不要求一次性改名。

### 流程节点

只负责组织后继节点、不直接识别或点击的节点，使用：

```text
<Domain><Subtask>Flow
```

示例：

```text
DailyRewardsEmailFlow
ShopPurchaseItemFlow
SimulationRoomBiosSettingFlow
ArenaBattleFlow
```

适用场景：

```json
{
    "ShopPurchaseItemFlow": {
        "next": [
            "ShopPurchaseDialogVisible",
            "[JumpBack]CommonConfirmAction",
            "CommonConfirmReward"
        ]
    }
}
```

### 进入页面节点

用于点击入口并进入某个页面的节点，使用：

```text
<Domain>Enter<Page>
```

示例：

```text
ShopEnterCommonShop
ArenaEnterRookie
SimulationRoomEnter
DailyRewardsEnterMission
```

不推荐：

```text
EnterShop
EnterArena
EnterSimulationRoom
```

原因是缺少模块前缀，在全局 pipeline 中不够明确。

### 页面状态节点

用于判断当前是否处于某页面、某界面、某弹窗的节点，使用：

```text
<Domain>On<Page>Page
```

或：

```text
<Domain><Object>Visible
```

示例：

```text
ShopOnCommonShopPage
ArenaOnRookiePage
SimulationRoomOnMainPage
CommonRewardDialogVisible
```

不推荐新增：

```text
FlagInShop
FlagInArena
FlagInSimulationRoom
```

旧节点可保留，但新节点应优先使用 `On...Page` / `Visible`。

### 纯检测节点

只负责识别某个元素、状态、文本、红点，不执行动作的节点，使用：

```text
<Domain><Object>Detected
<Domain><Object>Visible
<Domain><Object>Available
<Domain><Object>Claimed
<Domain><Object>Selected
```

示例：

```text
SimulationRoomBeginRedDotDetected
SimulationRoomBeginTextDetected
ArenaQuickBattleAvailable
ShopGemPurchaseAvailable
DailyRewardsMissionClaimed
```

根据语义选择后缀：

| 后缀        | 含义                       |
| ----------- | -------------------------- |
| `Detected`  | 识别到某个图像、文本或特征 |
| `Visible`   | UI 元素可见                |
| `Available` | 功能、按钮、次数可用       |
| `Claimed`   | 奖励或任务已领取           |
| `Selected`  | 选项已选中                 |
| `Completed` | 流程或收集已完成           |
| `Exhausted` | 次数已耗尽                 |

### 点击/选择节点

执行点击、选择、领取等动作的节点，使用动词前置：

```text
<Domain>Click<Object>
<Domain>Select<Object>
<Domain>Claim<Object>
<Domain>Purchase<Object>
<Domain>Open<Object>
<Domain>Close<Object>
```

示例：

```text
CommonClickMax
ShopPurchaseFreeGoods
SimulationRoomSelectBuff
SimulationRoomOpenBiosSetting
CommonClosePage
DailyRewardsClaimMission
```

不推荐：

```text
ClickMax
PassClick
FreeRecruitClick
```

除非该节点是历史节点且引用较多。

### 确认节点

确认弹窗、确认奖励、确认操作使用：

```text
<Domain>Confirm<Object>
```

示例：

```text
CommonConfirmAction
CommonConfirmReward
SimulationRoomConfirmBuff
SimulationRoomConfirmBattleEnd
ShopConfirmPurchase
```

不推荐：

```text
Confirm
ActionConfirm
RewardConfirm
ConfirmEnd
```

### 滚动/滑动节点

滚动、滑动节点使用：

```text
<Domain>Scroll<Direction>
<Domain>Swipe<Object>
```

示例：

```text
CommonScrollUp
CooperationSwipeBanner
SimulationRoomScrollOverclockOptionsUp
AdviseScrollEpisodeRewardDown
```

不推荐：

```text
ScrollUp
SlideBanner
SimulationRoomScrollUp
```

### 结束节点

通用结束节点可以保留：

```text
EndTask
```

如果未来希望完全统一，可迁移为：

```text
CommonEndTask
```

由于 `EndTask` 被大量引用，建议短期内保留。

## Common 节点命名

`Common` 节点是全局复用节点，应当显式带有 `Common` 前缀，避免和业务模块节点混淆。

推荐：

```text
CommonConfirmReward
CommonConfirmAction
CommonClosePage
CommonClickMax
CommonClickBlank
CommonGoBack
CommonScrollUp
CommonEndTask
```

对于导航类 Common 节点，可使用：

```text
NavigationEnterHall
NavigationEnterArk
NavigationClickHomeButton
NavigationClickHall
NavigationOnArkPage
NavigationArkVisibleInHall
```

如果不想引入 `Navigation` 域，也可统一使用 `Common`：

```text
CommonEnterHall
CommonEnterArk
CommonClickHomeButton
CommonOnArkPage
```

二者择一即可，不建议混用。

## 兼容旧节点的迁移策略

为降低风险，节点重命名应分批进行。

### 第一批：修正硬性不规范节点

优先处理不符合 PascalCase 的节点：

| 旧名                | 新名                                     |
| ------------------- | ---------------------------------------- |
| `_BeginSimulation1` | `SimulationRoomBeginRedDotDetected`      |
| `_BeginSimulation2` | `SimulationRoomBeginTextDetected`        |
| `25Check`           | `SimulationRoomOverclockLevel25Selected` |

### 第二批：整理 Common 节点

建议迁移：

| 旧名            | 新名                  |
| --------------- | --------------------- |
| `RewardConfirm` | `CommonConfirmReward` |
| `ActionConfirm` | `CommonConfirmAction` |
| `PageClose`     | `CommonClosePage`     |
| `ClickMax`      | `CommonClickMax`      |
| `BlankClick`    | `CommonClickBlank`    |
| `GoBack`        | `CommonGoBack`        |

`EndTask` 可暂时保留。

### 第三批：按模块逐步迁移

建议迁移顺序：

1. `SimulationRoom`
2. `Shop`
3. `Arena`
4. `DailyRewards`
5. 其他模块

每次迁移一个模块，避免跨模块大规模重命名导致引用遗漏。

## 重命名检查清单

每次重命名节点后，必须检查以下引用位置：

1. 节点定义 key。
2. `next`。
3. `on_error`。
4. `[JumpBack]NodeName`。
5. `[Anchor]AnchorName` 相关节点引用。
6. `target` 中的节点名引用。
7. `anchor` 对象中的节点名引用。
8. `recognition.type = "And"` 中的 `all_of`。
9. `recognition.type = "Or"` 中的 `any_of`。
10. `pipeline_override`。
11. 任务配置 `assets/tasks/**/*.json`。
12. 安装产物或镜像目录中的重复资源，如果项目要求同步维护。

## 命名示例

### 推荐

```json
{
    "SimulationRoomMain": {
        "next": [
            "SimulationRoomEnter",
            "[JumpBack]NavigationEnterArk"
        ]
    },
    "SimulationRoomEnter": {
        "desc": "进入模拟室",
        "recognition": {
            "type": "OCR"
        },
        "action": {
            "type": "Click"
        },
        "next": [
            "SimulationRoomOnMainPage",
            "SimulationRoomEnter"
        ]
    },
    "SimulationRoomOnMainPage": {
        "desc": "处于模拟室页面",
        "recognition": {
            "type": "OCR"
        },
        "next": [
            "[JumpBack]SimulationRoomBeginSimulation",
            "[JumpBack]SimulationRoomStartOverclock",
            "EndTask"
        ]
    }
}
```

### 不推荐

```json
{
    "SimulationRoom": {},
    "EnterSimulationRoom": {},
    "FlagInSimulationRoom": {},
    "_BeginSimulation1": {},
    "25Check": {}
}
```

## 总结

MDA Pipeline 节点命名统一采用：

```text
PascalCase + 模块域 + 功能语义 + 角色后缀
```

新增节点应优先表达“节点在流程中的功能”，而不是仅表达底层动作类型。

推荐风格：

```text
ShopEnterCommonShop
ShopOnCommonShopPage
SimulationRoomSelectBuff
SimulationRoomConfirmBuff
DailyRewardsMissionClaimed
ArenaQuickBattleAvailable
CommonConfirmReward
```

避免风格：

```text
FlagInShop
ClickMax
25Check
_BeginSimulation1
Confirm
```
