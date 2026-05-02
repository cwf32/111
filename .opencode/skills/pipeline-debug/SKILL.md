---
name: pipeline-debug
description: Debug and optimize MaaFramework pipeline JSON snippets. Validates against schema specifications, detects structural issues (missing references, circular dependencies), identifies common errors, and provides optimization recommendations. Use when debugging pipeline execution, validating configuration, or improving performance.
license: MIT
compatibility: opencode
metadata:
    version: "1.0"
    project: MDA
    author: MDA Team
allowed-tools: Read Grep Glob
---

## What I do

- Validate pipeline JSON snippets against `pipeline.schema.json` and related schemas
- Analyze node relationships and detect structural issues (missing references, circular dependencies)
- Identify common errors (missing required fields, type mismatches, invalid values)
- **Analyze node roles and semantics via naming conventions** - understand what each node does based on its name
- Detect naming-code mismatches (e.g., node named "Click" but missing action)
- Provide optimization recommendations for performance and maintainability
- Generate corrected pipeline snippets when issues are found
- **To ensure issues are fixed, adding or deleting nodes is permitted** - focus on achieving the correct pipeline behavior rather than preserving the original structure

## When to use me

Use this skill when:

- You have a pipeline JSON snippet that isn't working as expected
- You want to validate pipeline configuration before deployment
- You need to optimize pipeline performance or reduce complexity
- You're debugging pipeline execution issues
- **You want to understand what each node does based on its naming**
- **You need to verify if node names match their actual behavior**
- A user mentions "debug pipeline", "validate pipeline", "optimize pipeline", or "pipeline error"
- A user provides a pipeline fragment and asks "what's wrong with this?" or "how can I improve this?"

## Quick Start

1. **Provide your pipeline snippet** - Paste the JSON fragment you want to debug
2. **I'll analyze it** - Check for errors, validate against schema, and suggest improvements
3. **Get corrected version** - Receive a fixed and optimized pipeline snippet

## Workflow

### Step 1: Read Specifications

Read these files from project `tools/schema/` directory (relative path: `../../../tools/schema/`):

1. **`pipeline.schema.json`** — THE schema for pipeline nodes (recognition types, action types, node properties)
2. **`interface_import.schema.json`** — Schema for task files (understanding pipeline_override context)
3. **`interface.schema.json`** — Main PI schema (for understanding resource/controller context)
4. **`custom.recognition.schema.json`** — Custom recognition extensions
5. **`custom.action.schema.json`** — Custom action extensions

### Step 2: Analyze Pipeline Fragment

1. Parse the provided JSON snippet
2. Identify all node names and their properties
3. Build a node graph: parent→children mapping via `next[]`
4. Extract key configuration: `recognition`, `action`, `enabled`, `next`, `interrupt`, `sub`, `on_error`

### Step 3: Validate Against Schema

Check each node against `pipeline.schema.json`:

1. **Required fields**: Ensure all mandatory properties are present
2. **Type validation**: Verify field types match schema definitions
3. **Value constraints**: Check enum values, numeric ranges, string patterns
4. **Recognition validation**: Ensure recognition type and parameters are correct
5. **Action validation**: Ensure action type and parameters are correct

### Step 4: Analyze Node Semantics via Naming Conventions

Use [Pipeline Node Naming Specification](../../../docs/pipeline-node-naming.md) to understand node roles and relationships:

1. **Parse node name structure**: Extract `<Domain><ActionOrObject><Role>` components
2. **Identify node type by role suffix**:
    - `Main` → Entry point node, organizes subsequent nodes
    - `Flow` → Orchestration node, no direct recognition/action
    - `Enter<Page>` → Navigation node, clicks to enter a page
    - `On<Page>Page` / `Visible` → State detection node, checks if on page/UI visible
    - `Click<Object>` / `Select<Object>` / `Claim<Object>` → Action node, performs interaction
    - `Confirm<Object>` → Confirmation node, handles confirmation dialogs
    - `Scroll<Direction>` / `Swipe<Object>` → Scroll action node
    - `End` / `EndTask` → Terminal node, ends flow
    - `Entered` → Success sentinel, confirms navigation completed
