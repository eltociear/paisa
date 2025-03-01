# CHANGELOG

### 0.5.2 (2023-09-22)

* Add Desktop app
* Support password protected PDF on import page
* Bug fixes

## Breaking Changes :rotating_light:

* The structure of price code configuration has been updated to make
  it easier to add more price provider in the future. In addition to
  the code, the provider name also has to be added. Refer the
  [config](https://paisa.fyi/reference/config/) documentation for more details

```diff
     type: mutualfund
-    code: 122639
+    price:
+      provider: in-mfapi
+      code: 122639
     harvest: 365
```


### 0.5.0 (2023-09-16)

* Add config page
* Embed ledger binary inside paisa
* Bug fixes

### 0.4.9 (2023-09-09)

* Add [search query](https://paisa.fyi/reference/bulk-edit/#search) support in transaction page
* Spends at child accounts level would be included in the budget of
  parent account.
* Fix the windows build, which was broken by the recent changes to
  ledger import
* Bug fixes

### 0.4.8 (2023-09-01)

* Add budget
* Add hierarchial cash flow
* Switch from float64 to decimal
* Bug fixes


### 0.4.7 (2023-08-19)

* Add dark mode
* Add bulk transaction editor
