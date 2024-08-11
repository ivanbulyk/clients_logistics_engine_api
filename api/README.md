# Contribute instructions

Prerequisites:

- [Go](https://go.dev/doc/install)
- [Buf for proto-files management](https://buf.build/docs/installation)

Check the `Makefile` in root of the project for compiling go code from proto-files.

## Protobuf generated files

```text
├── api
│    ├── v1                     // Proto sources
│    └── logistics.swagger.json // Generated from proto sources
└── internal
    └── app
        └── logistics
            └── api
                └── v1          // Buf will generate go code here
```