3. **Verify domain consistency**: Ensure all nodes in same module use same domain prefix
4. **Validate naming correctness**:
    - PascalCase format
    - No prohibited patterns (underscore prefix, numeric prefix, snake_case, camelCase, generic names)
    - Semantic accuracy (name reflects function, not implementation)
5. **Detect naming-code mismatches**:
    - Node named `Click<Object>` but has no `action` → Possible error
    - Node named `Visible` but has `action: Click` → Should be `Click<Object>`
    - Node named `Flow` but has recognition parameters → Should be pure orchestration
    - Node named `Enter<Page>` but no `next` retry → Missing success sentinel

### Step 5: Analyze Node Relationships

Check structural integrity:

1. **Reference validation**: All `next[]`, `interrupt[]`, `sub[]` targets must exist in the pipeline
2. **Circular detection**: Identify infinite loops in `next` chains
3. **Orphan detection**: Find nodes that are never referenced
4. **Entry points**: Identify root nodes (not referenced by any other node)
5. **Dead ends**: Find nodes with no `next` that aren't terminal actions

### Step 6: Identify Common Issues

See [Debug Rules Reference](references/debug-rules.md) for comprehensive rules.

Quick checks:

- **Missing recognition**: Nodes that should match something but have no recognition
- **Wrong recognition type**: Using TemplateMatch when OCR would be better
- **Missing action**: Nodes that do nothing after recognition
- **Incorrect ROI**: Recognition area too large/small or incorrectly positioned
- **Threshold issues**: Template matching thresholds too high/low
- **Performance bottlenecks**: Unnecessary nodes, redundant recognitions
- **Logic errors**: Wrong `next` sequences, missing interrupt handling

### Step 7: Generate Optimization Recommendations

Suggest improvements:

1. **Performance**: Reduce unnecessary recognitions, optimize ROI, adjust thresholds
2. **Maintainability**: Simplify node structure, reduce nesting depth
3. **Reliability**: Add error handling, improve recognition accuracy
4. **Readability**: Add `desc` documentation, use meaningful node names

### Step 8: Provide Corrected Solution

If issues are found:

1. Generate corrected JSON snippet
2. Explain each change made
3. Provide before/after comparison
4. Suggest testing approach

## Output Format

Structure your response as:

```
## Analysis Summary
- Total nodes: X
- Issues found: Y
- Optimizations suggested: Z

## Node Semantic Analysis
| Node Name | Type | Domain | Role | Expected Behavior |
|-----------|------|--------|------|-------------------|
| ShopEnterExchangePage | Action | Shop | Enter<Page> | Click to enter exchange page |
| ShopOnExchangePage | Detection | Shop | On<Page>Page | Check if on exchange page |
| CommonConfirmReward | Action | Common | Confirm<Object> | Confirm reward dialog |

## Issues Found
### Critical Issues
1. [Issue description] - [Node name] - [Explanation]

### Warnings
1. [Warning description] - [Node name] - [Explanation]

### Naming-Code Mismatches
1. [Node name] - [Expected behavior from name] - [Actual behavior in code]

### Suggestions
1. [Suggestion] - [Benefit]

## Corrected Pipeline
[JSON snippet with fixes applied]

## Explanation of Changes
1. [Change 1]: [Why it was needed]
2. [Change 2]: [Why it was needed]

## Testing Recommendations
- [How to test the fixes]
```

## Gotchas

- **`pipeline.schema.json` uses V1/V2 patterns**: Some nodes use `recognition` (V1), others use `type` (V2) — both are valid
- **`next` is optional**: Nodes without `next` are valid terminal nodes
- **`enabled` default is true**: Nodes are enabled by default unless explicitly set to false
- **`interrupt` and `sub` have different semantics**: `interrupt` pauses current execution, `sub` runs in parallel
- **ROI format**: Can be `[x, y, w, h]` array or string reference to previous node
- **Template paths**: Relative to `image/` folder, not project root
- **OCR expected**: Supports regex patterns, not just literal strings
- **Custom nodes**: May have additional properties not in base schema
- **Node naming**: Must follow PascalCase with Domain + ActionOrObject + Role format

## References

- [Pipeline Schema](../../../tools/schema/pipeline.schema.json)
- [Pipeline Node Naming Specification](../../../docs/pipeline-node-naming.md)
- [Debug Rules](references/debug-rules.md)
- [Common Patterns](references/common-patterns.md)
