# Bitwarden Vault Bootstrap System

This vault subsystem restores:

- SSH keys  
- Multiple AWS SSO profiles  
- Multiple AWS credential profiles  
- Temporary session credentials  
- Per-profile regions  
- Arbitrary environment secrets  

## Usage


cd ~/workspace/dotfiles/vault
./bootstrap-vault.sh


This unlocks Bitwarden, caches the session, and restores:

- ~/.ssh/id_ed25519
- ~/.aws/config
- ~/.aws/credentials
- ~/.local/env.secrets

## Bitwarden Items Required

### *SSH-Primary*
Fields:
- private_key
- public_key

### *AWS-Master*
Fields:
- SSO_DEV_START_URL  
- SSO_DEV_REGION  
- SSO_PROD_START_URL  
- SSO_PROD_REGION  
- PERSONAL_AWS_ACCESS_KEY_ID  
- PERSONAL_AWS_SECRET_ACCESS_KEY  
- BWS_AWS_ACCESS_KEY_ID  
- BWS_AWS_SECRET_ACCESS_KEY  
- PROD_WB_ACCESS_KEY_ID  
- PROD_WB_SECRET_ACCESS_KEY  
- PROD_WB_SESSION_TOKEN  

### *Environment-Secrets*
Fields:  
Any key/value pairs you want exported.

---

