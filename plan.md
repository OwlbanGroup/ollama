### Plan to Address "import cycle not allowed in test" Issue

1. **Refactor `integration/llm_test.go`**:
   - Move the relevant test functions to a new package or file that does not create a circular dependency with the `api` package.
   - Ensure that any necessary interfaces or types used in the tests are defined in a way that avoids direct imports from `api` if possible.

2. **Review and Update Imports**:
   - Check all import statements in the affected files (`common/runner.go`, `cmd/cmd_test.go`, `integration/llm_test.go`, and `api/model.go`) to ensure that they do not create circular dependencies.
   - Remove or replace any imports that lead to cycles.

3. **Test the Changes**:
   - After making the changes, run the tests to ensure that the import cycle issue is resolved and that all tests pass successfully.

4. **Documentation**:
   - Update any relevant documentation to reflect the changes made to the structure of the code and the organization of tests.

### Follow-up Steps:

- Verify the changes in the files.
- Confirm with the user for any additional requirements or modifications.
