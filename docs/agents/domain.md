# Domain Docs

本專案採 single-context 配置：

- 根目錄 `GLOSSARY.md`：專案領域詞彙。
- `docs/adr/`：架構決策紀錄。

## 探索前閱讀

探索程式碼前，先讀取 `GLOSSARY.md`，以及 `docs/adr/` 中與本次工作相關的決策。

文件或目錄不存在時，直接繼續工作，不需回報缺漏或預先建立。
當詞彙或決策確定後，由 domain-modeling skill 按需建立。

## 使用專案詞彙

Issue 標題、重構提案、假設及測試名稱等輸出，使用 `GLOSSARY.md` 定義的名稱，避開其明列應避免的同義詞。

若需要的概念尚未定義，先確認是否符合專案語言；
確有缺口時，記錄供 domain-modeling 處理。

## 架構決策衝突

提案若與既有 ADR 衝突，明確指出該 ADR 與重新討論的理由，不直接覆蓋原決策。
