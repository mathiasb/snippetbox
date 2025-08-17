# For DevOps, you MUST follow these principles
- Changes are documented in a changelog document
- Use version control strictly and use semantic versioning
- Use Git for version control
- Version control according to these steps
 1. Do the implementation
 2. Ensure tests are passing
 3. Stage all changes in git
 4. Commit the staged changes with a succinct and clear message
 5. Figure out the most appropriate semantic version
 6. Update the changelog documentation
 7. Stage and commit the changelog
 8. Tag the repo with the version and a one line message such as "{{version}}: {{one line description of the version}}"
- If there is a remote staging repo, push the changes and the tag to the staging repo

