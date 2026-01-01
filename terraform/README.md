# Terraform Infrastructure

Production-grade Terraform setup using a **single environment** and **multiple reusable modules**.
This repository follows best practices for **state isolation**, **module design**, and **safe deployments**.

---

## 📁 Structure

```
terraform/
├── modules/                     # Reusable infrastructure modules
│   ├── network/                 # VPC, subnets, routing
│   ├── identity/                # Cognito / IAM / auth
│   ├── compute/                 # ECS / EC2 / Lambda
│   └── storage/                 # S3, backups
│
├── env/                         # Single environment (prod)
│   ├── backend.tf               # Remote state configuration
│   ├── providers.tf             # Provider configuration
│   ├── main.tf                  # Module composition
│   ├── variables.tf             # Environment inputs
│   ├── terraform.tfvars         # Environment values (not committed if sensitive)
│   └── outputs.tf               # Exposed outputs
│
├── .terraform.lock.hcl          # Provider lock file (committed)
├── .gitignore
└── README.md
```

---

## 🧠 Design Principles

- **Single environment** with a clear state boundary
- **Modules are pure** (no backend, no providers)
- **Environment (`env/`) owns state and providers**
- **Explicit inputs and outputs** between modules
- **Remote state with locking** (S3 + DynamoDB)

---

## 🔐 State Management

Terraform state is stored remotely using S3 and locked with DynamoDB.

**Location:** `env/backend.tf`

```hcl
terraform {
  backend "http" {}
}
````

> ⚠️ Never commit `.tfstate` files to Git.

---

## 📦 Modules

Modules are reusable building blocks and **do not**:

* Define backends
* Define providers
* Know about environments (dev/prod)

Example:

```
modules/network/
├── vpc.tf
├── variables.tf
└── outputs.tf
```

Outputs from one module are consumed by another using:

```hcl
module.network.vpc_id
```

---

## 🌍 Environment (`env/`)

The `env/` directory is the **only place Terraform commands are executed**.

It:

* Configures providers
* Configures backend
* Wires modules together
* Supplies real values

Example:

```hcl
module "network" {
  source   = "../modules/network"
  vpc_cidr = var.vpc_cidr
}
```

---

## ▶️ How to Run Terraform

```bash
cd terraform/env

terraform init
terraform plan -out=tfplan
terraform apply tfplan
```

---

## 🧪 Variables & Secrets

* `variables.tf` defines required inputs
* `terraform.tfvars` supplies environment values
* Sensitive values **should not be committed**
* Prefer:

  * AWS Secrets Manager
  * SSM Parameter Store
  * CI/CD secrets

---

## 🚫 What NOT to Commit

| File / Folder         | Commit     |
| --------------------- | ---------- |
| `.terraform/`         | ❌ No       |
| `*.tfstate`           | ❌ No       |
| `terraform.tfvars`    | ⚠️ Depends |
| `.terraform.lock.hcl` | ✅ Yes      |

---

## 🚀 Scaling Later

This structure scales cleanly to multiple environments:

```
env/
├── dev/
├── staging/
└── prod/
```

No module changes required.

---

## 🧭 Terraform Mental Model

| Concept | Meaning         |
| ------- | --------------- |
| Module  | Library         |
| Env     | Deployment      |
| State   | Source of truth |
| Plan    | Change diff     |
| Apply   | Release         |

---

## ✅ Status

✔ Production-ready
✔ Safe state management
✔ Clean module boundaries
✔ Easy to extend

---

## 📌 Maintainer

Subroto Kumar
subrotokumar.dev
