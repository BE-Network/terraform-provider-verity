# Preparing the generated SDK

`generate_openapi_sdk.sh` invokes this tool on temporary generator output before
both `--check` and `--write`. It keeps only top-level Go SDK sources, removes
generated tests and module files, and runs Go formatting. Generated docs, the
duplicate specification, publishing scripts, and generator bookkeeping are not
retained in `openapi/`.

`call_api.go.tmpl` is the provider-owned customization for the generated
`APIClient.callAPI` method. The tool locates the method through the Go parser and
refuses a changed signature. The customization is embedded in the tool; it never
copies behavior from the current generated SDK.

Authentication request and response bodies are excluded from debug logs.
Authorization and session headers are redacted on copies, leaving the actual
HTTP exchange intact. Non-authentication bodies remain available for debugging
and may contain resource secrets. Redacting resource fields marked sensitive in
the registry is a separate follow-up.
Regression tests live in `tests/unit/sdk/`; do not add tests inside `openapi/`,
which is replaced by regeneration.

The reconciled 6.6 baseline differs from earlier output mainly because canonical
JSON sorts keys: model fields and constructors are reordered, and identical
inline schemas acquire the name of the first sorted occurrence (`DenyInner`
instead of `PermitInner` for shared ACL models). Go formatting contributes
additional cosmetic differences. The schema additions in `cda5f49` are present
in the committed inputs; no model-field patch is required. The intentional auth
logging changes are now maintained here instead of manually restored after each
generation.

Use named fields for required string members such as authentication credentials;
generated constructor argument order can change without a compiler error.

Pipeline: canonical datacenter/campus inputs → `tools/process_swagger.py` → pinned
OpenAPI Generator → this tool → `openapi/`. Preprocessing merges modes, removes
excluded endpoints, changes numeric width/nullability, and transforms request
wrappers. Keep changes in those inputs or preprocessing rather than editing
generated models.
