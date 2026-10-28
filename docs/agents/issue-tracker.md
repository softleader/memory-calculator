# Issue tracker: GitHub

本專案的 issues 與 specs 存放在 `softleader/memory-calculator` 的 GitHub Issues，透過 `gh` CLI 操作。

## 操作慣例

- 首次建立 issue 前，以 `gh repo view --json nameWithOwner` 確認目標為 `softleader/memory-calculator`，避免誤寫到上游或 skills 來源儲存庫。
- 建立：`gh issue create --title "..." --body-file <file>`。
- 讀取：`gh issue view <number> --json number,title,body,labels,comments`。
- 列出：`gh issue list --state open --json number,title,body,labels`，依需求加入標籤與狀態篩選。
- 更新內文：`gh issue edit <number> --body-file <file>`。
- 留言：`gh issue comment <number> --body-file <file>`。
- 加上或移除標籤：`gh issue edit <number> --add-label "..."` 或 `--remove-label "..."`。
- 關閉：先留言說明結果，再執行 `gh issue close <number>`；即使 issue 已關閉，說明仍以獨立留言發布。
- 多行內文透過 UTF-8 檔案傳入，保留 Markdown 換行。

## 工作關係

- 用原生 sub-issue 表示工作拆分，用原生 issue dependency 表示必須等待的執行順序。
- 操作前查看已安裝 `gh` 的 help；CLI 未支援時，使用 GitHub API。
- API 要求 issue database ID 時，透過 `gh api repos/softleader/memory-calculator/issues/<number> --jq .id` 取得，不使用 issue number 或 node ID 代替。
- 若原生關係不可用，子項目內文以 `Part of #<parent>` 或 `Blocked by: #<number>` 記錄；父項目以工作清單連結子項目。
- 所有 blocker 關閉後，工作才算解除阻擋。
- 工作類型使用原生 issue Type；triage 標籤表示處理狀態。
- 完成 issue 的 PR 使用 `Closes #<number>`；僅相關的 PR 使用 `Refs #<number>`。

## Pull requests as a triage surface

**PRs as a request surface: no.**

## Skills 用語

- 「publish to the issue tracker」表示建立 GitHub issue。
- 「fetch the relevant ticket」表示讀取對應 GitHub issue。
- GitHub issues 與 PR 共用編號；遇到不明確的編號時，先確認項目類型。

## Wayfinding

- Map 是標記 `wayfinder:map` 的 issue，內文保存 Notes、Decisions-so-far 與 Fog。
- 子工作透過 sub-issue 連結 map，並使用 `wayfinder:research`、`wayfinder:prototype`、`wayfinder:grilling` 或 `wayfinder:task`。
- 選取 map 中尚未關閉、未指派且無未關閉 blocker 的第一個子工作。
- 認領時指派給執行者；完成時留言、關閉子工作，再將結果連結加入 map 的 Decisions-so-far。
