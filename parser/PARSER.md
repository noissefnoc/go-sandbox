# パーサーライブラリのテスト

## alecthomas/participle

`LL(k)` 用のパーサらしい。つまりEBNFで表されるものは対応可能。

で

* Lexer
    * Regexp or EBNF で書く
* Parser
    * タグ付き構造体にマッピング

という方法か。