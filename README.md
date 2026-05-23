# stackdiff

> Compare two Terraform state files and output a human-readable drift summary.

---

## Installation

```bash
go install github.com/yourusername/stackdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stackdiff.git
cd stackdiff
go build -o stackdiff .
```

---

## Usage

```bash
stackdiff <baseline.tfstate> <current.tfstate>
```

**Example:**

```bash
stackdiff prod-baseline.tfstate prod-current.tfstate
```

**Sample output:**

```
~ aws_instance.web_server
    instance_type: "t2.micro" → "t3.medium"

+ aws_s3_bucket.logs
    (new resource)

- aws_security_group.old_rule
    (removed resource)

Summary: 1 changed, 1 added, 1 removed
```

### Flags

| Flag | Description |
|------|-------------|
| `--json` | Output diff in JSON format |
| `--no-color` | Disable colored output |
| `--ignore-metadata` | Skip metadata-only changes |

---

## Requirements

- Go 1.21+
- Terraform state files in JSON format (`.tfstate`)

---

## License

This project is licensed under the [MIT License](LICENSE).