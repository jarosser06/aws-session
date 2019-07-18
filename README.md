AWS Session Tool
================
Easy way to use AWS IAM Account Federation to retrieve temporary credentials.

Sample Configuration File:
```yaml
---
accounts:
- aws_access_key_id: 'REDACTED
  aws_secret_access_key: 'REDACTED'
  mfa_role: arn:aws:iam::012345678:mfa/jim
  aliases:
    - name: sandbox
      account_number: 032453343343
      role: Administrator
    - name: production
      account_number: 203433434334
      role: Administrator
    - name: preprod
      account_number: 102343433034
      role: Administrator
- aws_access_key_id: 'REDACTED'
  aws_secret_access_key: 'REDACTED'
  aliases:
    - name: personal-sandbox-read
      account_number: 033430343343
      role: readonly
    - name: personal-sandbox-admin
      account_number: 033430343343
      role: administrator
```
